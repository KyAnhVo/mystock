package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KyAnhVo/mystock/config"
	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/jackc/pgx/v5"
)

type StockDataHandler struct {
	dbQuerier    *db.DBQueryMachine
	apiRequester *http.Client
	logger       *slog.Logger
}

func NewStockHandler(dbQuerier *db.DBQueryMachine, logger *slog.Logger) *StockDataHandler {
	client := &http.Client{}
	return &StockDataHandler{
		dbQuerier:    dbQuerier,
		apiRequester: client,
		logger:       logger,
	}
}

// Gets the information of a ticker (must be in US)
//
// path: /api/ticker/{ticker}
func (stockHandler *StockDataHandler) OverviewTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	var response struct {
		Ticker      string  `json:"ticker"`
		Name        string  `json:"name"`
		CIK         *string `json:"cik"`
		Description *string `json:"description"`
	}

	// Find in DB first. If not there, we find via handler.
	err := stockHandler.dbQuerier.Querier.QueryRow(
		r.Context(),
		`SELECT ticker, name, cik, description FROM market_data.ticker
 			WHERE ticker = $1 `,
		ticker,
	).Scan(
		&response.Ticker,
		&response.Name,
		&response.CIK,
		&response.Description,
	)

	// if found in DB, return.
	if err == nil {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// if error is "NOT FOUND", we try to query from Massive
		stockHandler.logger.Warn(
			"StockDataHandler: OverviewTicker: stock not found in DB",
			"ticker", ticker,
		)
		cfg := config.GetCfg()
		resp, err := stockHandler.apiRequester.Get(
			cfg.StockApiHeader +
				"/v3/reference/tickers/" + ticker +
				"?apiKey=" + cfg.StockApiKey,
		)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "server error")
			stockHandler.logger.Error(
				"StockDataHandler: OverviewTicker: cannot reach stock API",
			)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			// if found, we store it in our db first then return to user.
			type respContent struct {
				Ticker      string  `json:"ticker"`
				Name        string  `json:"name"`
				CIK         *string `json:"cik"`
				Description *string `json:"description"`
			}
			var respBody struct {
				Status string      `json:"status"`
				Result respContent `json:"results"`
			}
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintln(w, "server error")
				return
			}

			// shouldnt happen since code is 200, but oh well
			if respBody.Status != "OK" {
				w.WriteHeader(404)
				fmt.Fprintln(w, "no ticker found")
				return
			}

			// Store this in our db
			_, err := stockHandler.dbQuerier.Querier.Exec(
				r.Context(),
				"INSERT INTO market_data.ticker (ticker, name, cik, description) VALUES ($1, $2, $3, $4)",
				respBody.Result.Ticker,
				respBody.Result.Name,
				respBody.Result.CIK,
				respBody.Result.Description,
			)
			if err != nil {
				stockHandler.logger.Warn("failed to store ticker in db", "error", err)
			}

			// and return the data
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(respBody.Result)
		} else {
			// if not found, just say not found
			w.WriteHeader(404)
			fmt.Fprintln(w, "ticker not found")
		}
	} else {
		w.WriteHeader(500)
		fmt.Fprintln(w, "server error")
	}
}

// Get a bunch of tickers
//
// path: /api/tickers?page=<page>&page_size=<page_size>
func (stockHandler *StockDataHandler) GetTickers(w http.ResponseWriter, r *http.Request) {
	MAX_PAGE_SIZE := 100

	// Get page, page size, and starting index from query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "invalid page")
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "invalid page_size")
		return
	}
	queryFromIndex := page * pageSize
	if queryFromIndex < 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, "invalid page")
		return
	}
	if pageSize <= 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, "invalid page_size")
		return
	}
	if pageSize > MAX_PAGE_SIZE {
		w.WriteHeader(400)
		fmt.Fprintln(w, "page_size too large, limit to ", MAX_PAGE_SIZE)
		return
	}

	type ticker struct {
		Name   string `json:"name"`
		Ticker string `json:"ticker"`
	}

	// query the DB for tickers, starting from the calculated index
	rows, err := stockHandler.dbQuerier.Querier.Query(
		r.Context(),
		"SELECT name, ticker FROM market_data.ticker ORDER BY ticker LIMIT $1 OFFSET $2",
		pageSize,
		queryFromIndex,
	)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "server error")
		return
	}
	defer rows.Close()

	tickers := make([]ticker, 0, pageSize)
	for rows.Next() {
		var item ticker
		err = rows.Scan(&item.Name, &item.Ticker)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "server error")
			return
		}
		tickers = append(tickers, item)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "server error")
		return
	}

	// send all tickers queried from the DB to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(tickers)
}

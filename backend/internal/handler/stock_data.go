package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/KyAnhVo/mystock/config"
	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/jackc/pgx/v5"
)

type StockDataHandler struct {
	dbQuerier    db.DBQueryMachine
	apiRequester *http.Client
}

// Gets the information of a ticker (must be in US)
func (stockHandler *StockDataHandler) OverviewTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	var response struct {
		Ticker      string `json:"ticker"`
		Name        string `json:"name"`
		CIK         string `json:"cik"`
		Description string `json:"description"`
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
		cfg := config.GetCfg()
		resp, err := stockHandler.apiRequester.Get(
			cfg.StockApiHeader +
				"/v3/reference/ticker/" + ticker +
				"?apiKey=" + cfg.StockApiKey,
		)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "server error")
			return
		}

		if resp.StatusCode == 200 {
			// if found, we store it in our db first then return to user.
			type respContent struct {
				Ticker      string `json:"ticker"`
				Name        string `json:"name"`
				CIK         string `json:"cik"`
				Description string `json:"description"`
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
			stockHandler.dbQuerier.Querier.Exec(
				r.Context(),
				"INSERT INTO market_data.ticker (ticker, name, cik, description) VALUES ($1, $2, $3, $4)",
				respBody.Result.Ticker,
				respBody.Result.Name,
				respBody.Result.CIK,
				respBody.Result.Description,
			)

			// and return the data
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

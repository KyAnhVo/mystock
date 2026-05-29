package db

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KyAnhVo/mystock/config"
	"github.com/jackc/pgx/v5"
)

// Reset the schema using data in ./internal/db/migrations
func (db *DBQueryMachine) ResetSchema() error {
	ctx := context.Background()
	entries, err := os.ReadDir("internal/db/migrations")
	if err != nil {
		return err
	}

	tx, err := db.Querier.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".down.sql") {
			continue
		}
		content, err := os.ReadFile("./internal/db/migrations/" + entry.Name())
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, string(content))
		if err != nil {
			return err
		}
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		content, err := os.ReadFile("./internal/db/migrations/" + entry.Name())
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, string(content))
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// Get the tickers from API to DB
// Note: we do not care if the
func (db *DBQueryMachine) GetTickerInfoFromAPI() error {
	API_RATE := 5
	requestsMade := 0

	ctx := context.Background()
	cli := &http.Client{}
	conf := config.GetCfg()

	// Load tickers from DB first
	tickersInDb := make(map[string]struct{})
	rows, err := db.Querier.Query(ctx, "SELECT ticker FROM market_data.ticker")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var ticker string
		err := rows.Scan(&ticker)
		if err != nil {
			return err
		}
		tickersInDb[ticker] = struct{}{}
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}

	// Loop through all tickers, while there is still a new url to fetch
	responseStruct := &struct {
		Count   int     `json:"count"`
		NextUrl *string `json:"next_url"`
		Results []struct {
			CIK      *string `json:"cik"`
			Name     string  `json:"name"`
			Ticker   string  `json:"ticker"`
			Exchange *string `json:"primary_exchange"`
		} `json:"results"`
		Status string `json:"status"`
	}{}

	eofErr := errors.New("No more lines")
	for {
		err := func() error {
			// only API_LIMIT_PER_MINUTE requests per minute
			if requestsMade == API_RATE {
				time.Sleep(1 * time.Minute)
				requestsMade = 0
			}
			requestsMade += 1

			// Get the url for api querying
			var nextUrl string
			if responseStruct.NextUrl == nil {
				nextUrl = conf.StockApiHeader +
					"/v3/reference/tickers?market=stocks&active=true&order=asc&limit=1000&sort=ticker"
			} else {
				nextUrl = *responseStruct.NextUrl
			}

			// query from API then put in responseStruct
			resp, err := cli.Get(nextUrl + "&apiKey=" + conf.StockApiKey)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(&responseStruct)
			if err != nil {
				return err
			}

			// If there is any error, return ASAP that error.
			if resp.StatusCode != 200 {
				return errors.New("API access unsuccessful")
			} else if responseStruct.Status != "OK" {
				return errors.New("API access unsuccessful: status " + responseStruct.Status)
			}

			// Continuously build up the query. If something is in the DB,
			// we ignore it. Else we insert it into the DB.
			rows := [][]any{}
			for _, r := range responseStruct.Results {
				_, isInDb := tickersInDb[r.Ticker]
				if isInDb {
					continue
				}
				rows = append(rows, []any{r.Ticker, r.CIK, r.Name, r.Exchange})
				tickersInDb[r.Ticker] = struct{}{}
			}
			rowsRead, err := db.Querier.CopyFrom(
				ctx,
				pgx.Identifier{"market_data", "ticker"},
				[]string{"ticker", "cik", "name", "exchange"},
				pgx.CopyFromRows(rows),
			)
			if err != nil {
				return err
			}
			if rowsRead != int64(len(rows)) {
				return errors.New("Cannot write to DB")
			}

			// Return eofErr if this is the final one, else return nil.
			if responseStruct.NextUrl == nil || len(rows) == 0 || len(responseStruct.Results) < 1000 {
				return eofErr
			}
			return nil
		}()

		if err != nil {
			if errors.Is(err, eofErr) {
				break
			} else {
				return err
			}
		}
	}

	return nil
}

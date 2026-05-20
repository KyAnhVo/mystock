package db

import (
	"context"
	"os"
	"strings"

	"github.com/KyAnhVo/mystock/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBQueryMachine struct {
	Querier *pgxpool.Pool
}

// Creates a new DB connection using URL from .env
func Init() (*DBQueryMachine, error) {
	config := config.GetCfg()
	dbPool, err := pgxpool.New(context.Background(), config.DBConn)
	if err != nil {
		return nil, err
	}

	return &DBQueryMachine{Querier: dbPool}, nil
}

func (db *DBQueryMachine) RunQuery(query string) error {
	ctx := context.Background()
	_, err := db.Querier.Exec(ctx, query)
	return err
}

// Runs a sequence of queries, atomic
func (db *DBQueryMachine) RunQueries(queries []string) error {
	ctx := context.Background()
	tx, err := db.Querier.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, query := range queries {
		_, err := tx.Exec(ctx, query)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

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

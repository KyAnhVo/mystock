package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/KyAnhVo/mystock/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBQueryMachine struct {
	querier *pgxpool.Pool
}

func Init() (*DBQueryMachine, error) {
	dbPool, err := pgxpool.New(context.Background(), config.Cfg.DBConn)
	if err != nil {
		return nil, err
	}

	return &DBQueryMachine{querier: dbPool}, nil
}

func (db *DBQueryMachine) ResetSchema() error {
	ctx := context.Background()
	entries, err := os.ReadDir("internal/db/migrations")
	if err != nil {
		return err
	}

	tx, err := db.querier.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		content, err := os.ReadFile("./internal/db/migrations/" + entry.Name())
		if err != nil {
			return err
		}
		_, err = db.querier.Exec(ctx, string(content))
		if err != nil {
			fmt.Println("failed to exec:", entry.Name(), err.Error())
		}
	}

	return tx.Commit(ctx)
}

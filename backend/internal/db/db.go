package db

import (
	"context"

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

package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

import (
	"github.com/KyAnhVo/mystock/config"
)

type DBQueryMachine struct {
	querier *pgxpool.Pool
}

func Init() (*DBQueryMachine, error) {
	dbPool, err := pgxpool.New(context.Background(), config.Cfg.DBConn)
	if err != nil {
		return nil, errors.New("cannot create DB")
	}

	return &DBQueryMachine{querier: dbPool}, nil
}

func (*DBQueryMachine) resetSchema() error {

	return nil
}

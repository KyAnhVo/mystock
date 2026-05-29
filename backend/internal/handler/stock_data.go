package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/KyAnhVo/mystock/internal/util"
)

type StockDataHandler struct {
	db           *db.DBQueryMachine
	apiRequester *http.Client
	logger       *slog.Logger
}

func NewStockDataHandler(db *db.DBQueryMachine, logger *slog.Logger) *StockDataHandler {
	client := &http.Client{}
	return &StockDataHandler{
		db:           db,
		apiRequester: client,
		logger:       logger,
	}
}

// Returns the stock data (OHLCV + cummulated action)
//
// Path: /api/stock/<ticker>?fromdate=yyyy-mm-dd&todate=yyyy-mm-dd
func (handler *StockDataHandler) GetStockOn(w http.ResponseWriter, r *http.Request) {
	fromDateStr := r.URL.Query().Get("fromdate")
	toDateStr := r.URL.Query().Get("todate")

	var fromDateTime, toDateTime time.Time
	fromDate, err := util.StringToDate(fromDateStr)
	if err != nil {
		fromDateTime = time.Now()
	} else {
		fromDateTime = time.Date
	}

}

package handler

import (
	"encoding/json"
	"errors"
	"github.com/KyAnhVo/mystock/config"
	"net/http"
)

type Adress struct {
	Adress1    string `json:"address1"`
	Adress2    string `json:"address2"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

type TickerOverview struct {
	Name        string  `json:"name"`
	Market      string  `json:"market"`
	Description string  `json:"description"`
	Adress      *Adress `json:"address"`
}

func OverviewTicker(ticker string, handler *http.Client) (*TickerOverview, error) {
	config := config.GetCfg()
	resp, err := handler.Get(
		config.StockApiHeader +
			"/v3/reference/tickers/" + ticker +
			"?apiKey=" + config.StockApiKey,
	)

	if err != nil {
		return nil, errors.New("cannot get data")
	}

	type Response struct {
		Status string          `json:"status"`
		Result *TickerOverview `json:"results"`
	}

	var respStruct Response
	json.NewDecoder(resp.Body).Decode(&respStruct)
	if respStruct.Status != "OK" {
		return nil, errors.New("cannot get data")
	}
	if respStruct.Result == nil {
		return nil, errors.New("Cannot get data")
	}
	return respStruct.Result, nil
}

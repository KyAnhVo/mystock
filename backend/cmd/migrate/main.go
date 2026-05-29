package main

import (
	"fmt"

	"github.com/KyAnhVo/mystock/internal/db"
)

func main() {
	db, err := db.Init()
	if err != nil {
		fmt.Println("cannot connect to db")
		return
	}
	err = db.ResetSchema()
	if err != nil {
		fmt.Println("cannot reset schema:", err.Error())
		return
	}

	err = db.GetTickerInfoFromAPI()
	if err != nil {
		fmt.Println("cannot populate market_data.info:", err.Error())
		return
	}

	fmt.Println("reset schema completed")
}

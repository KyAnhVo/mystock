package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"sync"
)

import (
	"github.com/KyAnhVo/mystock/config"
	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/KyAnhVo/mystock/internal/handler"
)

func main() {
	godotenv.Load()
	config.Init()

	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		client := &http.Client{}
		res, err := handler.OverviewTicker("AAPL", client)
		if err != nil {
			mu.Lock()
			fmt.Println("Error getting ticker: ", err.Error())
			mu.Unlock()
		} else {
			data, _ := json.MarshalIndent(res, "", " ")
			mu.Lock()
			fmt.Println(string(data))
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		_, err := db.Init()
		if err != nil {
			mu.Lock()
			fmt.Println("Error getting DB: ", err.Error())
			mu.Unlock()
		} else {
			mu.Lock()
			fmt.Println("Connected to DB")
			mu.Unlock()
		}
	}()

	wg.Wait()

}

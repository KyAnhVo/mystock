package main

import (
	"fmt"

	"github.com/KyAnhVo/mystock/config"
	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.Init()

	db, err := db.Init()
	if err != nil {
		fmt.Println("cannot connect to db")
		return
	}
	db.ResetSchema()
}

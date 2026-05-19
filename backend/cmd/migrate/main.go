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
	} else {
		fmt.Println("reset schema successful")
	}
}

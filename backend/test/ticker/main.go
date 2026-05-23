package main

import (
	"fmt"
	"net/http"
)

func main() {
	client := &http.Client{}
	resp, err := client.Get("http://localhost:3001/api/ticker/AAPL")
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Println(resp.Body)
}

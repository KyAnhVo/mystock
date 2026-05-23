package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	client := &http.Client{}
	resp, err := client.Get("http://localhost:3001/api/ticker/AAPL")
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err reading body: " + err.Error())
		return
	}
	fmt.Println(string(body))

}

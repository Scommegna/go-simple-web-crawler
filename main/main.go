package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please insert exactly one URL")
    	return
	}

	httpLink := os.Args[1:]

	response, err := http.Get(httpLink[0])

	if err != nil {
		fmt.Printf("Error making http request: %s\n", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("Error reading the response body: %s\n", err)
		os.Exit(1)
	}

	htmlContent := string(body)

	links := GetLinks(htmlContent)

	fmt.Print(links)
}
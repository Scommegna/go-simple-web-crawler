package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	visited = make(map[string]struct{})
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	mu        sync.Mutex
	wg        sync.WaitGroup
	semaphore = make(chan struct{}, 10)
)

func crawl(url string) {
	defer wg.Done()

	mu.Lock()

	if _, ok := visited[url]; ok {
		mu.Unlock()
		return
	}

	visited[url] = struct{}{}
	mu.Unlock()

	semaphore <- struct{}{}
	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error making http request: %s\n", err)
		<-semaphore
		return
	}
	
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	<-semaphore

	if err != nil {
		return
	}

	links := GetLinks(string(body))

	for _, link := range links {
		wg.Add(1)
		go crawl(link)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please insert exactly one URL")
    	return
	}

	startURL := os.Args[1]

	wg.Add(1)
	go crawl(startURL)
	wg.Wait()

	fmt.Print(visited)
}
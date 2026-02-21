package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
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
	maxLinks int
)

func crawl(url string) {
	defer wg.Done()

	mu.Lock()

	if _, ok := visited[url]; ok {
		mu.Unlock()
		return
	}

	if maxLinks >0 && len(visited) >= maxLinks {
		mu.Unlock()
		return
	}

	visited[url] = struct{}{}
	mu.Unlock()

	semaphore <- struct{}{}
	response, err := client.Get(url)
	if err != nil {
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
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run . <url> [maxLinks]")
    	return
	}

	startURL := os.Args[1]

	if len(os.Args) == 3 {
		limit, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid maxLinks value. It should be an integer.")
			return
		}
		maxLinks = limit
	}

	wg.Add(1)
	go crawl(startURL)
	wg.Wait()

	fmt.Println("-----------")
	fmt.Println("Crawl finished")
	fmt.Println("-----------")

	fmt.Printf("Total unique links: %d\n\n", len(visited))

	var links []string
	for link := range visited {
		links = append(links, link)
	}

	sort.Strings(links)

	for i, link := range links {
		fmt.Printf("%3d. %s\n", i+1, link)
	}
}
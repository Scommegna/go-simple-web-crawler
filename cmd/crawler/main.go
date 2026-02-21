package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

var (
	visited = make(map[string]struct{})
	client  = &http.Client{
		Timeout: 5 * time.Second,
	}
	mu        sync.Mutex
	wg        sync.WaitGroup
	semaphore = make(chan struct{}, 10)
	maxLinks  int
)

func crawl(url string) {
	defer wg.Done()

	mu.Lock()

	if _, ok := visited[url]; ok {
		mu.Unlock()
		return
	}

	if maxLinks > 0 && len(visited) >= maxLinks {
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
	urlFlag := flag.String("url", "", "Starting URL (required)")
	limitFlag := flag.Int("limit", 0, "Maximum number of unique links (optional)")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Usage:")
		fmt.Println("  crawler -url=https://example.com [-limit=50]")
		return
	}

	maxLinks = *limitFlag

	wg.Add(1)
	go crawl(*urlFlag)
	wg.Wait()

	fmt.Println("==================================================")
	fmt.Println("Crawl finished")
	fmt.Println("==================================================")

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
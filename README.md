# Web Crawler (Go)

A simple concurrent web crawler written in Go.

This tool starts from a given URL, recursively discovers links, and
prints all unique URLs found.\
It supports an optional limit for the maximum number of links collected.

---

## Features

- Recursive crawling
- Concurrent requests using goroutines
- Thread-safe unique link tracking
- Semaphore-based concurrency control
- HTTP timeout protection
- Optional maximum link limit
- Clean and sorted output

---

## Installation

### Install locally

Clone the repository:

```bash
git clone https://github.com/yourusername/web-crawler.git
cd web-crawler
```

Install the binary:

```bash
go install ./cmd/crawler
```

Make sure your Go bin directory is in your PATH:

```bash
echo $PATH
```

If needed (macOS + zsh), add this to `~/.zshrc`:

```bash
export PATH=$PATH:$HOME/go/bin
```

Reload your shell:

```bash
source ~/.zshrc
```

Now you can run:

```bash
crawler -h
```

---

## Usage

```bash
crawler -url=<starting_url> [-limit=<max_links>]
```

### Required

- `-url` → Starting URL

### Optional

- `-limit` → Maximum number of unique links to collect\
  If omitted, crawling runs without limits.

---

## Examples

Crawl without limit:

```bash
crawler -url=https://go.dev
```

Crawl with limit of 50 links:

```bash
crawler -url=https://go.dev -limit=50
```

Show help:

```bash
crawler -h
```

---

## How It Works

- Uses goroutines for concurrent crawling
- Uses `sync.Mutex` to protect shared state
- Uses a `map[string]struct{}` as a Set to ensure uniqueness
- Uses a semaphore channel to limit concurrent HTTP requests
- Uses `http.Client` with timeout protection

---

## Project Structure

    web-crawler/
      go.mod
      cmd/
        crawler/
          main.go
      links.go

- `main.go` → CLI and crawling orchestration\
- `links.go` → HTML parsing and link extraction

---

## Notes

- The crawler does **not restrict domain scope**
- It may crawl external websites
- Be mindful when crawling large sites
- Concurrency is limited to avoid connection errors

---

## Requirements

- Go 1.20+ (recommended 1.22+)

Check your Go version:

```bash
go version
```

---

## License

MIT License

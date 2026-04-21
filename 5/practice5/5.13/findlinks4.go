package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func extractLinks(body []byte, baseURL string) []string {
	var links []string
	re := regexp.MustCompile(`href="([^"]+)"`)
	matches := re.FindAllSubmatch(body, -1)

	base, _ := url.Parse(baseURL)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		raw := string(match[1])
		ref, err := url.Parse(raw)
		if err != nil {
			continue
		}
		absolute := base.ResolveReference(ref).String()
		links = append(links, absolute)
	}
	return links
}

func savePage(urlStr string, body []byte) error {
	filename := strings.TrimPrefix(urlStr, "https://")
	filename = strings.TrimPrefix(filename, "http://")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	if filename == "" {
		filename = "index.html"
	}
	os.MkdirAll("downloaded", 0755)
	return os.WriteFile(filepath.Join("downloaded", filename), body, 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main <url>")
		return
	}

	startURL := os.Args[1]
	parsed, _ := url.Parse(startURL)
	domain := parsed.Hostname()

	queue := []string{startURL}
	visited := map[string]bool{}

	for len(queue) > 0 && len(visited) < 50 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		fmt.Println("Loading: ", current)

		resp, err := http.Get(current)
		if err != nil {
			fmt.Println(" ERROR:", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println(" STATUS:", resp.Status)
			continue
		}

		fmt.Printf(" Saved %d bytes\n", len(body))
		savePage(current, body)

		links := extractLinks(body, current)
		fmt.Printf(" Found %d links\n", len(links))

		for _, link := range links {
			parsedLink, _ := url.Parse(link)
			if parsedLink.Hostname() == domain && !visited[link] {
				queue = append(queue, link)
			}
		}
	}

	fmt.Println("\nDone. Check 'downloaded/' folder.")
}

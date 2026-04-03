package findlinks2

import (
	"fmt"
	findlinks "gobook/5/findlinks1"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		links, err := findLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "findlinks2: %v\n", err)
			continue
		}
		for _, links := range links {
			fmt.Println(links)
		}
	}
}

// findLinks выполняет HTTP-запрос GET для заданного url,
// выполняет синтаксический анализ ответа как HTML-документа
// и извлекает и возвращает ссылки.
func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("получение %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("анализ %s как HTML: %v", url, err)
	}
	return findlinks.Visit(nil, doc), nil
}

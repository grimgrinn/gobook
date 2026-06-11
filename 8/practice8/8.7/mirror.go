package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type item struct {
	url   string
	depth int
}

var tokens = make(chan struct{}, 20)
var domain string
var seen = struct {
	sync.Mutex
	m map[string]bool
}{m: make(map[string]bool)}

func main() {
	maxDepth := flag.Int("depth", 0, "максимальная глубина (0 = без ограничений)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("укажите URL")
	}
	startURL := args[0]

	// Нормализуем URL
	u, err := url.Parse(startURL)
	if err != nil {
		log.Fatal(err)
	}
	domain = u.Hostname()
	startURL = u.Scheme + "://" + domain + "/"

	worklist := make(chan []item)
	var n int

	n++
	go func() {
		worklist <- []item{{startURL, 0}}
	}()

	for ; n > 0; n-- {
		list := <-worklist
		for _, it := range list {
			seen.Lock()
			if seen.m[it.url] {
				seen.Unlock()
				continue
			}
			seen.m[it.url] = true
			seen.Unlock()

			if *maxDepth > 0 && it.depth >= *maxDepth {
				continue
			}

			n++
			go func(it item) {
				newItems := crawl(it)
				if len(newItems) > 0 {
					worklist <- newItems
				}
			}(it)
		}
	}
}

func crawl(it item) []item {
	fmt.Println(it.url)

	tokens <- struct{}{}
	resp, err := http.Get(it.url)
	<-tokens
	if err != nil {
		log.Printf("ошибка загрузки %s: %v", it.url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("неверный статус %s: %s", it.url, resp.Status)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ошибка чтения %s: %v", it.url, err)
		return nil
	}

	// Сохраняем страницу
	filename := savePage(it.url, body)
	if filename != "" {
		// Переписываем ссылки в сохраненном файле
		rewriteLocal(filename, it.url)
	}

	// Извлекаем ссылки
	links := extractLinks(body, it.url)
	var newItems []item
	for _, link := range links {
		if sameDomain(link, domain) {
			newItems = append(newItems, item{link, it.depth + 1})
		}
	}
	return newItems
}

func savePage(rawurl string, body []byte) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	// Убираем порт из host
	host := strings.Split(u.Host, ":")[0]
	path := u.Path
	if path == "" {
		path = "/"
	}
	// Добавляем .html, если путь не заканчивается на .html, .htm, .php и т.д.
	if !strings.Contains(path, ".") || strings.HasSuffix(path, "/") {
		if strings.HasSuffix(path, "/") {
			path += "index.html"
		} else {
			path += ".html"
		}
	}
	filename := filepath.Join("mirror", host, path)
	// Заменяем / на разделитель пути ОС
	filename = filepath.FromSlash(filename)

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Ошибка создания директории %s: %v", dir, err)
		return ""
	}
	if err := os.WriteFile(filename, body, 0644); err != nil {
		log.Printf("ошибка сохранения %s: %v", filename, err)
		return ""
	}
	return filename
}

func rewriteLocal(filename, originalURL string) {
	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("ошибка чтения %s: %v", filename, err)
		return
	}

	// Парсим HTML
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		log.Printf("ошибка парсинга %s: %v", filename, err)
		return
	}

	// Нормализуем оригинальный URL для сравнения
	base, _ := url.Parse(originalURL)

	var rewrite func(n *html.Node)
	rewrite = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Проверяем атрибуты href и src
			var attrName string
			if n.Data == "a" {
				attrName = "href"
			} else if n.Data == "img" || n.Data == "script" || n.Data == "link" {
				attrName = "src"
			}
			if attrName != "" {
				for i, a := range n.Attr {
					if a.Key == attrName {
						// Превращаем ссылку в абсолютную
						ref, err := url.Parse(a.Val)
						if err != nil {
							continue
						}
						absolute := base.ResolveReference(ref)
						// Если ссылка на тот же домен - переписываем в локальную
						if absolute.Hostname() == domain {
							localPath := urlToLocalPath(absolute.String())
							n.Attr[i].Val = localPath
						}
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rewrite(c)
		}
	}
	rewrite(doc)

	// Записываем обратно
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		log.Printf("ошибка рендера %s: %v", filename, err)
		return
	}
	if err := os.WriteFile(filename, buf.Bytes(), 06444); err != nil {
		log.Printf("Ошибка записи %s: %v", filename, err)
	}
}

func urlToLocalPath(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}
	host := strings.Split(u.Host, ":")[0]
	path := u.Path
	if path == "" {
		path = "/"
	}
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	} else {
		path += ".html"
	}

	fullPath := filepath.Join("mirror", host, path)
	absPath, _ := filepath.Abs(fullPath)
	return "file://" + filepath.ToSlash(absPath)
}

func extractLinks(body []byte, baseURL string) []string {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil
	}
	var links []string
	base, _ := url.Parse(baseURL)
	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					ref, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					absolute := base.ResolveReference(ref)
					links = append(links, absolute.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(doc)
	return links
}

func sameDomain(link, domain string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return u.Hostname() == domain
}

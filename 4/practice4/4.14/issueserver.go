package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Issue struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
}

var issues []Issue

func fetchIssues(owner, repo string) []Issue {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get error %v\n", err)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error %v\n", err)
	}

	var issues []Issue
	if err := json.Unmarshal(body, &issues); err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal error %v\n", err)
	}

	return issues
}

func issuesList(w http.ResponseWriter, r *http.Request) {
	const tmpl = `
	<h1>Issues</h1>
	<ul>
	{{range .}}
		<li><a href="/issue?number={{.Number}}">{{.Title}}</a></li>
	{{end}}
	</ul>`

	t := template.Must(template.New("list").Parse(tmpl))
	if err := t.Execute(w, issues); err != nil {
		log.Print(err)
	}
}

func issueDetail(w http.ResponseWriter, r *http.Request) {
	strNum := r.URL.Query().Get("number")
	if strNum == "" {
		http.NotFound(w, r)
		return
	}
	num, err := strconv.Atoi(strNum)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var found *Issue
	for i := range issues {
		if issues[i].Number == num {
			found = &issues[i]
			break
		}
	}
	if found == nil {
		http.NotFound(w, r)
		return
	}

	const tmpl = `
	<h1>{{.Title}}</h1>
	<p>Issue #{{.Number}}</p>
	<p>Author: {{.User.Login}}</p>
	<p>URL: <a href="{{.HTMLURL}}">{{.HTMLURL}}</a></p>
	<hr>
	<pre>{{.Body}}</pre>`

	t := template.Must(template.New("detail").Parse(tmpl))
	if err := t.Execute(w, found); err != nil {
		log.Print(err)
	}
}

func main() {
	owner := flag.String("o", "", "owner")
	repo := flag.String("r", "", "repo")
	flag.Parse()
	if *owner == "" || *repo == "" {
		fmt.Fprintf(os.Stderr, "need to enter owner -o and repo -r\n")
		os.Exit(1)
	}
	issues = fetchIssues(*owner, *repo)

	// fmt.Printf("%v", s)
	http.HandleFunc("/", issuesList)
	http.HandleFunc("/issue", issueDetail)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

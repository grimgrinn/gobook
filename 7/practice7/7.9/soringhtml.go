package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []Track{
	{"Go", "Moby", "Moby", 1992, time.Second * 3 * 37},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, time.Second * 4 * 36},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, time.Second * 4 * 24},
	{"Go", "Delilah", "From the Roots", 2012, time.Second * 3 * 30},
}

type ByColumns struct {
	t       []Track
	columns []string
}

func (b *ByColumns) Len() int      { return len(b.t) }
func (b *ByColumns) Swap(i, j int) { b.t[i], b.t[j] = b.t[j], b.t[i] }
func (b *ByColumns) Less(i, j int) bool {
	for _, col := range b.columns {
		switch col {
		case "Title":
			if b.t[i].Title != b.t[j].Title {
				return b.t[i].Title < b.t[j].Title
			}
		case "Artist":
			if b.t[i].Artist != b.t[j].Artist {
				return b.t[i].Artist < b.t[j].Artist
			}
		case "Album":
			if b.t[i].Album != b.t[j].Album {
				return b.t[i].Album < b.t[j].Album
			}
		case "Year":
			if b.t[i].Year != b.t[j].Year {
				return b.t[i].Year != b.t[j].Year
			}
		case "Length":
			if b.t[i].Length != b.t[j].Length {
				return b.t[i].Length < b.t[j].Length
			}
		}
	}
	return false
}

var tmpl = template.Must(template.New("table").Parse(`
<html>
<body>
<table border="1">
<tr>
	<th><a href="/?sort=Title,Artist,Year">Title</a></th>
    <th><a href="/?sort=Artist,Title,Year">Artist</a></th>
    <th><a href="/?sort=Album,Title,Year">Album</a></th>
    <th><a href="/?sort=Year,Title,Artist">Year</a></th>
    <th><a href="/?sort=Length,Title,Artist">Length</a></th>
</tr>
{{range .}}
<tr>
	<td>{{.Title}}</td>
	<td>{{.Artist}}</td>
	<td>{{.Album}}</td>
	<td>{{.Year}}</td>
	<td>{{.Length}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`))

func tracksHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "Title, Year" // по умолчанию
	}

	columns := strings.Split(sortBy, ",")

	// Копируем треки, чтобы не менять глобальные данные
	sorted := make([]Track, len(tracks))
	copy(sorted, tracks)

	// Сортируем
	bc := &ByColumns{t: sorted, columns: columns}

	sort.Sort(bc)

	err := tmpl.Execute(w, sorted)
	if err != nil {
		log.Print(err)
	}

}

func main() {
	http.HandleFunc("/", tracksHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

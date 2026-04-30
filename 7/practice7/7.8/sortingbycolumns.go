package main

import (
	"fmt"
	"sort"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// ByColumn реализует sort.Interface для сортировки по нескольким колонкам
type ByColumn struct {
	t       []Track
	columns []string // имена колонок в порядке убывания приоритета
}

// Len возвращает количество элементов
func (b *ByColumn) Len() int { return len(b.t) }

// Swap меняет местами элементы i и j
func (b *ByColumn) Swap(i, j int) { b.t[i], b.t[j] = b.t[j], b.t[i] }

// Less сравнивает элементы i и j по заданным колонкам
func (b *ByColumn) Less(i, j int) bool {
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
				return b.t[i].Year < b.t[j].Year
			}
		case "Length":
			if b.t[i].Length != b.t[j].Length {
				return b.t[i].Length < b.t[j].Length
			}
		}
	}
	return false
}

// SortByColumns - вспомогательная функция для создания сортировщика по колонкам
func SortByColumns(t []Track, columns ...string) {
	sort.Sort(&ByColumn{t, columns})
}

func main() {
	tracks := []Track{
		{"a", "x", "x", 2000, time.Minute},
		{"b", "y", "x", 2000, time.Minute},
		{"a", "z", "x", 2001, time.Minute},
	}

	SortByColumns(tracks, "Title", "Year")
	fmt.Println("After sort by Title, Year: ")
	for _, t := range tracks {
		fmt.Printf("%s %d\n", t.Title, t.Year)
	}

}

package main

import (
	"fmt"
	"gobook/4/github"
	"log"
	"os"
	"time"
)

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d тем: \n", result.TotalCount)

	for _, issue := range result.Items {
		switch {
		case daysAgo(issue.CreatedAt) < 30:
			fmt.Printf("#%-5d %9.9s %.55s Создано менее месяца назад \n", issue.Number, issue.User.Login, issue.Title)
		case daysAgo(issue.CreatedAt) < 365:
			fmt.Printf("#%-5d %9.9s %.55s Создано менее года назад \n", issue.Number, issue.User.Login, issue.Title)
		default:
			fmt.Printf("#%-5d %9.9s %.55s Создано более года назад \n", issue.Number, issue.User.Login, issue.Title)
		}
	}
}

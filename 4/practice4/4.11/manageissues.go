package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type IssueRequest struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

var (
	owner = flag.String("owner", "", "владелец репозитория")
	repo  = flag.String("repo", "", "название репозитория")
	title = flag.String("title", "", "заголовок issue")
	body  = flag.String("body", "", "текст issue")
)

func createIssue(owner, repo, title, body string) error {
	// Получаем токен из окружения
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("не задан GITHUB_TOKEN")
	}

	// Формируем URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
	fmt.Println(url)
	// Подготавливаем тело запроса
	issue := IssueRequest{
		Title: title,
		Body:  body,
	}
	jsonData, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	// Coздаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Отправляем
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем ответ
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка: %s\n%s", resp.Status, bodyBytes)
	}

	fmt.Println("Issue успешно создан!")
	return nil
}

func main() {
	flag.Parse()

	if *owner == "" || *repo == "" || *title == "" {
		fmt.Fprintf(os.Stderr, "необходимо указать -owner, -repo и -title \n")
		os.Exit(1)
	}

	fmt.Printf("Создаю issue в %s/%s с заголовком: %s\n", *owner, *repo, *title)

	err := createIssue(*owner, *repo, *title, *body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
}

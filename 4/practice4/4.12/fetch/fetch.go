package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Comic struct {
	Num         int    `json:"num"`
	Title       string `json:"title"`
	Trancscript string `json:"transcript"`
	Img         string `json:"img"`
	Alt         string `json:"alt"`
}

func fetchComic(num int) (*Comic, error) {
	url := fmt.Sprintf("http://xkcd.com/%d/info.0.json", num)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка %d для комикса %d", resp.StatusCode, num)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var comic Comic
	if err := json.Unmarshal(body, &comic); err != nil {
		return nil, err
	}

	return &comic, nil
}

func main() {
	comics := []Comic{}
	maxNum := 3000

	for num := 1; num <= maxNum; num++ {
		fmt.Printf("Загружаю %d...", num)
		comic, err := fetchComic(num)
		if err != nil {
			fmt.Printf("пропуск: %v\n", err)
			continue
		}
		comics = append(comics, *comic)
		fmt.Printf("OK (%s)\n", comic.Title)

		time.Sleep(200 * time.Microsecond)
	}

	// Сохраняем в JSON
	data, err := json.MarshalIndent(comics, "", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка маршалига: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile("xkcd.json", data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка записи: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Сохранено %d комиксов в xkcd.json\n", len(comics))
}

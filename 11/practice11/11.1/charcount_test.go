package main

import (
	"io"
	"os"
	"testing"
)

func TestCharcount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[rune]int
		invalid  int
	}{
		{
			name:  "ASCII only",
			input: "hello",
			expected: map[rune]int{
				'h': 1,
				'e': 1,
				'l': 2,
				'o': 1,
			},
			invalid: 0,
		},
		{
			name:  "Russian text",
			input: "привет",
			expected: map[rune]int{
				'п': 1,
				'р': 1,
				'и': 1,
				'в': 1,
				'е': 1,
				'т': 1,
			},
			invalid: 0,
		},
		{
			name:     "With invalid UTF-8",
			input:    string([]byte{0xFF, 0xFE, 0xFD}), // некорректные байты
			expected: map[rune]int{},
			invalid:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальный stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Создаем временный файл с входными данными
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			w.Write([]byte(tt.input))
			w.Close()

			// Перенаправляем stdin на наш pipe
			os.Stdin = r

			// Перехватываем stdout
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()

			rOut, wOut, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = wOut

			// Запускаем main
			main()

			// Закрываем stdout и читаем результат
			wOut.Close()
			out, err := io.ReadAll(rOut)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("output: %s", string(out))
		})
	}
}

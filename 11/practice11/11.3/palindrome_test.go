package palindrome

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomPalindrome(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000; i++ {
		s := randomString(rand.Intn(18) + 2)

		if IsPalindrome(s) {
			i--
			continue
		}

		if IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = true, exptected false", s)
		}
	}
}

// randomString генерирует случайную строку длины n
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

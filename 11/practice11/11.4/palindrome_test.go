package palindrome

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomPalindromesWithSpacesAndPunctiation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		s := randomPalindromeWithSpacesAndPunctuation(rng)
		if IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = true, expected false", s)
		}
	}
}

func randomPalindromeWithSpacesAndPunctuation(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		switch rng.Intn(3) {
		case 0:
			r = ' '
		case 1:
			extra := ".,;:!?"
			r = rune(extra[rng.Intn(len(extra))])
		case 2:
		}
		runes[i] = r
		runes[n-1-i] = r
	}
	result := string(runes)

	for i := 0; i < rng.Intn(3); i++ {
		if len(result) > 1 {
			pos := rng.Intn(len(result)-1) + 1
			extra := " ,.;:!?"
			ch := rune(extra[rng.Intn(len(extra))])
			result = result[:pos] + string(ch) + result[pos:]
		}
	}
	return result
}

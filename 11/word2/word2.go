// Пакет word предоставляет утилиты лоя игр со словами.
package word

import "unicode"

// IsPalindrome сообщает, является ли s палиндромом.
// Игнорируем регистар букв и символы, не являющиеся буквами.
func IsPalindrome(s string) bool {
	var letters []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	for i := range letters {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}

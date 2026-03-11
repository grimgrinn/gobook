// basename убирает компоненты каталога и суффикс типа файла.
// a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
package basename

func Basename(s string) string {
	// Отбрасываем последний символ '/' и все перед ним.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	// Созраняем все до последней точки '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

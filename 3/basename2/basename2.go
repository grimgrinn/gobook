package basename2

import "strings"

func Basename(s string) string {
	slash := strings.LastIndex(s, "/") // -1, если не "/" найден
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}

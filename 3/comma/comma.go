// comma вставляет запятые в строковое представление
// неотрицательного десятичного числа.
package comma

import "bytes"

func Comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return Comma(s[:n-3]) + "," + s[n-3:]
}

func CommaSlice(s string) string {
	n := len(s)
	commas := (n - 1) / 3
	totalLen := n + commas

	// создаем срез байт нужной длины
	result := make([]byte, totalLen)

	// указатель на позицию в result (начинаем с конца)
	pos := totalLen - 1

	// счетчик цифр в текущей тройке
	count := 0

	for i := n - 1; i >= 0; i-- {
		result[pos] = s[i]
		pos--
		count++

		// если положили три цифры и это не последняя цифра (i != 0)
		if count == 3 && i != 0 {
			result[pos] = ','
			pos--
			count = 0
		}
	}

	return string(result)
}

func CommaBuff(s string) string {
	var buf bytes.Buffer
	var counter int
	for i := len(s) - 1; i >= 0; i-- {
		buf.WriteByte(s[i])
		counter++
		if counter == 3 {
			buf.WriteByte(',')
			counter = 0
		}

	}

	return reverse(buf.String())
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

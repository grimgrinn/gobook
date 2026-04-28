package tempconv

import (
	"flag"
	"fmt"
	"gobook/2/tempconv"
)

// *celciusFlag соответствует интерфейсу flag.Value.
type celsiusFlag struct{ tempconv.Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // Проверки ошибок не нужны
	switch unit {
	case "C", "°C":
		f.Celsius = tempconv.Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(value))
		return nil
	case "K", "°K":
		f.Celsius = tempconv.KToC(tempconv.Kelvin(value))
		return nil
	}
	return fmt.Errorf("Неверная температура %q", s)
}

// CelsiusFlag определяет флаг Celsius с указанным именем, значением
// по умолчанию и строкой-инструкцией по применению и возвращает адрес
// переменной-флага. Аргумент флага должен содержать числовое значение
// и единицу измерения, например "100C".
func CelsiusFlag(name string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}

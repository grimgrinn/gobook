// Conventer конвертирует разные величины в другие
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Celsius float64
type Fahrenheit float64
type Meter float64
type Foot float64
type Kilogram float64
type Pound float64

func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }
func (m Meter) String() string      { return fmt.Sprintf("%gM", m) }
func (f Foot) String() string       { return fmt.Sprintf("%gFt", f) }
func (k Kilogram) String() string   { return fmt.Sprintf("%gKg", k) }
func (p Pound) String() string      { return fmt.Sprintf("%gPnds", p) }

func main() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			number, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "conventer: %v\n", err)
				os.Exit(1)
			}
			formatter(number)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			number, err := strconv.ParseFloat(line, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "converter: %v\n", err)
				os.Exit(1)
			}
			formatter(number)
		}
	}

}

// CToF преобразует температуру по Цельсию в температуру по Фаренгейту.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC преобразует температуру по Фаренгейту в температуру по Цельсию.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// MToF преобразует метры в футы
func MToF(m Meter) Foot { return Foot(m * 3.28084) }

// FToM преобразует футы в метры
func FToM(f Foot) Meter { return Meter(f / 3.28084) }

// KToP преобразует килограммы в фунты
func KToP(k Kilogram) Foot { return Foot(k * 2.20462) }

// PToK преобразует фунты в килограммы
func PToK(p Pound) Kilogram { return Kilogram(p / 2.20462) }

// formatter форматирует вывод разных величин
func formatter(number float64) {
	c := Celsius(number)
	f := Fahrenheit(number)
	m := Meter(number)
	ft := Foot(number)
	kg := Kilogram(number)
	p := Pound(number)
	FahrenheitFromCelsius := CToF(c)
	CelsiusFromFahrenheit := FToC(f)
	FootFromMeter := MToF(m)
	MeterFromFoot := FToM(ft)
	PoundFromKilogram := KToP(kg)
	KilogramFromPound := PToK(p)
	fmt.Printf("%g is: \n%s\n%s\n%s\n%s\n%s\n%s\n", number, c, f, m, ft, kg, p)
	fmt.Printf("%g Celsius is %g Fahrenheit\n", c, FahrenheitFromCelsius)
	fmt.Printf("%g Fahrenheit is %g Celsious\n", f, CelsiusFromFahrenheit)
	fmt.Printf("%g Meters is %g Foots\n", m, FootFromMeter)
	fmt.Printf("%g Foots is %g Meters\n", ft, MeterFromFoot)
	fmt.Printf("%g Kilograms is %g Pounds\n", kg, PoundFromKilogram)
	fmt.Printf("%g Pounds is %g Kilograms\n", p, KilogramFromPound)
}

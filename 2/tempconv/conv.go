package tempconv

// CToF преобразует температуру по Цельсию в температуру по Фаренгейту.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC преобразует температуру по Фаренгейту в температуру по Цельсию.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// CtoK преобразует темпертаруту по Цельсию в температуру по Кельвину.
func CtoK(c Celsius) Kelvin { return Kelvin(c + 273.15) }

// KtoC преобразует температуру по Келвину в температуру по Цельсию.
func KtoC(k Kelvin) Celsius { return Celsius(k - 273.15) }

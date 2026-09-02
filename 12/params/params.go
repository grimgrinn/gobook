package params

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type fieldInfo struct {
	name         string
	required     bool
	min          *int64
	max          *int64
	pattern      *regexp.Regexp
	isEmail      bool
	isCreditCard bool
}

// Unpack заполняет поля структуры, на которую указывают ptr,
// параметрами из HTTP-запроса в req.
func Unpack(req *http.Request, ptr interface{}) error {
	fmt.Println("1")
	if err := req.ParseForm(); err != nil {
		return err
	}
	// Строит отображение с ключом, являющимся эффективным именем.
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // Структурная переменная

	fieldInfos := make(map[string]fieldInfo)

	for i := 0; i < v.NumField(); i++ {
		mfieldInfo := v.Type().Field(i) // reflect.StructField
		tag := mfieldInfo.Tag.Get("http")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		if name == "" {
			name = strings.ToLower(mfieldInfo.Name)
		}

		info := fieldInfo{
			name: name,
		}

		for _, part := range parts[1:] {
			switch {
			case part == "required":
				info.required = true

			case strings.HasPrefix(part, "min="):
				val, err := strconv.ParseInt(part[4:], 10, 64)
				if err == nil {
					info.min = &val
				}

			case strings.HasPrefix(part, "max="):
				val, err := strconv.ParseInt(part[4:], 10, 64)
				if err == nil {
					info.max = &val
				}

			case strings.HasPrefix(part, "pattern="):
				re, err := regexp.Compile(part[8:])
				if err == nil {
					info.pattern = re
				}

			case part == "email":
				info.isEmail = true

			case part == "creditcard":
				info.isCreditCard = true
			}
		}

		fieldInfos[name] = info
		fields[name] = v.Field(i)
	}

	// Обновляет поля структуры для кадого параметра в запросе.
	for name, values := range req.Form {
		info, ok := fieldInfos[name]
		if !ok {
			continue
		}

		f := fields[name]
		if !f.IsValid() {
			continue
		}

		if info.required && len(values) == 0 {
			return fmt.Errorf("поле %s обязательно", name)
		}

		if f.Kind() == reflect.Slice {
			for _, value := range values {
				if err := populateWithValidation(f, value, info); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}

		} else {
			if len(values) > 0 {
				if err := populateWithValidation(f, values[0], info); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}

	return nil
}

func populateWithValidation(v reflect.Value, value string, info fieldInfo) error {
	if err := validateField(v, value, info); err != nil {
		return err
	}

	return populate(v, value)
}

func validateField(v reflect.Value, value string, info fieldInfo) error {
	if v.Kind() == reflect.String {
		length := int64(len(value))
		if info.min != nil && length < *info.min {
			return fmt.Errorf("минимально длина %d, получено %d", *info.min, length)
		}
		if info.max != nil && length > *info.max {
			return fmt.Errorf("максимальная длиа %d, получено %d", *info.max, length)
		}
	}

	if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("некорректное число: %s", value)
		}
		if info.min != nil && num < *info.min {
			return fmt.Errorf("минимальное значение %d, получено %d", *info.min, num)
		}
		if info.max != nil && num > *info.max {
			return fmt.Errorf("максимальное значение %d, получено %d", *info.max, num)
		}
	}

	if info.pattern != nil && v.Kind() == reflect.String && !info.pattern.MatchString(value) {
		return fmt.Errorf("не соответствует фаормату %s", info.pattern.String())
	}

	if info.isEmail && !isValidEmail(value) {
		return fmt.Errorf("некорректный email: %s", value)
	}

	if info.isCreditCard && !isValidCreditCard(value) {
		return fmt.Errorf("некорректный номер кредитной карты: %s", value)
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`[a-zA_Z0-9._%+-]+@a[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidCreditCard(card string) bool {
	card = strings.ReplaceAll(card, " ", "")
	card = strings.ReplaceAll(card, "-", "")

	if len(card) < 13 || len(card) > 19 {
		return false
	}

	sum := 0
	alternate := false

	for i := len(card) - 1; i >= 0; i-- {
		n := int(card[i] - '0')
		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}
	return sum%10 == 0
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("неподдерживаемый вид %s", v.Type())
	}
	return nil
}

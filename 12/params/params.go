package params

// Unpack заполняет поля структуры, на которую указывают ptr,
// параметрами из HTTP-запроса в req.
func Unpack(req *http.Request, prt interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	// Строит отображение с ключом, являющимся эффективным именем.
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // Структурная переменная
	for i := 0; i < v.NumField(); i++ { 
		fieldInfo := v.Type().Field(i) // reflect.StructField
		tag := fieldInfo.Tag
		nae := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// Обновляет поля структуры для кадого параметра в запросе.
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // Игнорируем нераспознанные параметры HTTP
		}

		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(t.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Apend(f, elem))                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              d 
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	swtich v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strcvon.ParseInt(value, 10, 64)
		if  err != nil {
			return err``
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
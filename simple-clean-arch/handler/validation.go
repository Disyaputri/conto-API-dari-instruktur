package handler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func validateRequest(request any) error {
	value := reflect.ValueOf(request)
	if !value.IsValid() {
		return fmt.Errorf("request cannot be nil")
	}

	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return fmt.Errorf("request cannot be nil")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("request must be a struct")
	}

	return validateStruct(value)
}

func validateStruct(value reflect.Value) error {
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		fieldType := valueType.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		fieldValue := value.Field(i)
		fieldName := strings.Split(fieldType.Tag.Get("json"), ",")[0]
		if fieldName == "" || fieldName == "-" {
			fieldName = fieldType.Name
		}

		rules := strings.Split(fieldType.Tag.Get("validate"), ",")
		for _, rule := range rules {
			if rule == "" {
				continue
			}

			if err := validateRule(fieldName, fieldValue, rule); err != nil {
				return err
			}
		}

		nested := fieldValue
		if nested.Kind() == reflect.Pointer {
			if nested.IsNil() {
				continue
			}
			nested = nested.Elem()
		}
		if nested.Kind() == reflect.Struct {
			if err := validateStruct(nested); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRule(name string, value reflect.Value, rule string) error {
	switch {
	case rule == "required":
		if value.IsZero() {
			return fmt.Errorf("%s is required", name)
		}
	case strings.HasPrefix(rule, "min="):
		minimum, err := strconv.Atoi(strings.TrimPrefix(rule, "min="))
		if err != nil || minimum < 0 {
			return fmt.Errorf("invalid min rule for %s", name)
		}

		switch value.Kind() {
		case reflect.String:
			if len([]rune(value.String())) < minimum {
				return fmt.Errorf("%s must have at least %d characters", name, minimum)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(minimum) {
				return fmt.Errorf("%s must be at least %d", name, minimum)
			}
		default:
			return fmt.Errorf("min rule is unsupported for %s", name)
		}
	default:
		return fmt.Errorf("unsupported validation rule %q for %s", rule, name)
	}

	return nil
}

package genv

import (
	"fmt"
	"reflect"
	"strings"
)

func Marshal(v interface{}) (EnvSet, error) {
	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Ptr || refValue.IsNil() {
		return nil, ErrInvalidValue
	}

	refValue = refValue.Elem()
	if refValue.Kind() != reflect.Struct {
		return nil, ErrInvalidValue
	}

	es := make(EnvSet)
	refType := refValue.Type()

	for i := 0; i < refType.NumField(); i++ {
		valueField := refValue.Field(i)
		switch valueField.Kind() {
		case reflect.Struct:
			if !valueField.Addr().CanInterface() {
				continue
			}

			iface := valueField.Addr().Interface()
			nes, err := Marshal(iface)
			if err != nil {
				return nil, err
			}

			for k, v := range nes {
				es[k] = v
			}
		}

		typeField := refType.Field(i)
		tag := typeField.Tag.Get("env")

		if tag == "" {
			continue
		}

		envKeys := strings.Split(tag, ",")

		var el interface{}
		if typeField.Type.Kind() == reflect.Ptr {
			if valueField.IsNil() {
				continue
			}
			el = valueField.Elem().Interface()
		} else {
			el = valueField.Interface()
		}

		var (
			err      error
			envValue string
		)

		if m, ok := el.(Marshaler); ok {
			envValue, err = m.MarshalEnvironmentValue()
			if err != nil {
				return nil, err
			}
		} else {
			envValue = fmt.Sprintf("%v", el)
		}

		for _, key := range envKeys {
			es[key] = envValue
		}
	}

	return es, nil
}

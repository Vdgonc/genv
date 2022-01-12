package genv

import (
	"os"
	"reflect"
	"strings"
)

func Unmarshal(es EnvSet, v interface{}) error {
	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Ptr || refValue.IsNil() {
		return ErrInvalidValue
	}

	refValue = refValue.Elem()
	if refValue.Kind() != reflect.Struct {
		return ErrInvalidValue
	}

	refType := refValue.Type()
	for i := 0; i < refType.NumField(); i++ {
		valueField := refValue.Field(i)

		switch valueField.Kind() {
		case reflect.Struct:
			if !valueField.Addr().CanInterface() {
				continue
			}

			iface := valueField.Addr().Interface()
			err := Unmarshal(es, iface)
			if err != nil {
				return err
			}
		}

		typeField := refType.Field(i)
		tag := typeField.Tag.Get("env")
		if tag == "" {
			continue
		}

		if !valueField.CanSet() {
			return ErrUnsupportedType
		}

		envTag := parseTag(tag)

		var (
			envValue string
			ok       bool
		)

		for _, key := range envTag.Keys {
			envValue, ok = es[key]

			if ok {
				break
			}
		}

		if !ok {
			if envTag.Default != "" {
				envValue = envTag.Default
			} else if envTag.Required {
				return &ErrMissingRequiredValue{Value: envTag.Keys[0]}
			} else {
				continue
			}

		}

		err := set(valueField.Type(), valueField, envValue)
		if err != nil {
			return err
		}

		delete(es, tag)
	}

	return nil
}

func UnmarshalFromEnviron(v interface{}) (EnvSet, error) {
	es, err := environToEnvSet(os.Environ())
	if err != nil {
		return nil, err
	}

	return es, Unmarshal(es, v)
}

func environToEnvSet(environ []string) (EnvSet, error) {
	env := make(EnvSet)
	for _, e := range environ {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			return nil, ErrInvalidValue
		}
		env[parts[0]] = parts[1]
	}

	return env, nil
}

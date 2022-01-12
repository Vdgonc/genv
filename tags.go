package genv

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type tag struct {
	Keys     []string
	Default  string
	Required bool
}

func parseTag(tagString string) tag {
	var tag tag
	envKeys := strings.Split(tagString, ",")
	for _, key := range envKeys {
		if strings.Contains(key, "=") {
			keyData := strings.SplitN(key, "=", 2)

			switch strings.ToLower(keyData[0]) {
			case "default":
				tag.Default = keyData[1]
			case "required":
				tag.Required = strings.ToLower(keyData[1]) == "true"
			default:
				continue
			}
		} else {
			tag.Keys = append(tag.Keys, key)
		}
	}

	return tag
}

func set(t reflect.Type, f reflect.Value, value string) error {
	var isUnmarshaler bool
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		isUnmarshaler = t.Implements(UnmarshalType) && f.CanInterface()
	} else if f.CanAddr() {
		isUnmarshaler = f.Addr().Type().Implements(UnmarshalType) && f.Addr().CanInterface()
	}

	if isUnmarshaler {
		var ptr reflect.Value
		if isPtr {
			ptr = reflect.New(t.Elem())
		} else {
			ptr = f.Addr()
		}
		if u, ok := ptr.Interface().(Unmarshaler); ok {
			if err := u.UnmarshalEnvironmentValue(value); err != nil {
				return err
			}
			if isPtr {
				f.Set(ptr)
			}
			return nil
		}
	}

	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := set(t.Elem(), ptr.Elem(), value)
		if err != nil {
			return err
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		f.SetFloat(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		f.SetFloat(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.PkgPath() == "time" && t.Name() == "Duration" {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return err
			}

			f.Set(reflect.ValueOf(duration))
			break
		}

		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		f.SetInt(int64(v))
	default:
		return ErrUnsupportedType
	}

	return nil
}

package genv

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidValue    = errors.New("value must be a non-nil pointer to a struct")
	ErrUnsupportedType = errors.New("field is an unsupported type")
	ErrUnexportedField = errors.New("field must be exported")
	ErrInvalidField    = errors.New("field must be a struct, pointer to a struct, or pointer to a pointer to a struct")
	UnmarshalType      = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

type ErrMissingRequiredValue struct {
	Value string
}

func (e ErrMissingRequiredValue) Error() string {
	return fmt.Sprintf("value for this field is required [%s]", e.Value)
}

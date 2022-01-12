package genv

type Unmarshaler interface {
	UnmarshalEnvironmentValue(data string) error
}

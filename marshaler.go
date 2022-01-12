package genv

type Marshaler interface {
	MarshalEnvironmentValue() (string, error)
}

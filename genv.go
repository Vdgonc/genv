package genv

import "log"

type EnvSet map[string]string

func LoadEnvironmentVars(v interface{}) {

	_, err := UnmarshalFromEnviron(v)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

}

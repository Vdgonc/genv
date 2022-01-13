# genv
A  simple Golang library to manage your environment variables using structs.


## How to use


```golang
package main

import (
	"github.com/vdgonc/genv"
)

type Environments struct {
	RabbitMQ struct {
		Host string `env:"RABBITMQ_HOST"`
		Port int    `env:"RABBITMQ_PORT"`
	}
	Redis struct {
		Host string `env:"REDIS_HOST"`
		Port int    `env:"REDIS_PORT"`
	}
	Mongo struct {
		Host string `env:"MONGO_HOST"`
		Port int    `env:"MONGO_PORT"`
	}
	MySQL struct {
		Host string `env:"MYSQL_HOST,required=true"`
		Port int    `env:"MYSQL_PORT,required=true"`
	}
}

func main() {
	environment := Environments{}
	genv.LoadEnvironmentVars(&environment)
}


```

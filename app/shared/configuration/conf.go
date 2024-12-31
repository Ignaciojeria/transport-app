package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewConf)
}

type Conf struct {
	VERSION           string `env:"version,required"`
	PORT              string `env:"PORT" envDefault:"8080"`
	ENVIRONMENT       string `env:"ENVIRONMENT" envDefault:"development"`
	PROJECT_NAME      string `env:"PROJECT_NAME,required"`
	GOOGLE_PROJECT_ID string `env:"GOOGLE_PROJECT_ID"`
}

func NewConf() (Conf, error) {
	return Parse[Conf]()
}

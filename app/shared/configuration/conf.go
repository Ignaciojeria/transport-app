package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewConf)
}

type Conf struct {
	VERSION                  string `env:"version,required"`
	PORT                     string `env:"PORT" envDefault:"8080"`
	ENVIRONMENT              string `env:"ENVIRONMENT" envDefault:"development"`
	PROJECT_NAME             string `env:"PROJECT_NAME" envDefault:"transport-app"`
	GOOGLE_PROJECT_ID        string `env:"GOOGLE_PROJECT_ID" envDefault:"einar-404623"`
	OUTBOX_TOPIC_NAME        string `env:"OUTBOX_TOPIC_NAME" envDefault:"transport-app-events"`
	OPTIMIZATION_STRATEGY    string `env:"OPTIMIZATION_STRATEGY" envDefault:"locationiq"`
	LOCATION_IQ_ACCESS_TOKEN string `env:"LOCATION_IQ_ACCESS_TOKEN"`
	LOCATION_IQ_DNS          string `env:"LOCATION_IQ_DNS"`
	FIREBASE_API_KEY         string `env:"FIREBASE_API_KEY"`
}

func NewConf() (Conf, error) {
	return Parse[Conf]()
}

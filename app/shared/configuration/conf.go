package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewConf)
}

type Conf struct {
	VERSION                             string `env:"version,required"`
	PORT                                string `env:"PORT" envDefault:"8080"`
	ENVIRONMENT                         string `env:"ENVIRONMENT" envDefault:"development"`
	CACHE_STRATEGY                      string `env:"CACHE_STRATEGY" envDefault:"redis"`
	PROJECT_NAME                        string `env:"PROJECT_NAME" envDefault:"transport-app"`
	GOOGLE_PROJECT_ID                   string `env:"GOOGLE_PROJECT_ID" envDefault:"einar-404623"`
	TRANSPORT_APP_TOPIC                 string `env:"TRANSPORT_APP_TOPIC" envDefault:"transport-app-events"`
	OPTIMIZATION_STRATEGY               string `env:"OPTIMIZATION_STRATEGY" envDefault:"locationiq"`
	GEOCODING_STRATEGY                  string `env:"GEOCODING_STRATEGY" envDefault:"locationiq"`
	GOOGLE_MAPS_API_KEY                 string `env:"GOOGLE_MAPS_API_KEY"`
	LOCATION_IQ_ACCESS_TOKEN            string `env:"LOCATION_IQ_ACCESS_TOKEN"`
	LOCATION_IQ_DNS                     string `env:"LOCATION_IQ_DNS"`
	CACHE_URL                           string `env:"CACHE_URL"`
	FIREBASE_API_KEY                    string `env:"FIREBASE_API_KEY"`
	ORDER_SUBMITTED_SUBSCRIPTION        string `env:"ORDER_SUBMITTED_SUBSCRIPTION"`
	TENANT_SUBMITTED_SUBSCRIPTION       string `env:"TENANT_SUBMITTED_SUBSCRIPTION"`
	REGISTRATION_SUBMITTED_SUBSCRIPTION string `env:"REGISTRATION_SUBMITTED_SUBSCRIPTION"`
}

func NewConf() (Conf, error) {
	return Parse[Conf]()
}

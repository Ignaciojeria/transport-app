package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewConf)
}

type Conf struct {
	VERSION                                   string `env:"version,required"`
	PORT                                      string `env:"PORT" envDefault:"8080"`
	ENVIRONMENT                               string `env:"ENVIRONMENT" envDefault:"development"`
	CACHE_STRATEGY                            string `env:"CACHE_STRATEGY" envDefault:"redis"`
	PROJECT_NAME                              string `env:"PROJECT_NAME" envDefault:"transport-app"`
	GOOGLE_PROJECT_ID                         string `env:"GOOGLE_PROJECT_ID" envDefault:"einar-404623"`
	TRANSPORT_APP_TOPIC                       string `env:"TRANSPORT_APP_TOPIC" envDefault:"transport-app-events"`
	OPTIMIZATION_STRATEGY                     string `env:"OPTIMIZATION_STRATEGY" envDefault:"locationiq"`
	GEOCODING_STRATEGY                        string `env:"GEOCODING_STRATEGY" envDefault:"locationiq"`
	GOOGLE_MAPS_API_KEY                       string `env:"GOOGLE_MAPS_API_KEY"`
	LOCATION_IQ_ACCESS_TOKEN                  string `env:"LOCATION_IQ_ACCESS_TOKEN"`
	LOCATION_IQ_DNS                           string `env:"LOCATION_IQ_DNS"`
	CACHE_URL                                 string `env:"CACHE_URL"`
	FIREBASE_API_KEY                          string `env:"FIREBASE_API_KEY"`
	ORDER_SUBMITTED_SUBSCRIPTION              string `env:"ORDER_SUBMITTED_SUBSCRIPTION"`
	NODE_SUBMITTED_SUBSCRIPTION               string `env:"NODE_SUBMITTED_SUBSCRIPTION"`
	TENANT_SUBMITTED_SUBSCRIPTION             string `env:"TENANT_SUBMITTED_SUBSCRIPTION"`
	REGISTRATION_SUBMITTED_SUBSCRIPTION       string `env:"REGISTRATION_SUBMITTED_SUBSCRIPTION"`
	DELIVERIES_SUBMITTED_SUBSCRIPTION         string `env:"DELIVERIES_SUBMITTED_SUBSCRIPTION"`
	ORDER_CANCELLATION_SUBMITTED_SUBSCRIPTION string `env:"ORDER_CANCELLATION_SUBMITTED_SUBSCRIPTION"`
	ROUTE_STARTED_SUBMITTED_SUBSCRIPTION      string `env:"ROUTE_STARTED_SUBMITTED_SUBSCRIPTION"`
	OPTIMIZATION_REQUESTED_SUBSCRIPTION       string `env:"OPTIMIZATION_REQUESTED_SUBSCRIPTION"`
	VROOM_OPTIMIZER_URL                       string `env:"VROOM_OPTIMIZER_URL"`
	VROOM_PLANNER_URL                         string `env:"VROOM_PLANNER_URL"`
	MASTER_NODE_URL                           string `env:"MASTER_NODE_URL"`
	MASTER_NODE_API_KEY                       string `env:"MASTER_NODE_API_KEY"`
	JWT_ISSUER                                string `env:"JWT_ISSUER" envDefault:"transport-app"`
	JWT_PRIVATE_KEY                           string `env:"JWT_PRIVATE_KEY"`
	JWT_PUBLIC_KEY                            string `env:"JWT_PUBLIC_KEY"`
}

func NewConf() (Conf, error) {
	return Parse[Conf]()
}

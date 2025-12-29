package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewConf)
}

type Conf struct {
	VERSION                 string `env:"version,required"`
	PORT                    string `env:"PORT" envDefault:"8080"`
	ENVIRONMENT             string `env:"ENVIRONMENT" envDefault:"development"`
	PROJECT_NAME            string `env:"PROJECT_NAME,required" envDefault:"micartapro"`
	NGROK_AUTHTOKEN         string `env:"NGROK_AUTHTOKEN"`
	GOOGLE_PROJECT_ID       string `env:"GOOGLE_PROJECT_ID" envDefault:"einar-404623"`
	GOOGLE_PROJECT_LOCATION string `env:"GOOGLE_PROJECT_LOCATION" envDefault:"us-central1"`
	SUPABASE_WEBHOOK_SECRET string `env:"SUPABASE_WEBHOOK_SECRET" envDefault:""`
	SUPABASE_JWKS_URL       string `env:"SUPABASE_JWKS_URL" envDefault:"https://rbpdhapfcljecofrscnj.supabase.co/auth/v1/.well-known/jwks.json"`
}

func NewConf() (Conf, error) {
	return Parse[Conf]()
}

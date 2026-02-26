package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewConf)
}

type Conf struct {
	VERSION                      string `env:"version,required"`
	PORT                         string `env:"PORT" envDefault:"8080"`
	ENVIRONMENT                  string `env:"ENVIRONMENT" envDefault:"development"`
	PROJECT_NAME                 string `env:"PROJECT_NAME,required" envDefault:"micartapro"`
	NGROK_AUTHTOKEN              string `env:"NGROK_AUTHTOKEN"`
	GOOGLE_PROJECT_ID            string `env:"GOOGLE_PROJECT_ID" envDefault:"einar-404623"`
	GOOGLE_PROJECT_LOCATION      string `env:"GOOGLE_PROJECT_LOCATION" envDefault:"us-central1"`
	SPEECH_TO_TEXT_PROVIDER      string `env:"SPEECH_TO_TEXT_PROVIDER" envDefault:"chirp"` // "chirp" | "gemini"
	SUPABASE_WEBHOOK_SECRET      string `env:"SUPABASE_WEBHOOK_SECRET" envDefault:""`
	SUPABASE_JWKS_URL            string `env:"SUPABASE_JWKS_URL" envDefault:"https://rbpdhapfcljecofrscnj.supabase.co/auth/v1/.well-known/jwks.json"`
	CREEM_PRODUCT_ID             string `env:"CREEM_PRODUCT_ID" envDefault:""`
	CREEM_API_KEY                string `env:"CREEM_API_KEY" envDefault:""`
	CREEM_DNS                    string `env:"CREEM_DNS" envDefault:""`
	CREEM_SUCCESS_URL            string `env:"CREEM_SUCCESS_URL" envDefault:"https://console.micartapro.com/"`
	CREEM_WEBHOOK_SIGNING_SECRET string `env:"CREEM_WEBHOOK_SIGNING_SECRET" envDefault:""`
	SUPABASE_BACKEND_API_KEY     string `env:"SUPABASE_BACKEND_API_KEY" envDefault:""`
	SUPABASE_PROJECT_URL         string `env:"SUPABASE_PROJECT_URL" envDefault:"https://rbpdhapfcljecofrscnj.supabase.co"`
	MERCADOPAGO_ACCESS_TOKEN     string `env:"MERCADOPAGO_ACCESS_TOKEN" envDefault:""`
	MERCADOPAGO_PUBLIC_KEY       string `env:"MERCADOPAGO_PUBLIC_KEY" envDefault:""`
	MERCADOPAGO_WEBHOOK_SECRET   string `env:"MERCADOPAGO_WEBHOOK_SECRET" envDefault:""`
	MERCADOPAGO_SUCCESS_URL      string `env:"MERCADOPAGO_SUCCESS_URL" envDefault:"https://console.micartapro.com/"`
	MERCADOPAGO_FAILURE_URL      string `env:"MERCADOPAGO_FAILURE_URL" envDefault:"https://console.micartapro.com/"`
	MERCADOPAGO_PENDING_URL      string `env:"MERCADOPAGO_PENDING_URL" envDefault:"https://console.micartapro.com/"`
	MERCADOPAGO_WEBHOOK_URL      string `env:"MERCADOPAGO_WEBHOOK_URL" envDefault:""`
}

func NewConf() (Conf, error) {
	return Parse[Conf]()
}

package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewGeminiConfiguration)
}

type GeminiConfiguration struct {
	GEMINI_API_KEY string `env:"GEMINI_API_KEY,required"`
}

func NewGeminiConfiguration() (GeminiConfiguration, error) {
	return Parse[GeminiConfiguration]()
}

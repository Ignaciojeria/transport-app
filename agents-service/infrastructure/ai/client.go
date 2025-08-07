package ai

import (
	"context"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

func NewClient(conf configuration.Conf) (*genai.Client, error) {
	return genai.NewClient(context.Background(), &genai.ClientConfig{
		Project:  conf.GOOGLE_PROJECT_ID,
		Location: conf.GOOGLE_PROJECT_LOCATION,
		Backend:  genai.BackendVertexAI,
	})
}

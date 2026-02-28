package ai

import (
	"context"
	"fmt"
	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	"google.golang.org/genai"
)

func init() {
	ioc.Register(NewClient)
}

func NewClient(conf configuration.Conf) (*genai.Client, error) {
	if conf.GOOGLE_PROJECT_ID == "" {
		fmt.Println("Vertex AI will not be initialized because GOOGLE_PROJECT_ID is not set")
		return nil, nil
	}
	if conf.GOOGLE_PROJECT_LOCATION == "" {
		fmt.Println("Vertex AI will not be initialized because GOOGLE_PROJECT_LOCATION is not set")
		return nil, nil
	}
	return genai.NewClient(context.Background(), &genai.ClientConfig{
		Project:  conf.GOOGLE_PROJECT_ID,
		Location: conf.GOOGLE_PROJECT_LOCATION,
		Backend:  genai.BackendVertexAI,
	})
}

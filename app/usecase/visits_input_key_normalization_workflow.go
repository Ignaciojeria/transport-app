package usecase

import (
	"context"
	"transport-app/app/adapter/out/agents"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type VisitsInputKeyNormalizationWorkflow func(ctx context.Context, input interface{}) (map[string]string, error)

func init() {
	ioc.Registry(
		NewVisitsInputKeyNormalizationWorkflow,
		agents.NewVisitFieldNamesNormalizer)
}
func NewVisitsInputKeyNormalizationWorkflow(
	visitFieldNamesNormalizer agents.VisitFieldNamesNormalizer) VisitsInputKeyNormalizationWorkflow {
	return func(ctx context.Context, input interface{}) (map[string]string, error) {
		return visitFieldNamesNormalizer(ctx, input)
	}
}

package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func CollectSelectedPaths(ctx context.Context) []string {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	var result []string
	collectFieldPaths(ctx, fields, "", &result)
	return result
}

func collectFieldPaths(ctx context.Context, fields []graphql.CollectedField, prefix string, out *[]string) {
	for _, f := range fields {
		path := f.Name
		if prefix != "" {
			path = prefix + "." + f.Name
		}
		*out = append(*out, path)

		subFields := graphql.CollectFields(graphql.GetOperationContext(ctx), f.SelectionSet, nil)
		collectFieldPaths(ctx, subFields, path, out)
	}
}

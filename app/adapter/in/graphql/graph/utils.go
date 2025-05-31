package graph

import (
	"context"
	"strings"

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

// ConvertSelectedPathsToMap convierte los paths seleccionados a un mapa donde la clave es el path y el valor es true
// Elimina el prefijo "edges.node" de los paths si existe
func ConvertSelectedPathsToMap(ctx context.Context) map[string]any {
	paths := CollectSelectedPaths(ctx)
	result := make(map[string]any, len(paths))
	for _, path := range paths {
		// Eliminar el prefijo "edges.node" si existe
		cleanPath := strings.TrimPrefix(path, "edges.node.")
		result[cleanPath] = true
	}
	return result
}

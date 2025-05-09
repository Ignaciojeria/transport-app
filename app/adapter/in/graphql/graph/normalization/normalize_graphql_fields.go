package normalization

import "strings"

func normalizeGraphQLField(graphqlField string) string {
	// Elimina el prefijo "edges.node." que viene de la consulta GraphQL
	return strings.TrimPrefix(graphqlField, "edges.node.")
}

func NormalizeGraphQLFields(graphqlFields []string) map[string]struct{} {
	normalized := make(map[string]struct{}, len(graphqlFields))
	for _, field := range graphqlFields {
		normalized[normalizeGraphQLField(field)] = struct{}{}
	}
	return normalized
}

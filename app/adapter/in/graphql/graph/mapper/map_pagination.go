package mapper

import (
	"encoding/base64"
	"strconv"
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/domain"
)

func MapOrderPagination(p *model.OrderPagination) domain.Pagination {
	limit := 10 // default
	if p != nil && p.First != nil {
		limit = *p.First
	}

	offset := 0
	if p != nil && p.After != nil {
		offset = decodeCursor(*p.After)
	}

	return domain.Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

// decodeCursor convierte el cursor Relay-style (base64) a offset entero
func decodeCursor(cursor string) int {
	if cursor == "" {
		return 0
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0
	}

	offset, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0
	}

	return offset
}

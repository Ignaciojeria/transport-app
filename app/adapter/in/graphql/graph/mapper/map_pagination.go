package mapper

import (
	"encoding/base64"
	"strconv"
	"transport-app/app/domain"
)

func MapRelayPagination(first, last *int, after, before *string) domain.Pagination {
	return domain.Pagination{
		First:  first,
		Last:   last,
		After:  after,
		Before: before,
	}
}

func DecodeCursor(cursor string) int {
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

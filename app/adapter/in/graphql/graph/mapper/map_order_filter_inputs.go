package mapper

import (
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/domain"
)

func MapOrderFilterWithPagination(
	filter *model.OrderFilterInput,
	pagination domain.Pagination,
	requestedFields []string,
) domain.OrderFilterInput {
	if filter == nil {
		filter = &model.OrderFilterInput{}
	}
	return domain.OrderFilterInput{
		Pagination:           pagination,
		ReferenceIds:         filter.ReferenceIds,
		ReferenceType:        filter.ReferenceType,
		ReferenceValue:       filter.ReferenceValue,
		Lpns:                 filter.Lpns,
		GroupBy:              filter.GroupBy,
		LabelType:            filter.LabelType,
		LabelValue:           filter.LabelValue,
		Commerces:            filter.Commerces,
		Consumers:            filter.Consumers,
		OriginNodeReferences: filter.OriginNodeReferences,
		RequestedFields:      requestedFields,
	}
}

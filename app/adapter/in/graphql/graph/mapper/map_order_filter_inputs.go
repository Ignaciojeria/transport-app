package mapper

import (
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/domain"
)

func MapOrderFilterWithPagination(filter *model.OrderFilterInput, pagination *model.OrderPagination) domain.OrderFilterInput {
	if filter == nil {
		filter = &model.OrderFilterInput{}
	}
	domainFilter := domain.OrderFilterInput{
		Pagination:     MapOrderPagination(pagination),
		ReferenceIds:   filter.ReferenceIds,
		ReferenceType:  filter.ReferenceType,
		ReferenceValue: filter.ReferenceValue,
		Lpns:           filter.Lpns,
		GroupBy:        filter.GroupBy,
		LabelType:      filter.LabelType,
		LabelValue:     filter.LabelValue,
		Commerces:      filter.Commerces,
		Consumers:      filter.Consumers,
	}
	return domainFilter
}

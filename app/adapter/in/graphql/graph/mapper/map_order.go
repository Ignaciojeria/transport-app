package mapper

import (
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/domain"
)

func MapOrder(o domain.Order) *model.Order {
	return &model.Order{
		ReferenceID: o.ReferenceID.String(),
	}
}

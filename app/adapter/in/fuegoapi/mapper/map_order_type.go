package mapper

import "transport-app/app/domain"

func MapOrderTypeToDomain(orderType struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}) domain.OrderType {
	return domain.OrderType{
		Type:        orderType.Type,
		Description: orderType.Description,
	}
}

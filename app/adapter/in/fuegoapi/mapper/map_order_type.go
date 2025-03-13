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

func MapOrderTypeFromDomain(orderType domain.OrderType) struct {
	Description string `json:"description"`
	Type        string `json:"type"`
} {
	return struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	}{
		Description: orderType.Description,
		Type:        orderType.Type,
	}
}

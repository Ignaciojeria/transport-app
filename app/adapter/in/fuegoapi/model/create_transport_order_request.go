package model

import "transport-app/app/domain"

type CreateTransportOrderRequest struct {
	ReferenceID             domain.ReferenceID             `json:"referenceID"`
	OrderType               domain.OrderType               `json:"orderType"`
	References              []domain.References            `json:"references"`
	Origin                  domain.Origin                  `json:"origin"`
	Destination             domain.Destination             `json:"destination"`
	Items                   []domain.Items                 `json:"items"`
	Packages                []domain.Packages              `json:"packages"`
	CollectAvailabilityDate domain.CollectAvailabilityDate `json:"collectAvailabilityDate"`
	PromisedDate            domain.PromisedDate            `json:"promisedDate"`
	Visit                   domain.Visit                   `json:"visit"`
	TransportRequirements   []domain.References            `json:"transportRequirements"`
}

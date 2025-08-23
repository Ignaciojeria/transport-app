package request

import (
	"context"
	"transport-app/app/domain"
)

type RouteStartedRequest struct {
	StartedAt string `json:"startedAt" example:"2025-06-06T14:30:00Z"`
	Carrier   struct {
		Name       string `json:"name" example:"Transportes ABC"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"carrier"`
	Driver struct {
		Email      string `json:"email" example:"juan@example.com"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"driver"`
	Vehicle struct {
		Plate string `json:"plate" example:"ABC123"`
	} `json:"vehicle"`
	Route struct {
		ReferenceID string `json:"referenceID"`
	} `json:"route"`
}

func (r RouteStartedRequest) Map(ctx context.Context) domain.Route {
	route := domain.Route{
		ReferenceID: r.Route.ReferenceID,
		Vehicle: domain.Vehicle{
			Plate: r.Vehicle.Plate,
			Carrier: domain.Carrier{
				Name:       r.Carrier.Name,
				NationalID: r.Carrier.NationalID,
				Driver: domain.Driver{
					Email:      r.Driver.Email,
					NationalID: r.Driver.NationalID,
				},
			},
		},
	}
	return route
}

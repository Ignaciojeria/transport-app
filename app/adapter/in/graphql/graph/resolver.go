package graph

import (
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func init() {
	ioc.Registry(
		NewResolver,
		usecase.NewSearchOrders)
}

type Resolver struct {
	usecase.SearchOrders
}

func NewResolver(
	searchOrders usecase.SearchOrders) *Resolver {
	return &Resolver{
		SearchOrders: searchOrders,
	}
}

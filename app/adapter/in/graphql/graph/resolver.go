package graph

import (
	"transport-app/app/adapter/out/tidbrepository"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func init() {
	ioc.Registry(
		NewResolver,
		tidbrepository.NewFindDeliveryUnitsProjectionResult)
}

type Resolver struct {
	findDeliveryUnitsProjectionResult tidbrepository.FindDeliveryUnitsProjectionResult
}

func NewResolver(
	findDeliveryUnitsProjectionResult tidbrepository.FindDeliveryUnitsProjectionResult) *Resolver {
	return &Resolver{
		findDeliveryUnitsProjectionResult: findDeliveryUnitsProjectionResult,
	}
}

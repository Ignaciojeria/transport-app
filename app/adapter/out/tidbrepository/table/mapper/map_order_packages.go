package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderPackages(ctx context.Context, order domain.Order) []table.OrderPackage {
	var orderPackages []table.OrderPackage
	for _, p := range order.Packages {
		orderPackages = append(orderPackages, table.OrderPackage{
			OrderDoc:   order.DocID(ctx).String(),
			PackageDoc: p.DocID(ctx, order.ReferenceID.String()).String(),
		})
	}
	return orderPackages
}

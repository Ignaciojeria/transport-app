package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/mapper"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewFindOrdersByFilters,
		tidb.NewTIDBConnection)
}

type FindOrdersByFilters func(
	context.Context,
	domain.OrderSearchFilters) ([]domain.Order, error)

func NewFindOrdersByFilters(
	conn tidb.TIDBConnection,
) FindOrdersByFilters {
	return func(ctx context.Context, filters domain.OrderSearchFilters) ([]domain.Order, error) {
		db := conn.WithContext(ctx)
		var organizationCountries []table.OrganizationCountry

		// Paso 1: Buscar organizationCountry
		if err := db.Table("organization_countries").
			Joins("JOIN api_keys ON api_keys.organization_id = organization_countries.organization_id").
			Where("api_keys.key = ? AND organization_countries.country = ?", filters.Organization.Key, filters.Organization.Country.Alpha2()).
			Find(&organizationCountries).Error; err != nil {
			return nil, err
		}

		if len(organizationCountries) == 0 {
			return nil, nil // No hay resultados
		}

		organizationCountryIDs := make([]int64, len(organizationCountries))
		for i, org := range organizationCountries {
			organizationCountryIDs[i] = org.ID
		}

		// Preparar listas para consolidar las órdenes encontradas
		var ordersByReferenceID []table.Order
		var ordersByPackages []table.Packages

		// Paso 2: Filtrar por referenceIds si están presentes
		if len(filters.ReferenceIDs) > 0 {
			query := db.Model(&table.Order{}).
				Preload("Commerce").
				Preload("Consumer").
				Preload("OriginContact").
				Preload("DestinationContact").
				Preload("OriginAddressInfo").
				Preload("OrderStatus").
				Preload("OrganizationCountry").
				Preload("DestinationAddressInfo").
				Preload("OriginNodeInfo").
				Preload("DestinationNodeInfo").
				Preload("Packages").
				Where("reference_id IN (?) AND organization_country_id IN (?)", filters.ReferenceIDs, organizationCountryIDs)

			// Agregar filtro de commerces si está presente
			if len(filters.Commerces) > 0 {
				query = query.Where("commerce_id IN (?)", filters.Commerces)
			}

			if err := query.Find(&ordersByReferenceID).Error; err != nil {
				return nil, err
			}
		}

		// Paso 3: Filtrar por packageLpns si están presentes
		if len(filters.Packages) > 0 {
			lpns := make([]string, len(filters.Packages))
			for i, pkg := range filters.Packages {
				lpns[i] = pkg.Lpn
			}

			query := db.Model(&table.Packages{}).
				Joins("JOIN orders ON packages.order_id = orders.id").
				Where("packages.lpn IN (?) AND orders.organization_country_id IN (?)", lpns, organizationCountryIDs).
				Preload("Order").
				Preload("Order.OrganizationCountry").
				Preload("Order.Commerce").
				Preload("Order.Consumer").
				Preload("Order.OrderStatus").
				Preload("Order.OriginContact").
				Preload("Order.DestinationContact").
				Preload("Order.OriginAddressInfo").
				Preload("Order.DestinationAddressInfo").
				Preload("Order.OriginNodeInfo").
				Preload("Order.DestinationNodeInfo").
				Preload("Order.Packages").
				Preload("Order.Items")

			// Agregar filtro por commerces si está presente
			if len(filters.Commerces) > 0 {
				query = query.Where("orders.commerce_id IN (?)", filters.Commerces)
			}

			if err := query.Find(&ordersByPackages).Error; err != nil {
				return nil, err
			}
		}

		// Consolidar resultados y eliminar duplicados
		orderMap := make(map[int64]table.Order)
		for _, order := range ordersByReferenceID {
			orderMap[order.ID] = order
		}

		for _, pkg := range ordersByPackages {
			if pkg.Order.ID != 0 {
				orderMap[pkg.Order.ID] = pkg.Order
			}
		}

		consolidatedOrders := make([]table.Order, 0, len(orderMap))
		for _, order := range orderMap {
			consolidatedOrders = append(consolidatedOrders, order)
		}

		// Aplicar paginación
		start := (filters.Pagination.Page - 1) * filters.Pagination.Size
		end := start + filters.Pagination.Size
		if start > len(consolidatedOrders) {
			start = len(consolidatedOrders)
		}
		if end > len(consolidatedOrders) {
			end = len(consolidatedOrders)
		}

		paginatedOrders := consolidatedOrders[start:end]

		domainOrders := make([]domain.Order, len(paginatedOrders))
		for i, order := range paginatedOrders {
			domainOrders[i] = mapper.MapOrderDomain(order)
		}

		return domainOrders, nil
	}
}

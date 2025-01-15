package tidbrepository

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewSaveOrder,
		tidb.NewTIDBConnection,
		NewLoadOrderStatuses)
}

type SaveOrder func(
	ctx context.Context,
	existingOrder, orderToCreate domain.Order) (domain.Order, error)

func NewSaveOrder(
	conn tidb.TIDBConnection,
	loadOrderSorderStatuses LoadOrderStatuses,
) SaveOrder {
	return func(ctx context.Context, existingOrder, orderToCreate domain.Order) (domain.Order, error) {
		available := loadOrderSorderStatuses().Available()
		type QueryResult struct {
			OrganizationCountryID int64
			CommerceID            int64
			ConsumerID            int64
			OrderTypeID           int64
			OriginContactID       int64
			DestinationContactID  int64
			OriginAddressID       int64
			DestinationAddressID  int64
			OriginNodeInfoID      int64
			DestinationNodeInfoID int64
		}

		var result QueryResult
		err := conn.Raw(`
		SELECT 
		  c.id AS commerce_id,
		  con.id AS consumer_id,
		  ot.id AS order_type_id,
		  o_ct.id AS origin_contact_id,
		  d_ct.id AS destination_contact_id,
		  o_ai.id AS origin_address_id,
		  d_ai.id AS destination_address_id,
		  o_ni.id AS origin_node_info_id,
		  d_ni.id AS destination_node_info_id
		FROM 
		  organization_countries org
		  LEFT JOIN commerces c 
			ON c.organization_country_id = org.id
			AND c.name = ?
		  LEFT JOIN consumers con 
			ON con.organization_country_id = org.id
			AND con.name = ?
		  LEFT JOIN order_types ot 
			ON ot.organization_country_id = org.id
			AND ot.type = ?
		  LEFT JOIN contacts o_ct 
			ON o_ct.organization_country_id = org.id
			AND o_ct.full_name = ? 
			AND o_ct.email = ? 
			AND o_ct.phone = ?
		  LEFT JOIN contacts d_ct 
			ON d_ct.organization_country_id = org.id
			AND d_ct.full_name = ? 
			AND d_ct.email = ? 
			AND d_ct.phone = ?
		  LEFT JOIN address_infos o_ai 
			ON o_ai.organization_country_id = org.id
			AND o_ai.raw_address = ?
		  LEFT JOIN address_infos d_ai 
			ON d_ai.organization_country_id = org.id
			AND d_ai.raw_address = ?
		  LEFT JOIN node_infos o_ni 
			ON o_ni.organization_country_id = org.id
			AND o_ni.reference_id = ?
		  LEFT JOIN node_infos d_ni 
			ON d_ni.organization_country_id = org.id
			AND d_ni.reference_id = ?
		WHERE 
		  org.id = ?;
		`,
			// Argumentos en el mismo orden que los placeholders en la query
			orderToCreate.BusinessIdentifiers.Commerce,        // Nombre del commerce
			orderToCreate.BusinessIdentifiers.Consumer,        // Nombre del consumer
			orderToCreate.OrderType.Type,                      // Tipo de orden
			orderToCreate.Origin.AddressInfo.Contact.FullName, // Contacto origen
			orderToCreate.Origin.AddressInfo.Contact.Email,
			orderToCreate.Origin.AddressInfo.Contact.Phone,
			orderToCreate.Destination.AddressInfo.Contact.FullName, // Contacto destino
			orderToCreate.Destination.AddressInfo.Contact.Email,
			orderToCreate.Destination.AddressInfo.Contact.Phone,
			orderToCreate.Origin.AddressInfo.RawAddress(),      // Dirección origen
			orderToCreate.Destination.AddressInfo.RawAddress(), // Dirección destino
			orderToCreate.Origin.NodeInfo.ReferenceID,          // Nodo origen
			orderToCreate.Destination.NodeInfo.ReferenceID,     // Nodo destino
			orderToCreate.Organization.OrganizationCountryID,   // ID de la organización
		).Scan(&result).Error

		if err != nil {
			return domain.Order{}, err
		}
		orderTable := mapper.MapOrderToTable(orderToCreate)
		orderTable.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
		orderTable.CommerceID = result.CommerceID
		orderTable.ConsumerID = result.ConsumerID
		orderTable.OriginContactID = result.OriginContactID
		orderTable.DestinationContactID = result.DestinationContactID
		orderTable.OriginAddressInfoID = result.OriginAddressID
		orderTable.DestinationAddressInfoID = result.DestinationAddressID
		orderTable.OriginNodeInfoID = result.OriginNodeInfoID
		orderTable.DestinationNodeInfoID = result.DestinationNodeInfoID
		orderTable.OrderTypeID = result.OrderTypeID
		orderTable.OrderStatusID = available.ID

		return domain.Order{}, conn.Transaction(func(tx *gorm.DB) error {
			// Guardar entidades que no existen y actualizar relaciones en orderTable
			if result.CommerceID == 0 {
				orderTable.Commerce.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
				if err := tx.Save(&orderTable.Commerce).Error; err != nil {
					return err
				}
				orderTable.CommerceID = orderTable.Commerce.ID
			}

			if result.ConsumerID == 0 {
				orderTable.Consumer.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
				if err := tx.Save(&orderTable.Consumer).Error; err != nil {
					return err
				}
				orderTable.ConsumerID = orderTable.Consumer.ID
			}

			if result.OriginContactID == 0 {
				orderTable.OriginContact.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
				if err := tx.Save(&orderTable.OriginContact).Error; err != nil {
					return err
				}
				orderTable.OriginContactID = orderTable.OriginContact.ID
			}

			if result.DestinationContactID == 0 {
				if orderToCreate.IsOriginAndDestinationContactEqual() {
					orderTable.DestinationContactID = orderTable.OriginContactID
				} else {
					orderTable.DestinationContact.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
					if err := tx.Save(&orderTable.DestinationContact).Error; err != nil {
						return err
					}
					orderTable.DestinationContactID = orderTable.DestinationContact.ID
				}
			}

			if result.OriginAddressID == 0 {
				orderTable.OriginAddressInfo.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
				if err := tx.Save(&orderTable.OriginAddressInfo).Error; err != nil {
					return err
				}
				orderTable.OriginAddressInfoID = orderTable.OriginAddressInfo.ID
			}

			if result.DestinationAddressID == 0 {
				if orderToCreate.IsOriginAndDestinationAddressEqual() {
					orderTable.DestinationAddressInfoID = orderTable.OriginAddressInfoID
				} else {
					orderTable.DestinationAddressInfo.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
					if err := tx.Save(&orderTable.DestinationAddressInfo).Error; err != nil {
						return err
					}
					orderTable.DestinationAddressInfoID = orderTable.DestinationAddressInfo.ID
				}
			}

			if result.OriginNodeInfoID == 0 {
				orderTable.OriginNodeInfo.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
				orderTable.OriginNodeInfo.AddressID = orderTable.OriginAddressInfoID
				if err := tx.Save(&orderTable.OriginNodeInfo).Error; err != nil {
					return err
				}
				orderTable.OriginNodeInfoID = orderTable.OriginNodeInfo.ID
			}

			if result.DestinationNodeInfoID == 0 {
				if orderToCreate.IsOriginAndDestinationNodeEqual() {
					orderTable.DestinationNodeInfoID = orderTable.OriginNodeInfoID
				} else {
					orderTable.DestinationNodeInfo.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
					orderTable.DestinationNodeInfo.AddressID = orderTable.DestinationAddressInfoID
					if err := tx.Save(&orderTable.DestinationNodeInfo).Error; err != nil {
						return err
					}
					orderTable.DestinationNodeInfoID = orderTable.DestinationNodeInfo.ID
				}
			}

			if result.OrderTypeID == 0 {
				orderTable.OrderType.OrganizationCountryID = orderToCreate.Organization.OrganizationCountryID
				if err := tx.Create(&orderTable.OrderType).Error; err != nil {
					return err
				}
				orderTable.OrderTypeID = orderTable.OrderType.ID
			}

			orderTable.Commerce = table.Commerce{}
			orderTable.Consumer = table.Consumer{}
			orderTable.OriginContact = table.Contact{}
			orderTable.DestinationContact = table.Contact{}
			orderTable.OriginAddressInfo = table.AddressInfo{}
			orderTable.DestinationAddressInfo = table.AddressInfo{}
			orderTable.OriginNodeInfo = table.NodeInfo{}
			orderTable.DestinationNodeInfo = table.NodeInfo{}
			orderTable.OrderType = table.OrderType{}

			if err := tx.Create(&orderTable).Error; err != nil {
				return err
			}

			if err := saveOrderPackages(tx, orderTable.ID, orderToCreate.Organization.OrganizationCountryID, orderToCreate.Packages); err != nil {
				return fmt.Errorf("failed to save packages: %w", err)
			}

			return nil
		})
	}
}

// Guardar paquetes asociados a la orden
func saveOrderPackages(tx *gorm.DB, orderID int64, organizationCountry int64, incomingPackages []domain.Package) error {
	// Procesar los paquetes entrantes
	for _, incomingPkg := range incomingPackages {
		// Usar createOrUpdatePackage para manejar cada paquete
		if err := createOrUpdatePackage(tx, orderID, organizationCountry, incomingPkg); err != nil {
			return fmt.Errorf("failed to process package LPN %s: %w", incomingPkg.Lpn, err)
		}
	}

	return nil
}

// Buscar paquetes existentes desde la base de datos
func fetchExistingPackages(tx *gorm.DB, incomingPackages []domain.Package) (map[string]table.Package, error) {
	lpns := extractLPNs(incomingPackages)

	var existingPackages []table.Package
	if err := tx.Where("lpn IN (?)", lpns).Find(&existingPackages).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch existing packages: %w", err)
	}

	// Convertir la lista en un mapa para acceso rápido
	existingPackagesMap := make(map[string]table.Package)
	for _, pkg := range existingPackages {
		existingPackagesMap[pkg.Lpn] = pkg
	}

	return existingPackagesMap, nil
}

// Extraer LPNs de los paquetes entrantes
func extractLPNs(packages []domain.Package) []string {
	lpns := make([]string, len(packages))
	for i, pkg := range packages {
		lpns[i] = pkg.Lpn
	}
	return lpns
}

func createOrUpdatePackage(tx *gorm.DB, orderID int64, organizationCountry int64, incomingPkg domain.Package) error {
	var existingPkg table.Package

	// Intentar obtener el paquete existente por LPN
	if err := tx.Where("lpn = ?", incomingPkg.Lpn).First(&existingPkg).Error; err == nil {
		// Mapear el paquete existente al dominio
		domainExistingPkg := mapper.MapPackageDomain(existingPkg)

		// Verificar si necesita actualización
		updatedPkg, needsUpdate := domainExistingPkg.UpdateIfChanged(incomingPkg)
		if needsUpdate {
			// Mapear el paquete actualizado a la tabla
			updatedTablePkg := mapper.MapPackageToTable(updatedPkg)

			// Actualizar el paquete existente
			if err := tx.Updates(&updatedTablePkg).Error; err != nil {
				return fmt.Errorf("failed to update package LPN %s: %w", incomingPkg.Lpn, err)
			}
		}

		// Verificar si el paquete ya está asociado con la misma orden
		var existingOrderPkg table.OrderPackage
		if err := tx.Where("order_id = ? AND package_id = ?", orderID, existingPkg.ID).First(&existingOrderPkg).Error; err == nil {
			// Si ya está asociado con la misma orden, retornar sin error
			return nil
		}

		// Asociar el paquete existente (actualizado o no) con la orden
		orderPkg := table.OrderPackage{
			OrderID:   orderID,
			PackageID: existingPkg.ID,
		}
		if err := tx.Create(&orderPkg).Error; err != nil {
			return fmt.Errorf("failed to associate existing package LPN %s with order: %w", incomingPkg.Lpn, err)
		}
	} else if err == gorm.ErrRecordNotFound {
		// Si el paquete no existe, crearlo
		newPkgTable := mapper.MapPackageToTable(incomingPkg)
		newPkgTable.OrganizationCountryID = organizationCountry
		// Crear el nuevo paquete
		if err := tx.Create(&newPkgTable).Error; err != nil {
			return fmt.Errorf("failed to create new package LPN %s: %w", incomingPkg.Lpn, err)
		}

		// Asociar el nuevo paquete con la orden
		orderPkg := table.OrderPackage{
			OrderID:   orderID,
			PackageID: newPkgTable.ID,
		}
		if err := tx.Create(&orderPkg).Error; err != nil {
			return fmt.Errorf("failed to create order-package relation for LPN %s: %w", incomingPkg.Lpn, err)
		}
	} else {
		// Error inesperado al intentar obtener el paquete
		return fmt.Errorf("failed to query package LPN %s: %w", incomingPkg.Lpn, err)
	}

	return nil
}

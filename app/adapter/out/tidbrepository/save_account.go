package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewSaveAccount, tidb.NewTIDBConnection)
}

type SaveAccount func(
	context.Context,
	domain.Operator) (domain.Operator, error)

func NewSaveAccount(conn tidb.TIDBConnection) SaveAccount {
	type SaveAccountQuery struct {
		OrganizationCountryID int64
		ContactID             int64
		AddressInfoID         int64
		NodeInfoID            int64
	}

	return func(ctx context.Context, o domain.Operator) (domain.Operator, error) {
		/*
			var query SaveAccountQuery

			// Ejecutar la consulta con LEFT JOINs
			err := conn.DB.Table("organization_countries").
				Select(`
					organization_countries.id AS organization_country_id,
					contacts.id AS contact_id,
					address_infos.id AS address_info_id,
					node_infos.id AS node_info_id
				`).
				Joins(`
					JOIN api_keys
						ON api_keys.organization_id = organization_countries.organization_id
					LEFT JOIN contacts
						ON contacts.organization_country_id = organization_countries.id
						AND contacts.email = ?
					LEFT JOIN address_infos
						ON address_infos.organization_country_id = organization_countries.id
						AND address_infos.raw_address = ?
					LEFT JOIN node_infos
						ON node_infos.organization_country_id = organization_countries.id
						AND node_infos.reference_id = ?
				`,
					o.Contact.Email,
					o.Origin.AddressInfo.RawAddress(),
					o.Origin.ReferenceID,
				).
				Where("organization_countries.country = ? AND api_keys.key = ?",
					o.Organization.Country.Alpha2(),
					o.Organization.Key,
				).
				Scan(&query).Error

			if err != nil {
				return domain.Account{}, err
			}

			contactTable := mapper.MapContactToTable(o.Contact, query.OrganizationCountryID)
			addressInfoTable := mapper.MapAddressInfoTable(o.Origin.AddressInfo, query.OrganizationCountryID)

			if err := conn.Transaction(func(tx *gorm.DB) error {
				// Crear o actualizar Contact
				if query.ContactID == 0 {
					if err := tx.Create(&contactTable).Error; err != nil {
						return err
					}
					query.ContactID = contactTable.ID // Asignar el ID generado
				}

				// Crear o actualizar AddressInfo
				if query.AddressInfoID == 0 {
					if err := tx.Create(&addressInfoTable).Error; err != nil {
						return err
					}
					query.AddressInfoID = addressInfoTable.ID // Asignar el ID generado
				}

				// Crear o actualizar NodeInfo
				nodeInfoTable := mapper.MapNodeInfoTable(o.Origin)
				if query.NodeInfoID == 0 {
					if err := tx.Create(&nodeInfoTable).Error; err != nil {
						return err
					}
					query.NodeInfoID = nodeInfoTable.ID // Asignar el ID generado
				}
				// Crear o actualizar Account
				accountTable := mapper.MapAccountTable(
					o,
					query.NodeInfoID,
					query.ContactID,
					query.OrganizationCountryID)
				if err := tx.Save(&accountTable).Error; err != nil {
					return err
				}

				return nil
			}); err != nil {
				return domain.Account{}, err
			}

			// Devolver los datos actualizados
			return domain.Account{
				Organization: o.Organization,
				Origin:       o.Origin,
				Contact:      o.Contact,
			}, nil*/
		return domain.Operator{}, nil
	}
}

package tidbrepository

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(NewSaveAccount, tidb.NewTIDBConnection)
}

type SaveAccount func(
	context.Context,
	domain.Account) (domain.Account, error)

func NewSaveAccount(conn tidb.TIDBConnection) SaveAccount {
	type SaveAccountQuery struct {
		OrganizationCountryID int64
		ContactID             int64
		AddressInfoID         int64
		NodeInfoID            int64
	}

	return func(ctx context.Context, o domain.Account) (domain.Account, error) {
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
				o.Origin.NodeInfo.ReferenceID,
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
		//addressInfoTable := mapper.MapAddressInfoTable(o.Origin.AddressInfo)
		//nodeInfoTable := mapper.MapNodeInfoTable(o.Origin.NodeInfo)

		if err := conn.Transaction(func(tx *gorm.DB) error {
			if query.ContactID == 0 {
				if err := tx.Save(&contactTable).Error; err != nil {
					fmt.Println(err.Error())
					return err
				}
			}
			return nil
		}); err != nil {
			return domain.Account{}, err
		}

		// Devuelve el resultado o úsalo según tus necesidades
		return domain.Account{
			Organization: o.Organization,
			Origin: domain.Origin{
				NodeInfo: o.Origin.NodeInfo,
				AddressInfo: domain.AddressInfo{
					Contact: o.Contact,
				},
			},
			Contact: o.Contact,
		}, nil
	}
}

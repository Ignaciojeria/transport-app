package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchOperatorByEmail func(context.Context, domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(NewSearchOperatorByEmail, tidb.NewTIDBConnection)
}
func NewSearchOperatorByEmail(conn tidb.TIDBConnection) SearchOperatorByEmail {
	return func(ctx context.Context, o domain.Operator) (domain.Operator, error) {
		var account table.Account
		if err := conn.DB.WithContext(ctx).
			Preload("OriginNodeInfo").
			Preload("OriginNodeInfo.AddressInfo").
			Preload("OriginNodeInfo.NodeType").
			Preload("Contact").
			Joins("JOIN organization_countries oc ON accounts.organization_country_id = oc.id").
			Joins("JOIN organizations org ON oc.organization_id = org.id").
			Joins("JOIN api_keys ak ON org.id = ak.organization_id").
			Joins("JOIN contacts c ON accounts.contact_id = c.id"). // Unión con la tabla Contact
			Where("ak.key = ? AND oc.country = ?", o.Organization.Key, o.Organization.Country.Alpha2()).
			Where("c.email = ? AND accounts.type = ?", o.Contact.Email, "operator"). // Filtro por email del Contact y type del Account
			First(&account).Error; err != nil {
			return domain.Operator{}, err
		}
		return account.MapOperator(), nil
	}
}

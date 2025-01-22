package tidbrepository

import (
	"context"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNode func(context.Context, domain.NodeInfo) error

func init() {
	ioc.Registry(NewUpsertNode, tidb.NewTIDBConnection)
}
func NewUpsertNode(conn tidb.TIDBConnection) UpsertNode {
	return func(ctx context.Context, o domain.NodeInfo) error {
		addressInfo := mapper.MapAddressInfoTable(o.AddressInfo, o.Organization.OrganizationCountryID)
		contact := mapper.MapContactToTable(o.Operator.Contact, o.Organization.OrganizationCountryID)
		nodeInfo := mapper.MapNodeInfoTable(o, o.Organization.OrganizationCountryID, addressInfo.ID)
		nodeInfo.OrganizationCountryID = o.Organization.OrganizationCountryID
		//nodeInfo.OperatorID = o.NodeInfo.Operator.ID
		nodeInfo.Operator.ID = o.Operator.ID
		nodeInfo.Operator.Type = o.Operator.Type
		nodeInfo.Operator.Contact = contact
		nodeInfo.Operator.OrganizationCountryID = o.Organization.OrganizationCountryID
		nodeInfo.AddressInfo = addressInfo
		nodeInfo.Name = o.Name
		nodeInfo.Type = o.Type
		nodeInfo.AddressID = addressInfo.ID

		now := time.Now()
		nodeInfo.CreatedAt = now
		nodeInfo.UpdatedAt = now
		conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&nodeInfo).Error; err != nil {
				return err
			}
			return nil
		})

		return nil
	}
}

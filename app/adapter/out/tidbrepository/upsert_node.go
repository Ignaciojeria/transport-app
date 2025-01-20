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

type UpsertNode func(context.Context, domain.Origin) error

func init() {
	ioc.Registry(NewUpsertNode, tidb.NewTIDBConnection)
}
func NewUpsertNode(conn tidb.TIDBConnection) UpsertNode {
	return func(ctx context.Context, o domain.Origin) error {
		addressInfo := mapper.MapAddressInfoTable(o.AddressInfo, o.OrganizationCountryID)
		contact := mapper.MapContactToTable(o.NodeInfo.Operator.Contact, o.OrganizationCountryID)
		nodeInfo := mapper.MapNodeInfoTable(o.NodeInfo, o.OrganizationCountryID, addressInfo.ID)
		nodeInfo.OrganizationCountryID = o.OrganizationCountryID
		nodeInfo.OperatorID = o.NodeInfo.Operator.ID
		nodeInfo.Operator.ID = o.NodeInfo.Operator.ID
		nodeInfo.Operator.Type = o.NodeInfo.Operator.Type
		nodeInfo.Operator.Contact = contact
		nodeInfo.Operator.OrganizationCountryID = o.OrganizationCountryID
		nodeInfo.AddressInfo = addressInfo
		nodeInfo.Name = o.NodeInfo.Name
		nodeInfo.Type = o.NodeInfo.Type
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

package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNodeInfo func(context.Context, domain.NodeInfo) (domain.NodeInfo, error)

func init() {
	ioc.Registry(NewUpsertNodeInfo, tidb.NewTIDBConnection)
}
func NewUpsertNodeInfo(conn tidb.TIDBConnection) UpsertNodeInfo {
	return func(ctx context.Context, ni domain.NodeInfo) (domain.NodeInfo, error) {
		nodeInfo := table.NodeInfo{}
		err := conn.DB.WithContext(ctx).Table("node_infos").
			Where("reference_id = ? AND organization_id = ?",
				string(ni.ReferenceID), ni.Organization.ID).First(&nodeInfo).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NodeInfo{}, err
		}
		m := nodeInfo.Map()
		nodeWithChanges, _ := m.UpdateIfChanged(ni)
		dbNodeToUpsert := mapper.MapNodeInfoTable(nodeWithChanges)
		dbNodeToUpsert.CreatedAt = nodeInfo.CreatedAt
		upsertQuery := conn.DB

		if err := upsertQuery.
			Omit("Organization").
			Omit("NodeType").
			Omit("Contact").
			Omit("AddressInfo").
			Save(&dbNodeToUpsert).Error; err != nil {
			return domain.NodeInfo{}, err
		}
		return dbNodeToUpsert.Map(), nil
	}
}

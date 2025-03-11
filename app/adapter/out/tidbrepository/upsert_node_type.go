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

type UpsertNodeType func(context.Context, domain.NodeType) (domain.NodeType, error)

func init() {
	ioc.Registry(NewUpsertNodeType, tidb.NewTIDBConnection)
}
func NewUpsertNodeType(conn tidb.TIDBConnection) UpsertNodeType {
	return func(ctx context.Context, nt domain.NodeType) (domain.NodeType, error) {
		var nodeType table.NodeType
		err := conn.DB.WithContext(ctx).
			Table("node_types").
			Where("name = ? AND organization_id = ?",
				nt.Value,
				nt.Organization.ID).
			First(&nodeType).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NodeType{}, err
		}
		nodeTypeWithChanges := nodeType.Map().UpdateIfChanged(nt)
		DBNodeTypeToUpdate := mapper.MapNodeType(nodeTypeWithChanges)
		DBNodeTypeToUpdate.CreatedAt = nodeType.CreatedAt
		if err := conn.Save(&DBNodeTypeToUpdate).Error; err != nil {
			return domain.NodeType{}, err
		}
		return DBNodeTypeToUpdate.Map(), nil
	}
}

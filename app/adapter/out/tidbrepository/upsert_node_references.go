package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNodeReferences func(context.Context, domain.NodeInfo) error

func init() {
	ioc.Registry(
		NewUpsertNodeReferences,
		database.NewConnectionFactory)
}

func NewUpsertNodeReferences(conn database.ConnectionFactory) UpsertNodeReferences {
	return func(ctx context.Context, node domain.NodeInfo) error {
		nodeDocID := node.DocID(ctx)
		nodeReferences := mapper.MapNodeReferences(ctx, node)

		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("node_doc = ?", nodeDocID).
				Delete(&table.NodeReferences{}).Error; err != nil {
				return err
			}
			if len(nodeReferences) > 0 {
				if err := tx.Save(&nodeReferences).Error; err != nil {
					return err
				}
			}
			if len(nodeReferences) == 0 {
				if err := tx.Create(&table.NodeReferences{
					NodeDoc: nodeDocID.String(),
				}).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
}

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

type UpsertNodeType func(context.Context, domain.NodeType) error

func init() {
	ioc.Registry(NewUpsertNodeType, tidb.NewTIDBConnection)
}

func NewUpsertNodeType(conn tidb.TIDBConnection) UpsertNodeType {
	return func(ctx context.Context, nt domain.NodeType) error {
		var existing table.NodeType

		err := conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newRecord := mapper.MapNodeType(nt)
			return conn.Omit("Organization").Create(&newRecord).Error
		}

		// Ya existe → actualizar solo si hubo cambios
		updated, changed := existing.Map().UpdateIfChange(nt)
		if !changed {
			return nil
		}

		updateData := mapper.MapNodeType(updated)
		updateData.ID = existing.ID
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Organization").Save(&updateData).Error
	}
}

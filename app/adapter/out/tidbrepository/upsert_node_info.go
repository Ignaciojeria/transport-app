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

type UpsertNodeInfo func(context.Context, domain.NodeInfo) error

func init() {
	ioc.Registry(NewUpsertNodeInfo, tidb.NewTIDBConnection)
}

func NewUpsertNodeInfo(conn tidb.TIDBConnection) UpsertNodeInfo {
	return func(ctx context.Context, ni domain.NodeInfo) error {
		var existing table.NodeInfo

		err := conn.DB.WithContext(ctx).
			Table("node_infos").
			Where("document_id = ?", ni.DocID()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRecord := mapper.MapNodeInfoTable(ni)
			return conn.Omit("Organization", "NodeType", "Contact", "AddressInfo").
				Create(&newRecord).Error
		}

		existingMapped := existing.Map()
		updated, changed := existingMapped.UpdateIfChanged(ni)

		// Verificación por hash de relaciones anidadas
		if ni.NodeType.DocID().ShouldUpdate(existing.NodeTypeDoc) {
			updated.NodeType = ni.NodeType
			changed = true
		}
		/*
			if ni.Contact.DocID().ShouldUpdate(existing.ContactDoc) {
				updated.Contact = ni.Contact
				changed = true
			}
			if ni.AddressInfo.DocID().ShouldUpdate(existing.AddressDoc) {
				updated.AddressInfo = ni.AddressInfo
				changed = true
			}
		*/

		if !changed {
			return nil
		}

		updateData := mapper.MapNodeInfoTable(updated)
		updateData.ID = existing.ID
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Organization", "NodeType", "Contact", "AddressInfo").
			Save(&updateData).Error
	}
}

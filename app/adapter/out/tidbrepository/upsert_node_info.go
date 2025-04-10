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

		docID := ni.DocID(ctx)
		if docID.IsZero() {
			return errors.New("cannot persist node info with empty ReferenceID")
		}

		err := conn.DB.WithContext(ctx).
			Table("node_infos").
			Where("document_id = ?", docID).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRecord := mapper.MapNodeInfoTable(ctx, ni)
			// Use the same Omit pattern for consistency
			return conn.
				Create(&newRecord).Error
		}

		existingMapped := existing.Map()
		updated, changed := existingMapped.UpdateIfChanged(ni)

		// Verificación por hash de relaciones anidadas
		if ni.NodeType.DocID(ctx).ShouldUpdate(existing.NodeTypeDoc) {
			updated.NodeType = ni.NodeType
			changed = true
		}

		contactHash := ni.Contact.DocID(ctx)
		if contactHash.ShouldUpdate(existing.ContactDoc) {
			updated.Contact = ni.Contact
			changed = true
		}

		if ni.AddressInfo.DocID(ctx).ShouldUpdate(existing.AddressInfoDoc) {
			updated.AddressInfo = ni.AddressInfo
			changed = true
		}

		if !changed {
			return nil
		}

		updateData := mapper.MapNodeInfoTable(ctx, updated)
		updateData.ID = existing.ID
		updateData.CreatedAt = existing.CreatedAt

		// Use the same Omit pattern for consistency
		return conn.
			Save(&updateData).Error
	}
}

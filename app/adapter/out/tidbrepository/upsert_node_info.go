package tidbrepository

import (
	"context"
	"errors"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNodeInfo func(context.Context, domain.NodeInfo, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertNodeInfo, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertNodeInfo(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertNodeInfo {
	return func(ctx context.Context, ni domain.NodeInfo, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.NodeInfo

			docID := ni.DocID(ctx)
			if docID.IsZero() {
				return errors.New("cannot persist node info with empty ReferenceID")
			}

			err := tx.WithContext(ctx).
				Table("node_infos").
				Where("document_id = ?", docID).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				newRecord := mapper.MapNodeInfoTable(ctx, ni)
				// Use the same Omit pattern for consistency
				if err := tx.Create(&newRecord).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			existingMapped := existing.Map()
			updated, changed := existingMapped.UpdateIfChanged(ni)

			// Verificación por hash de relaciones anidadas
			if ni.NodeType.DocID(ctx).ShouldUpdate(existing.NodeTypeDoc) {
				updated.NodeType = ni.NodeType
				changed = true
			}

			contactHash := ni.AddressInfo.Contact.DocID(ctx)
			if contactHash.ShouldUpdate(existing.ContactDoc) {
				updated.AddressInfo.Contact = ni.AddressInfo.Contact
				changed = true
			}

			if ni.AddressInfo.DocID(ctx).ShouldUpdate(existing.AddressInfoDoc) {
				updated.AddressInfo = ni.AddressInfo
				changed = true
			}

			// Manejo de headers según los tres escenarios
			if !ni.Headers.IsEmpty() {
				// Si las cabeceras nuevas no están vacías, actualizamos
				updated.Headers = ni.Headers
				changed = true
			} else {
				// Si las cabeceras nuevas están vacías, mantenemos las existentes
				updated.Headers = existingMapped.Headers
			}

			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapNodeInfoTable(ctx, updated)
			updateData.ID = existing.ID
			updateData.CreatedAt = existing.CreatedAt

			// Use the same Omit pattern for consistency
			if err := tx.Table("node_infos").
				Where("document_id = ?", docID).
				Updates(updateData).Error; err != nil {
				return err
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}

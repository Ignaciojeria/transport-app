package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertDeliveryUnits func(context.Context, []domain.DeliveryUnit) error

func init() {
	ioc.Registry(NewUpsertDeliveryUnits, database.NewConnectionFactory)
}

func NewUpsertDeliveryUnits(conn database.ConnectionFactory) UpsertDeliveryUnits {
	return func(ctx context.Context, pcks []domain.DeliveryUnit) error {

		// 1. Expandimos paquetes sin LPN en paquetes individuales
		var normalized []domain.DeliveryUnit
		for _, pkg := range pcks {
			normalized = append(normalized, pkg)
		}

		if len(pcks) == 0 {
			normalized = append(normalized, domain.DeliveryUnit{})
		}

		// 2. Construimos una lista de DocumentIDs
		docIDs := make([]string, 0, len(normalized))
		docIDToPackage := make(map[string]domain.DeliveryUnit, len(normalized))
		for _, p := range normalized {
			docID := string(p.DocID(ctx))
			docIDs = append(docIDs, docID)
			docIDToPackage[docID] = p
		}

		// 3. Traemos todos los paquetes existentes con un IN
		var existingDBPackages []table.DeliveryUnit
		err := conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id IN ?", docIDs).
			Find(&existingDBPackages).Error

		if err != nil {
			return err
		}

		// 4. Creamos un map de paquetes existentes por documentID
		existingMap := make(map[string]table.DeliveryUnit)
		for _, pkg := range existingDBPackages {
			existingMap[pkg.DocumentID] = pkg
		}

		// 5. Preparamos los paquetes a upsertear
		var DBpackagesToUpsert []table.DeliveryUnit
		for _, docID := range docIDs {
			domainPkg := docIDToPackage[docID]
			if existingPkg, found := existingMap[docID]; found {
				updatedDomainPkg, _ := existingPkg.Map().UpdateIfChanged(domainPkg)
				updatedTablePkg := mapper.MapPackageToTable(ctx, updatedDomainPkg)
				updatedTablePkg.SizeCategoryDoc = domainPkg.SizeCategory.DocumentID(ctx).String()
				// Preservar campos importantes
				updatedTablePkg.ID = existingPkg.ID
				updatedTablePkg.CreatedAt = existingPkg.CreatedAt
				updatedTablePkg.DocumentID = existingPkg.DocumentID

				DBpackagesToUpsert = append(DBpackagesToUpsert, updatedTablePkg)
			} else {
				newTablePkg := mapper.MapPackageToTable(ctx, domainPkg)
				DBpackagesToUpsert = append(DBpackagesToUpsert, newTablePkg)
			}
		}

		// 6. Guardamos
		return conn.Save(&DBpackagesToUpsert).Error
	}
}

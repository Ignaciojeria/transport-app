package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertPackages func(context.Context, []domain.Package, string) error

func init() {
	ioc.Registry(NewUpsertPackages, tidb.NewTIDBConnection)
}
func NewUpsertPackages(conn tidb.TIDBConnection) UpsertPackages {
	return func(ctx context.Context, pcks []domain.Package, orderReference string) error {
		if len(pcks) == 0 {
			return nil
		}

		// 1. Construimos una lista de DocumentIDs
		docIDs := make([]string, 0, len(pcks))
		docIDToPackage := make(map[string]domain.Package, len(pcks))
		for _, p := range pcks {
			docID := string(p.DocID(ctx, orderReference))
			docIDs = append(docIDs, docID)
			docIDToPackage[docID] = p
		}

		// 2. Traemos todos los paquetes existentes con un IN
		var existingDBPackages []table.Package
		err := conn.DB.WithContext(ctx).
			Table("packages").
			Where("document_id IN ?", docIDs).
			Find(&existingDBPackages).Error
		if err != nil {
			return err
		}

		// 3. Creamos un map de paquetes existentes por documentID
		existingMap := make(map[string]table.Package)
		for _, pkg := range existingDBPackages {
			existingMap[pkg.DocumentID] = pkg
		}

		// 4. Preparamos los paquetes a upsertear
		var DBpackagesToUpsert []table.Package
		for _, docID := range docIDs {
			domainPkg := docIDToPackage[docID]
			if existingPkg, found := existingMap[docID]; found {
				updatedDomainPkg, _ := existingPkg.Map().UpdateIfChanged(domainPkg)
				updatedTablePkg := mapper.MapPackageToTable(ctx, updatedDomainPkg, orderReference)

				// Preservar campos importantes
				updatedTablePkg.ID = existingPkg.ID
				updatedTablePkg.CreatedAt = existingPkg.CreatedAt
				updatedTablePkg.DocumentID = existingPkg.DocumentID

				DBpackagesToUpsert = append(DBpackagesToUpsert, updatedTablePkg)
			} else {
				newTablePkg := mapper.MapPackageToTable(ctx, domainPkg, orderReference)
				DBpackagesToUpsert = append(DBpackagesToUpsert, newTablePkg)
			}
		}

		// 5. Guardamos
		return conn.Save(&DBpackagesToUpsert).Error
	}
}

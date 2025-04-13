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

		// 1. Expandimos paquetes sin LPN en paquetes individuales
		var normalized []domain.Package
		for _, pkg := range pcks {
			normalized = append(normalized, pkg.ExplodeIfNoLpn()...)
		}

		if len(pcks) == 0 {
			normalized = append(normalized, domain.Package{})
		}

		// 2. Construimos una lista de DocumentIDs
		docIDs := make([]string, 0, len(normalized))
		docIDToPackage := make(map[string]domain.Package, len(normalized))
		for _, p := range normalized {
			docID := string(p.DocID(ctx, orderReference))
			docIDs = append(docIDs, docID)
			docIDToPackage[docID] = p
		}

		// 3. Traemos todos los paquetes existentes con un IN
		var existingDBPackages []table.Package
		err := conn.DB.WithContext(ctx).
			Table("packages").
			Where("document_id IN ?", docIDs).
			Find(&existingDBPackages).Error
		if err != nil {
			return err
		}

		// 4. Creamos un map de paquetes existentes por documentID
		existingMap := make(map[string]table.Package)
		for _, pkg := range existingDBPackages {
			existingMap[pkg.DocumentID] = pkg
		}

		// 5. Preparamos los paquetes a upsertear
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

		// 6. Guardamos
		return conn.Save(&DBpackagesToUpsert).Error
	}
}

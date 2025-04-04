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

type UpsertPackages func(context.Context, []domain.Package, domain.Organization) error

func init() {
	ioc.Registry(NewUpsertPackages, tidb.NewTIDBConnection)
}
func NewUpsertPackages(conn tidb.TIDBConnection) UpsertPackages {
	return func(ctx context.Context, pcks []domain.Package, org domain.Organization) error {
		if len(pcks) == 0 {
			return nil
		}

		var DBpackagesToUpsert []table.Package

		// Procesar cada paquete individualmente
		for _, pck := range pcks {
			// Buscar directamente por lpn y organization_id
			var existingPackage table.Package
			err := conn.DB.WithContext(ctx).
				Table("packages").
				Preload("Organization").
				Where("document_id = ?", pck.DocID()).
				First(&existingPackage).Error

			// Si no hay error, significa que encontramos el registro
			if err == nil {
				// Actualizar existente
				updatedDomainPkg := existingPackage.Map().UpdateIfChanged(pck)
				updatedTablePkg := mapper.MapPackageToTable(updatedDomainPkg)

				// Preservar campos importantes
				updatedTablePkg.ID = existingPackage.ID
				updatedTablePkg.CreatedAt = existingPackage.CreatedAt
				updatedTablePkg.DocumentID = existingPackage.DocumentID

				DBpackagesToUpsert = append(DBpackagesToUpsert, updatedTablePkg)
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe - insertar nuevo
				newTablePkg := mapper.MapPackageToTable(pck)
				DBpackagesToUpsert = append(DBpackagesToUpsert, newTablePkg)
			} else {
				// Error de BD distinto a "no encontrado"
				return err
			}
		}

		// Guardar todos los paquetes
		if err := conn.Save(&DBpackagesToUpsert).Error; err != nil {
			return err
		}

		return nil
	}
}

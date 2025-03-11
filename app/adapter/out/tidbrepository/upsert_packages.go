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

type UpsertPackages func(context.Context, []domain.Package, domain.Organization) ([]domain.Package, error)

func init() {
	ioc.Registry(NewUpsertPackages, tidb.NewTIDBConnection)
}
func NewUpsertPackages(conn tidb.TIDBConnection) UpsertPackages {
	return func(ctx context.Context, pcks []domain.Package, org domain.Organization) ([]domain.Package, error) {
		if len(pcks) == 0 {
			return []domain.Package{}, nil
		}

		var packagesTable []table.Package
		var lpns []string
		for _, v := range pcks {
			lpns = append(lpns, v.Lpn)
		}

		err := conn.DB.WithContext(ctx).Table("packages").
			Where("lpn IN ? AND organization_id = ?", lpns, org.ID).Find(&packagesTable).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return []domain.Package{}, err
		}

		// Crear mapa de paquetes existentes
		existingPackages := make(map[string]table.Package)
		for _, v := range packagesTable {
			existingPackages[v.Lpn] = v
		}

		var DBpackagesToUpsert []table.Package
		// Procesar todos los paquetes
		for _, pck := range pcks {
			if existing, ok := existingPackages[pck.Lpn]; ok {
				// Actualizar existente
				DBpackagesToUpsert = append(DBpackagesToUpsert,
					mapper.MapPackageToTable(existing.Map().UpdateIfChanged(pck), org.ID))
			} else {
				// Insertar nuevo
				DBpackagesToUpsert = append(DBpackagesToUpsert,
					mapper.MapPackageToTable(pck, org.ID))
			}
		}

		if err := conn.Save(&DBpackagesToUpsert).Error; err != nil {
			return []domain.Package{}, err
		}

		var upsertedPackagesDomain []domain.Package
		for _, v := range DBpackagesToUpsert {
			upsertedPackagesDomain = append(upsertedPackagesDomain, v.Map())
		}

		return upsertedPackagesDomain, nil
	}
}

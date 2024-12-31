package table

import (
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		MigrateTables,
		tidb.NewTIDBConnection,
		configuration.NewTiDBConfiguration)
}

func MigrateTables(
	conn tidb.TIDBConnection,
	conf configuration.TiDBConfiguration) error {
	if conf.TIDB_RUN_MIGRATIONS != "true" {
		return nil
	}
	// Lista de tablas que tienen un campo ID como clave primaria
	tables := []interface{}{
		&TransportOrder{},
		&Origin{},
		&Destination{},
		&AddressInfo{},
		&TransportOrderReferences{},
		&TransportRequirementsReferences{},
		&NodeReferences{},
		&NodeInfo{},
		&Operator{},
		&Items{},
		&Packages{},
		&ItemReferences{},
		&OrderStatus{},
		&Visit{},
		&Commerce{},
		&Consumer{},
		&Organization{},
		&OrderType{},
	}

	// Opcional: Eliminar tablas si existen
	for _, table := range tables {
		if err := conn.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	// Crear las tablas nuevamente
	if err := conn.AutoMigrate(tables...); err != nil {
		return err
	}

	return nil
}

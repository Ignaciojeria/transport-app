package table

import (
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewRunMigrations,
		tidb.NewTIDBConnection,
		configuration.NewTiDBConfiguration)
}

type RunMigrations func() error

func NewRunMigrations(
	conn tidb.TIDBConnection,
	conf configuration.TiDBConfiguration) RunMigrations {
	return func() error {
		if conf.TIDB_RUN_MIGRATIONS != "true" {
			return nil
		}
		// Lista de tablas que tienen un campo ID como clave primaria
		tables := []interface{}{
			&OrdersOutbox{},
			&NodesOutbox{},
			&VehiclesOutbox{},
			&Order{},
			&AddressInfo{},
			&OrderReferences{},
			&NodeInfo{},
			&Operator{},
			&Package{},
			&OrderPackage{},
			&OrderStatus{},
			&Commerce{},
			&Consumer{},
			&Organization{},
			&OrderType{},
			&ApiKey{},
			&OrganizationCountry{},
			&Account{},
			&Contact{},
			&Carrier{},
			&Vehicle{},
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
}

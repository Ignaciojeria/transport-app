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
			&EventOutbox{},
			&PlanType{},
			&PlanningStatus{},
			&PlanType{},
			&Route{},
			&OrderHeaders{},
			&VehicleHeaders{},
			&VehicleCategory{},
			&NodeHeaders{},
			&NodeType{},
			&Order{},
			&AddressInfo{},
			&OrderReferences{},
			&NodeInfo{},
			&Package{},
			&OrderPackage{},
			&OrderStatus{},
			&Organization{},
			&OrderType{},
			//&ApiKey{},
			//&OrganizationCountry{},
			&Account{},
			&Contact{},
			&Carrier{},
			&Vehicle{},
			&CheckoutRejection{},
			&CheckoutHistory{},
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

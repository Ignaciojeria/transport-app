package table

import (
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewRunMigrations,
		database.NewConnectionFactory,
		configuration.NewTiDBConfiguration)
}

type RunMigrations func() error

func NewRunMigrations(
	conn database.ConnectionFactory,
	conf configuration.DBConfiguration) RunMigrations {
	return func() error {
		if conf.DB_RUN_MIGRATIONS != "true" {
			return nil
		}
		// Lista de tablas que tienen un campo ID como clave primaria
		tables := []interface{}{
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
			&DeliveryUnit{},
			&OrderDeliveryUnit{},
			&Status{},
			&Organization{},
			&OrderType{},
			&Account{},
			&AccountOrganization{},
			&Contact{},
			&Carrier{},
			&Vehicle{},
			&CheckoutRejection{},
			&CheckoutHistory{},
			&State{},
			&District{},
			&Province{},
			&NonDeliveryReason{},
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

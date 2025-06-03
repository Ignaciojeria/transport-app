package table

import (
	"log"
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
			log.Println("Migrations disabled, skipping...")
			return nil
		}
		log.Println("Starting migrations...")
		// Lista de tablas que tienen un campo ID como clave primaria
		tables := []interface{}{
			&Route{},
			&OrderHeaders{},
			&VehicleHeaders{},
			&VehicleCategory{},
			&NodeInfoHeaders{},
			&NodeType{},
			&Order{},
			&AddressInfo{},
			&OrderReferences{},
			&NodeInfo{},
			&DeliveryUnit{},
			&OrderDeliveryUnit{},
			&Status{},
			&Tenant{},
			&OrderType{},
			&Account{},
			&AccountTenant{},
			&Contact{},
			&Carrier{},
			&Vehicle{},
			&Driver{},
			&State{},
			&District{},
			&Province{},
			&NonDeliveryReason{},
			&Plan{},
			&PlanHeaders{},
			&Route{},
			&DeliveryUnitsStatusHistory{},
			&SizeCategory{},
			&DeliveryUnitsLabels{},
		}

		// Crear las tablas nuevamente
		log.Println("Creating tables...")
		if err := conn.AutoMigrate(tables...); err != nil {
			log.Printf("Error creating tables: %v", err)
			return err
		}
		log.Println("Migrations completed successfully")
		return nil
	}
}

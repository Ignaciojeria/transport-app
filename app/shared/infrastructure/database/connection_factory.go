package database

import (
	"fmt"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

// connectionFactory es un tipo que representa una función que retorna una conexión GORM lista para usar.
type connectionFactory func() (*gorm.DB, error)

// ConnectionFactory representa la conexión activa y la estrategia seleccionada (ej: "postgresql", "tidb").
type ConnectionFactory struct {
	*gorm.DB
	Strategy string
}

func init() {
	ioc.Registry(
		NewConnectionFactory,
		configuration.NewTiDBConfiguration,
		NewPostgreSQLConnectionFactory,
		NewTiDBConnectionFactory,
	)
}

// NewConnectionFactory selecciona y crea una conexión basada en la estrategia definida en la configuración.
func NewConnectionFactory(
	env configuration.DBConfiguration,
	postgresqlStrategy connectionFactory,
	tidbStrategy connectionFactory,
) (ConnectionFactory, error) {
	var (
		db  *gorm.DB
		err error
	)

	switch env.DB_STRATEGY {
	case "tidb":
		db, err = tidbStrategy()
	case "postgresql":
		db, err = postgresqlStrategy()
	case "disabled":
		// Retorna una conexión nula para deshabilitar la persistencia
		return ConnectionFactory{
			DB:       nil,
			Strategy: "disabled",
		}, nil
	default:
		return ConnectionFactory{}, fmt.Errorf("unknown DB_STRATEGY: %s", env.DB_STRATEGY)
	}

	if err != nil {
		return ConnectionFactory{}, err
	}

	return ConnectionFactory{
		DB:       db,
		Strategy: env.DB_STRATEGY,
	}, nil
}

package tidb

import (
	"errors"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewTIDBConnection,
		configuration.NewTiDBConfiguration,
		NewPostgreSQLConnectionStrategy,
		NewTIDBConnectionStrategy)
}

type TIDBConnection struct {
	*gorm.DB
}

func NewTIDBConnection(
	env configuration.DBConfiguration,
	postgresqlStrategy connectionStrategy,
	tidbStrategy connectionStrategy,
) (TIDBConnection, error) {
	var strategy *gorm.DB
	var err error
	if env.DB_STRATEGY == "tidb" {
		strategy, err = tidbStrategy()
	}
	if env.DB_STRATEGY == "postgresql" {
		strategy, err = postgresqlStrategy()
	}
	if err != nil {
		return TIDBConnection{}, err
	}
	if strategy == nil {
		return TIDBConnection{}, errors.New("unknown strategy: " + env.DB_STRATEGY)
	}
	return TIDBConnection{strategy}, err
}

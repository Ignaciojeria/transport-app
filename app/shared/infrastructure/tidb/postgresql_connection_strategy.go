package tidb

import (
	"fmt"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

func init() {
	ioc.Registry(newPostgreSQLConnectionStrategy, configuration.NewTiDBConfiguration)
}
func newPostgreSQLConnectionStrategy(env configuration.DBConfiguration) connectionStrategy {
	return func() (*gorm.DB, error) {
		username := env.DB_USERNAME
		pwd := env.DB_PASSWORD
		host := env.DB_HOSTNAME
		port := env.DB_PORT
		dbname := env.DB_NAME
		sslMode := env.DB_SSL_MODE
		config := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		}
		db, err := gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, pwd, host, port, dbname, sslMode)), config)
		if err != nil {
			return nil, err
		}
		if err := db.Use(tracing.NewPlugin()); err != nil {
			return nil, err
		}
		return db, nil
	}
}

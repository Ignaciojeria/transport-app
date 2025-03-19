package tidb

import (
	"crypto/tls"
	"fmt"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	tidbmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(newTIDBConnectionStrategy, configuration.NewTiDBConfiguration)
}
func newTIDBConnectionStrategy(env configuration.DBConfiguration) connectionStrategy {
	return func() (*gorm.DB, error) {
		err := tidbmysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: env.DB_HOSTNAME,
		})
		if err != nil {
			return &gorm.DB{}, fmt.Errorf("failed to register TLS config: %w", err)
		}

		// Create the DSN (Data Source Name)
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?tls=tidb&parseTime=true&charset=utf8mb4&loc=Local",
			env.DB_USERNAME,
			env.DB_PASSWORD,
			env.DB_HOSTNAME,
			env.DB_PORT,
			env.DB_NAME,
		)

		// Open a GORM connection
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return &gorm.DB{}, fmt.Errorf("failed to open GORM connection: %w", err)
		}

		// Optional: Ping the database to ensure the connection works
		sqlDB, err := db.DB()
		if err != nil {
			return &gorm.DB{}, fmt.Errorf("failed to retrieve raw DB connection: %w", err)
		}

		if err := sqlDB.Ping(); err != nil {
			return &gorm.DB{}, fmt.Errorf("failed to ping database: %w", err)
		}
		return db, nil
	}
}

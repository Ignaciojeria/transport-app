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
	ioc.Registry(NewTIDBConnection, configuration.NewTiDBConfiguration)
}

type TIDBConnection struct {
	*gorm.DB
}

func NewTIDBConnection(env configuration.TiDBConfiguration) (TIDBConnection, error) {
	// Register the custom TLS configuration
	err := tidbmysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: env.TIDB_HOSTNAME,
	})
	if err != nil {
		return TIDBConnection{}, fmt.Errorf("failed to register TLS config: %w", err)
	}

	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=tidb&parseTime=true&charset=utf8mb4&loc=Local",
		env.TIDB_USERNAME,
		env.TIDB_PASSWORD,
		env.TIDB_HOSTNAME,
		env.TIDB_PORT,
		env.TIDB_DATABASE,
	)

	// Open a GORM connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return TIDBConnection{}, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	// Optional: Ping the database to ensure the connection works
	sqlDB, err := db.DB()
	if err != nil {
		return TIDBConnection{}, fmt.Errorf("failed to retrieve raw DB connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return TIDBConnection{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return TIDBConnection{db}, nil
}

package tidb

import (
	"transport-app/app/shared/configuration"
	"crypto/tls"
	"database/sql"
	"fmt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-sql-driver/mysql"
)

func init() {
	ioc.Registry(NewTiDBConnection, configuration.NewTiDBConfiguration)
}

func NewTiDBConnection(env configuration.TiDBConfiguration) (*sql.DB, error) {

	err := mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: env.TIDB_HOSTNAME,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register TLS config: %w", err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=tidb",
		env.TIDB_USERNAME,
		env.TIDB_PASSWORD,
		env.TIDB_HOSTNAME,
		env.TIDB_PORT,
		env.TIDB_DATABASE,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

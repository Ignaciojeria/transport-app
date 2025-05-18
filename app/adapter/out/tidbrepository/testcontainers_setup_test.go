package tidbrepository

import (
	"context"
	"testing"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestContainersSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tidb Repository Suite")
}

var container *tcpostgres.PostgresContainer
var connection database.ConnectionFactory

var noTablesContainerConnection database.ConnectionFactory
var noTablesMigrationContainer *tcpostgres.PostgresContainer

var _ = BeforeSuite(func() {
	ctx := context.Background()
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	Expect(err).ToNot(HaveOccurred())
	container = postgresContainer

	// Obtener host y puerto del contenedor
	host, err := container.Host(ctx)
	Expect(err).ToNot(HaveOccurred())

	port, err := container.MappedPort(ctx, "5432")
	Expect(err).ToNot(HaveOccurred())

	connection, err = database.NewConnectionFactory(
		configuration.DBConfiguration{DB_STRATEGY: "postgresql"},
		database.NewPostgreSQLConnectionFactory(configuration.DBConfiguration{
			DB_HOSTNAME:       host,
			DB_PORT:           port.Port(), // devuelve string
			DB_SSL_MODE:       "disable",   // SSL deshabilitado para test local
			DB_NAME:           dbName,
			DB_USERNAME:       dbUser,
			DB_PASSWORD:       dbPassword,
			DB_RUN_MIGRATIONS: "true",
		}),
		nil,
	)
	Expect(err).ToNot(HaveOccurred())
	err = table.NewRunMigrations(connection, configuration.DBConfiguration{
		DB_RUN_MIGRATIONS: "true",
	})()
	Expect(err).ToNot(HaveOccurred())

	// No tables container setup (remains unchanged)
	noTablesContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	Expect(err).ToNot(HaveOccurred())
	noTablesMigrationContainer = noTablesContainer

	noTablesHost, err := noTablesMigrationContainer.Host(ctx)
	Expect(err).ToNot(HaveOccurred())

	noTablesPort, err := noTablesMigrationContainer.MappedPort(ctx, "5432")
	Expect(err).ToNot(HaveOccurred())

	noTablesContainerConnection, err = database.NewConnectionFactory(
		configuration.DBConfiguration{DB_STRATEGY: "postgresql"},
		database.NewPostgreSQLConnectionFactory(configuration.DBConfiguration{
			DB_HOSTNAME:       noTablesHost,
			DB_PORT:           noTablesPort.Port(),
			DB_SSL_MODE:       "disable",
			DB_NAME:           dbName,
			DB_USERNAME:       dbUser,
			DB_PASSWORD:       dbPassword,
			DB_RUN_MIGRATIONS: "false",
		}),
		nil,
	)
	Expect(err).ToNot(HaveOccurred())
})

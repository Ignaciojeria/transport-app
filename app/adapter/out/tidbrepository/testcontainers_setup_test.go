package tidbrepository

import (
	"context"
	"testing"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/tidb"
	"transport-app/app/shared/sharedcontext"

	"go.opentelemetry.io/otel/baggage"

	"github.com/biter777/countries"
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
var connection tidb.TIDBConnection
var organization1 domain.Organization
var organization2 domain.Organization

var noTablesContainerConnection tidb.TIDBConnection
var noTablesMigrationContainer *tcpostgres.PostgresContainer

var _ = Describe("TidbRepository", func() {
	It("dummy test", func() {
		Expect(true).To(BeTrue())
	})
})

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

	connection, err = tidb.NewTIDBConnection(
		configuration.DBConfiguration{DB_STRATEGY: "postgresql"},
		tidb.NewPostgreSQLConnectionStrategy(configuration.DBConfiguration{
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

	// Create test account first
	err = NewUpsertAccount(connection)(ctx, domain.Operator{
		Contact: domain.Contact{
			PrimaryEmail: "ignaciovl.j@gmail.com",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	// Create baggage member for tenant country
	countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, countries.Chile.String())
	bag, _ := baggage.New(countryMember)
	ctxWithCountry := baggage.ContextWithBaggage(ctx, bag)

	// Create first organization using the new function
	saveOrganization := NewSaveOrganization(connection)
	organization1, err = saveOrganization(
		ctxWithCountry,
		domain.Operator{
			Contact: domain.Contact{
				PrimaryEmail: "ignaciovl.j@gmail.com",
			},
		},
		"Organization 1",
	)
	Expect(err).ToNot(HaveOccurred())

	// Create second organization using the new function
	organization2, err = saveOrganization(
		ctxWithCountry,
		domain.Operator{
			Contact: domain.Contact{
				PrimaryEmail: "ignaciovl.j@gmail.com",
			},
		},
		"Organization 2",
	)
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

	noTablesContainerConnection, err = tidb.NewTIDBConnection(
		configuration.DBConfiguration{DB_STRATEGY: "postgresql"},
		tidb.NewPostgreSQLConnectionStrategy(configuration.DBConfiguration{
			DB_HOSTNAME:       noTablesHost,
			DB_PORT:           noTablesPort.Port(),
			DB_SSL_MODE:       "disable",
			DB_NAME:           dbName,
			DB_USERNAME:       dbUser,
			DB_PASSWORD:       dbPassword,
			DB_RUN_MIGRATIONS: "false", // importante: no ejecutar migraciones
		}),
		nil,
	)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if container != nil {
		_ = container.Terminate(context.Background())
	}
	if noTablesMigrationContainer != nil {
		_ = noTablesMigrationContainer.Terminate(context.Background())
	}
})

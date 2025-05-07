package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/database"

	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "embed"

	_ "transport-app/app/adapter/in/fuegoapi"
	_ "transport-app/app/adapter/out/tidbrepository"
	_ "transport-app/app/adapter/out/tidbrepository/table"
	_ "transport-app/app/shared/configuration"

	_ "transport-app/app/shared/infrastructure/httpserver"

	_ "transport-app/app/shared/infrastructure/database"

	_ "transport-app/app/usecase"

	_ "transport-app/app/adapter/in/gcpsubscription"
	_ "transport-app/app/adapter/out/gcppublisher"
	_ "transport-app/app/adapter/out/restyclient/locationiq"
	_ "transport-app/app/onload"
	_ "transport-app/app/shared/infrastructure/gcppubsub"
	_ "transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	_ "transport-app/app/shared/infrastructure/httpresty"
	_ "transport-app/app/shared/infrastructure/observability"
	_ "transport-app/app/shared/infrastructure/observability/strategy"
	_ "transport-app/mocks"

	_ "transport-app/app/adapter/out/geocoding"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func TestContainersSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tidb Repository Suite")
}

var container *tcpostgres.PostgresContainer
var connection database.ConnectionFactory
var organization1 domain.Organization
var organization2 domain.Organization

var noTablesContainerConnection database.ConnectionFactory
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
	os.Setenv("version", "testing")
	os.Setenv("DB_STRATEGY", "postgresql")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("DB_HOSTNAME", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_USERNAME", dbUser)
	os.Setenv("DB_PASSWORD", dbPassword)
	os.Setenv("DB_RUN_MIGRATIONS", "true")
	Expect(err).ToNot(HaveOccurred())

	Expect(err).ToNot(HaveOccurred())
	err = table.NewRunMigrations(connection, configuration.DBConfiguration{
		DB_RUN_MIGRATIONS: "true",
	})()
	Expect(err).ToNot(HaveOccurred())

	// Create test account first
	err = tidbrepository.NewUpsertAccount(connection)(ctx, domain.Operator{
		Contact: domain.Contact{
			PrimaryEmail: "ignaciovl.j@gmail.com",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	// Setup context with country information

	// Create first organization using the new function signature
	saveOrganization := tidbrepository.NewSaveOrganization(connection)

	// Create organization entity with required fields
	orgToSave1 := domain.Organization{
		Country: countries.CL,
		Name:    "Organization 1",
		Operator: domain.Operator{
			Contact: domain.Contact{
				PrimaryEmail: "ignaciovl.j@gmail.com",
			},
		},
	}

	// Save the first organization
	organization1, err = saveOrganization(ctx, orgToSave1)
	Expect(err).ToNot(HaveOccurred())

	// Create second organization
	orgToSave2 := domain.Organization{
		Country: countries.CL,
		Name:    "Organization 2",
		Operator: domain.Operator{
			Contact: domain.Contact{
				PrimaryEmail: "ignaciovl.j@gmail.com",
			},
		},
	}

	// Save the second organization
	organization2, err = saveOrganization(ctx, orgToSave2)
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
			DB_RUN_MIGRATIONS: "false", // importante: no ejecutar migraciones
		}),
		nil,
	)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		if err := ioc.LoadDependencies(); err != nil {
			log.Fatalf("failed to load app: %v", err)
		}
	}()

	waitForAppReady := func() {
		maxRetries := 5
		for i := 1; i <= maxRetries; i++ {
			resp, err := http.Get("http://localhost:8080/health") // o tu endpoint real
			if err == nil && resp.StatusCode == http.StatusOK {
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
		Fail("application did not start in time")
	}

	waitForAppReady()

})

var _ = AfterSuite(func() {
	os.Exit(0)
})

type EmbeddedRequest struct {
	Headers map[string]string `json:"headers"`
	Body    json.RawMessage   `json:"body"`
}

func loadRequest(data []byte) (EmbeddedRequest, error) {
	var req EmbeddedRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		return EmbeddedRequest{}, fmt.Errorf("error unmarshaling request: %w", err)
	}
	return req, nil
}

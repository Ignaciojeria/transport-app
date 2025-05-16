package usecase

import (
	"context"
	"testing"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.opentelemetry.io/otel/baggage"
)

func TestDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Suite")
}

var _ = BeforeSuite(func() {
	postgreSQl()
})
var _ = AfterSuite(func() {
	if pgContainer != nil {
		ctx := context.Background()
		err := pgContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to terminate PostgreSQL container")
	}
})

var pgContainer *tcpostgres.PostgresContainer
var connection database.ConnectionFactory

func buildCtx(tenantID, country string) context.Context {
	ctx := context.Background()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID)
	cntry, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, country)
	bag, _ := baggage.New(tID, cntry)
	return baggage.ContextWithBaggage(ctx, bag)
}

func postgreSQl() {
	ctx := context.Background()
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"
	var err error
	pgContainer, err = tcpostgres.Run(ctx,
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
	// Obtener host y puerto del contenedor
	host, err := pgContainer.Host(ctx)
	Expect(err).ToNot(HaveOccurred())

	port, err := pgContainer.MappedPort(ctx, "5432")
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
}

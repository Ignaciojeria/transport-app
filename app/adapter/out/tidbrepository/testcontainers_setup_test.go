package tidbrepository

import (
	"context"
	"testing"
	"time"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/tidb"

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
	tidb.NewTIDBConnection(configuration.DBConfiguration{}, nil, nil)
	container = postgresContainer
})

var _ = AfterSuite(func() {
	if container != nil {
		_ = container.Terminate(context.Background())
	}
})

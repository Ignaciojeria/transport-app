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
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/infrastructure/gcppubsub"

	"cloud.google.com/go/pubsub"
	pubsubtc "github.com/testcontainers/testcontainers-go/modules/gcloud/pubsub"

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

var pgContainer *tcpostgres.PostgresContainer
var connection database.ConnectionFactory
var organization1 domain.Tenant
var organization2 domain.Tenant
var pubsubContainer *pubsubtc.Container

const TRANSPORT_APP_TOPIC = "transport-app-events"
const ORDER_SUBMITTED_SUBSCRIPTION = "transport-app-events-order-submitted"
const TENANT_SUBMITTED_SUBSCRIPTION = "transport-app-events-tenant-submitted"
const REGISTRATION_SUBMITTED_SUBSCRIPTION = "transport-app-events-registration-submitted"

var noTablesContainerConnection database.ConnectionFactory
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

	pubsubContainer, err = pubsubtc.Run(
		ctx,
		"gcr.io/google.com/cloudsdktool/cloud-sdk:367.0.0-emulators",
		testcontainers.WithEnv(map[string]string{
			"PUBSUB_PROJECT_ID": "test-project",
		}),
	)
	Expect(err).ToNot(HaveOccurred())

	pubsubHost, err := pubsubContainer.Host(ctx)
	Expect(err).ToNot(HaveOccurred())

	pubsubPort, err := pubsubContainer.MappedPort(ctx, "8085")
	Expect(err).ToNot(HaveOccurred())
	os.Setenv("GOOGLE_PROJECT_ID", "test-project")
	os.Setenv("PUBSUB_EMULATOR_HOST", fmt.Sprintf("%s:%s", pubsubHost, pubsubPort.Port()))
	os.Setenv("ORDER_SUBMITTED_SUBSCRIPTION", ORDER_SUBMITTED_SUBSCRIPTION)
	os.Setenv("TENANT_SUBMITTED_SUBSCRIPTION", TENANT_SUBMITTED_SUBSCRIPTION)
	os.Setenv("REGISTRATION_SUBMITTED_SUBSCRIPTION", REGISTRATION_SUBMITTED_SUBSCRIPTION)

	os.Setenv("TRANSPORT_APP_TOPIC", TRANSPORT_APP_TOPIC)

	port, err := pgContainer.MappedPort(ctx, "5432")
	Expect(err).ToNot(HaveOccurred())

	os.Setenv("version", "testing")
	os.Setenv("DB_STRATEGY", "postgresql")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("DB_HOSTNAME", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_USERNAME", dbUser)
	os.Setenv("DB_PASSWORD", dbPassword)
	os.Setenv("DB_RUN_MIGRATIONS", "true")

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
			time.Sleep(2000 * time.Millisecond)
		}
		Fail("application did not start in time")
	}

	waitForAppReady()

})

var _ = AfterSuite(func() {
	ctx := context.Background()
	if pgContainer != nil {
		_ = pgContainer.Terminate(ctx)
	}
	if pubsubContainer != nil {
		_ = pubsubContainer.Terminate(ctx)
	}
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

func init() {
	ioc.Registry(SetupPubsubTestResources, gcppubsub.NewClient)
}
func SetupPubsubTestResources(client *pubsub.Client) error {
	ctx := context.Background()

	// Crear el tópico si no existe
	topic := client.Topic(TRANSPORT_APP_TOPIC)
	if _, err := client.CreateTopic(ctx, TRANSPORT_APP_TOPIC); err != nil {
		return err
	}

	_, err := client.CreateSubscription(ctx, ORDER_SUBMITTED_SUBSCRIPTION, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		return err
	}

	_, err = client.CreateSubscription(ctx, TENANT_SUBMITTED_SUBSCRIPTION, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		return err
	}

	_, err = client.CreateSubscription(ctx, REGISTRATION_SUBMITTED_SUBSCRIPTION, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		return err
	}

	return nil
}

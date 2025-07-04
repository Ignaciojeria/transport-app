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
	os.Setenv("JWT_PRIVATE_KEY", `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCd/4LVEzxXI+6o
3PXlU0iF4cjwdp+C/5srH3uUh848cstha0vv5VaUC1+z6AJHiYWn6ZSOM/IyWTr7
QyZkfzs4hwHF0DkhXw222dENmrbo27W5edDZhDMfcFvfI8eU9UTfMqksDCRWEA9p
SYKQpsvkSmjHjCR7s0Hu853/vF3Q1JDM4C6Da8AywO1+ytMew6AbnUlmtZA+sfH0
PGVJ3oJlIPkddO5EVQfZrdZ5a+WJbANg/ZjVtWlx0QxYKPZC0WUfC3hK+Y/M8kHM
mdBqX4nqqkqlzu8LZUpewHsr+2OW14IBF2lFls5SzZzPFXC1Xlk0C8sMJ4yypEYV
VxVyMJNDAgMBAAECggEADbWKbh+VqD/7aNaHVYck70z4iPjZ/V1dYt//8pLYD8Gs
KP9M3vtoyD62Zp9Wd0uU981aMP6L4NeCOSQJ9EVf55c7TaU0F9OuFnQO4h3PCiRq
P7Y0q3L/lpZCunlZy3B+sdks+Z/yFS+ejrWsqQ13+o6ESfA8iCk1Kt0jk/mNsF1m
YYN2jRu46x2NEsVZU14nwbSxD6u5WT5TATAJHV7urFX3mCPs7KHV7kP5OpQ5Lkt7
W0XJmQ0WC1k/yCu2B+NYRoCih5mFmjT+fFm70mK68FYBpIzb84j4XOa1RMHFR5nf
eNmC5/QP23fq4XsIN2f09K18fsfMDyp24GFaD1E+nQKBgQDULCqOm8PM6mywjUFi
Ikf/fAaobxKQewlAf1VRBxvXXwS9GVcfEtGQO3Q2mx7o43u0voGeMv25epxekqY1
cE/j2olj65FoYuw/E6R3ddixOltr04uxF3AS+UV8UHz9fQz0ha8FwZZKYIX+h3Ag
4uQtuNYoL6XbO/CfYga8f5BLRQKBgQC+oo/w3kHeZqxU4Qp703i31j5hNElu31Z6
PYzTPDVC/uDhUhaaeUMQnDnfUfTZx72WKIWHZP0DHlYk/DWDoxtIPEvUkgF6sFay
xqJesJGAVgwyx29+2jus8lGjVeX8VajUiM04NYNTjLos+MyPNinuDAPzNUVOvSfA
xfyGYwCI5wKBgQCdPH3tYZIhcjlKPeSOjUk+FPP6LxZa7FNW8QaRHeuMGGaynOzr
ok6bzPO65ApsHOm4cNYuHyvZIPxxOczjHXCXM4VN/22rJmRd+niP702/SbgmmIeV
ngD4jrLoBd4bHWlUbR3f7i8qv42Nq2F1fbAMEkbjUSxg5HLWKxdC6mZM4QKBgDyA
265W3BD6BTfrNKiYNXgjRykSrzvBJnEll8xzD3Rz8GuS4hmk4uQisTtvh4aXHlTK
B5cxNhwHRM/4PAPLgJ0sheSxcka+MMYMxPvIjmVs6fIz2e8o1EfPoJl2acfsZ+kM
ghWU5lleELi1Pjc1uZkTty05ewXCgxtruvnr8f+BAoGAEJf0xBEqx7eMItB3uLU/
GuGrXMM+6PsTE8CLp8yXFXjVTde3XVvYhGWeNQSW8S0Pn9NnSdyUqtg4cDBzRoFG
1HsRLEp6zC6+0Ak3ngCg1tPfWAvxdtYtB/yYu9Wud8Ee9XAuZJR9AaF9VXOTQdKp
pLbS9a8Ae1bHgtd2gDW2K9A=
-----END PRIVATE KEY-----
`)
	os.Setenv("JWT_PUBLIC_KEY", `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnf+C1RM8VyPuqNz15VNI
heHI8Hafgv+bKx97lIfOPHLLYWtL7+VWlAtfs+gCR4mFp+mUjjPyMlk6+0MmZH87
OIcBxdA5IV8NttnRDZq26Nu1uXnQ2YQzH3Bb3yPHlPVE3zKpLAwkVhAPaUmCkKbL
5Epox4wke7NB7vOd/7xd0NSQzOAug2vAMsDtfsrTHsOgG51JZrWQPrHx9DxlSd6C
ZSD5HXTuRFUH2a3WeWvliWwDYP2Y1bVpcdEMWCj2QtFlHwt4SvmPzPJBzJnQal+J
6qpKpc7vC2VKXsB7K/tjlteCARdpRZbOUs2czxVwtV5ZNAvLDCeMsqRGFVcVcjCT
QwIDAQAB
-----END PUBLIC KEY-----`)
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

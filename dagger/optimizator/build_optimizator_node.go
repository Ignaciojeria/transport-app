package main

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	godotenv.Load()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// ✅ Read secret from environment variable
	natsCreds := client.SetSecret("nats-creds", os.Getenv("NATS_CONNECTION_CREDS_FILECONTENT"))

	println("🚀 Building Transport App...")

	// ✅ Build Transport App container
	println("=== BUILDING TRANSPORT APP ===")
	appContainer := client.Container().
		From("ghcr.io/ignaciojeria/transport-app/transport-app-d0a6ffdd2b5a22c2c0423e7b340b3900@sha256:e1ee7fd378916720caaa961352cf287bede6bfc117638cbf64e47dcd0876abed").
		WithEnvVariable("PORT", "8080").
		WithEnvVariable("TRANSPORT_APP_TOPIC", "transport-app-events").
		WithEnvVariable("VERSION", "1.4.0").
		WithEnvVariable("ENVIRONMENT", "prod").
		WithEnvVariable("DB_STRATEGY", "disabled").
		WithEnvVariable("OBSERVABILITY_STRATEGY", "none").
		WithEnvVariable("OPTIMIZATION_REQUESTED_SUBSCRIPTION", "transport-app-events-optimization-requested").
		WithEnvVariable("MASTER_NODE_URL", "https://einar-main-f0820bc.d2.zuplo.dev").
		WithEnvVariable("VROOM_PLANNER_URL", "http://planner:3000").
		WithEnvVariable("VROOM_OPTIMIZER_URL", "http://optimizer:3000").
		WithEnvVariable("NATS_CONNECTION_URL", "connect.ngs.global").
		WithSecretVariable("NATS_CONNECTION_CREDS_FILECONTENT", natsCreds).
		WithExposedPort(8080)

	println("✅ Transport App container built successfully")

	// ✅ Push Transport App to registry
	println("=== PUSHING TRANSPORT APP TO REGISTRY ===")

	// Get image tag from environment (default to latest if not set)
	imageTag := os.Getenv("IMAGE_TAG")
	if imageTag == "" {
		imageTag = uuid.New().String()
	}

	// Push the Transport App container to registry
	_, err = appContainer.Publish(ctx, "ghcr.io/ignaciojeria/transport-app/transport-app-optimizator-worker:"+imageTag)
	if err != nil {
		panic(err)
	}

	println("✅ Successfully pushed Transport App to registry:")
	println("   Image: ghcr.io/ignaciojeria/transport-app/transport-app-optimizator:" + imageTag)

}

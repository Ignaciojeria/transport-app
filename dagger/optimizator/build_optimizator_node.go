package main

import (
	"context"
	"os"

	"dagger.io/dagger"
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

	println("🚀 Building Transport App from source...")

	// ✅ Build Transport App from source code
	println("=== BUILDING TRANSPORT APP FROM SOURCE ===")

	// Get image tag from environment (default to latest if not set)
	imageTag := os.Getenv("IMAGE_TAG")
	if imageTag == "" {
		imageTag = "latest"
	}

	// Build the Go application using the ko-built image as base
	appContainer := client.Container().
		From("ghcr.io/ignaciojeria/transport-app:"+imageTag).
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

	println("✅ Transport App built successfully from source")

	// ✅ Health check - Validate that the app starts correctly
	println("=== STEP 3: Validating Transport App ===")

	// Validate Transport App with logs
	appTester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("app", appContainer.AsService()).
		WithExec([]string{
			"sh", "-c",
			`echo "Validating Transport App...";
			echo "Starting health check...";
			sleep 5;
			for i in $(seq 1 10); do
				echo "Attempt $i/10 - Checking http://app:8080/health";
				res=$(curl -s -o /dev/null -w "%{http_code}" "http://app:8080/health");
				echo "Response code: $res";
				if [ "$res" = "200" ]; then 
					echo "✅ Transport App healthy - Health check passed!"; 
					break; 
				fi;
				echo "⏳ waiting for Transport App... ($i/10)";
				sleep 2;
			done;
			if [ "$res" != "200" ]; then
				echo "❌ Transport App not ready after 10 attempts";
				echo "❌ Health check failed - App did not respond with 200";
				exit 1;
			fi;`,
		})

	appOutput, err := appTester.Stdout(ctx)
	if err != nil {
		println("❌ Health check failed - App did not start correctly")
		panic(err)
	}
	println("Transport App validation output:")
	println(appOutput)

	println("✅ Health check passed - Transport App is ready!")

	// Use the same imageTag for publishing

	// Push the Transport App container to registry
	_, err = appContainer.Publish(ctx, "ghcr.io/ignaciojeria/transport-app/transport-app-optimizator-worker:"+imageTag)
	if err != nil {
		panic(err)
	}

	println("✅ Successfully pushed Transport App to registry:")
	println("   Image: ghcr.io/ignaciojeria/transport-app/transport-app-optimizator-worker:" + imageTag)

}

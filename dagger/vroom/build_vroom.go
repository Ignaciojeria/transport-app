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

	println("🚀 Building VROOM containers from Docker registry with configurations...")

	// ✅ STEP 1: Create VROOM Optimizer from Docker registry with config
	println("=== STEP 1: Building VROOM Optimizer ===")

	optimizer := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithMountedDirectory("/conf", client.Host().Directory("./conf/optimizer")).
		WithExposedPort(3000).
		WithEntrypoint([]string{"/usr/local/bin/vroom", "--port", "3000"})

	// ✅ STEP 2: Create VROOM Planner from Docker registry with config
	println("=== STEP 2: Building VROOM Planner ===")

	planner := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithMountedDirectory("/conf", client.Host().Directory("./conf/planner")).
		WithExposedPort(3000).
		WithEntrypoint([]string{"/usr/local/bin/vroom", "--port", "3000"})

	// ✅ STEP 3: Validate VROOM containers
	println("=== STEP 3: Validating VROOM containers ===")

	// Validate Optimizer
	optimizerTester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("optimizer", optimizer.AsService()).
		WithExec([]string{
			"sh", "-c",
			`echo "Validating VROOM Optimizer...";
			sleep 5;
			for i in $(seq 1 10); do
				res=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://optimizer:3000/" -H "Content-Type: application/json" -d '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}');
				if [ "$res" = "200" ]; then 
					echo "✅ VROOM Optimizer healthy"; 
					break; 
				fi;
				echo "⏳ waiting for VROOM Optimizer... ($i/10)";
				sleep 2;
			done;
			if [ "$res" != "200" ]; then
				echo "❌ VROOM Optimizer not ready";
				exit 1;
			fi;`,
		})

	optimizerOutput, err := optimizerTester.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	println("VROOM Optimizer:", optimizerOutput)

	// Validate Planner
	plannerTester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("planner", planner.AsService()).
		WithExec([]string{
			"sh", "-c",
			`echo "Validating VROOM Planner...";
			sleep 5;
			for i in $(seq 1 10); do
				res=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://planner:3000/" -H "Content-Type: application/json" -d '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}');
				if [ "$res" = "200" ]; then 
					echo "✅ VROOM Planner healthy"; 
					break; 
				fi;
				echo "⏳ waiting for VROOM Planner... ($i/10)";
				sleep 2;
			done;
			if [ "$res" != "200" ]; then
				echo "❌ VROOM Planner not ready";
				exit 1;
			fi;`,
		})

	plannerOutput, err := plannerTester.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	println("VROOM Planner:", plannerOutput)

	// ✅ STEP 4: Push VROOM containers to registry
	println("=== STEP 4: Pushing VROOM containers to registry ===")

	// Get image tag from environment (default to latest if not set)
	imageTag := os.Getenv("IMAGE_TAG")
	if imageTag == "" {
		imageTag = uuid.New().String()
	}

	// Push VROOM Optimizer
	println("Pushing VROOM Optimizer...")
	_, err = optimizer.Publish(ctx, "ghcr.io/ignaciojeria/transport-app/vroom-optimizer:"+imageTag)
	if err != nil {
		panic(err)
	}
	println("✅ Successfully pushed VROOM Optimizer to registry:")
	println("   Image: ghcr.io/ignaciojeria/transport-app/vroom-optimizer:" + imageTag)

	// Push VROOM Planner
	println("Pushing VROOM Planner...")
	_, err = planner.Publish(ctx, "ghcr.io/ignaciojeria/transport-app/vroom-planner:"+imageTag)
	if err != nil {
		panic(err)
	}
	println("✅ Successfully pushed VROOM Planner to registry:")
	println("   Image: ghcr.io/ignaciojeria/transport-app/vroom-planner:" + imageTag)

	println("🎉 All VROOM containers pushed successfully!")
	println("   Source: ghcr.io/vroom-project/vroom-docker:v1.14.0")
	println("   Configuration: With mounted config directories")
	println("   Ports: 3000 (both optimizer and planner)")
	println("   Configs: ./conf/optimizer and ./conf/planner")
}

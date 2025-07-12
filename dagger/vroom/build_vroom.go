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

	println("🚀 Starting sequential service validation...")

	// === STEP 1: Configurar OSRM ===
	println("=== STEP 1: Configuring OSRM ===")
	osrm := client.Container().
		From("ghcr.io/ignaciojeria/transport-app/osrm-chile:latest").
		WithExposedPort(5000).
		AsService()

	// Validar OSRM inmediatamente
	osrmTester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("osrm", osrm).
		WithExec([]string{
			"sh", "-c",
			`echo "Validating OSRM...";
			for i in $(seq 1 10); do
				res=$(curl -s -o /dev/null -w "%{http_code}" "http://osrm:5000/route/v1/driving/-70.65,-33.45;-70.66,-33.46");
				if [ "$res" = "200" ]; then 
					echo "✅ OSRM healthy"; 
					exit 0; 
				fi;
				echo "⏳ waiting for OSRM... ($i/10)";
				sleep 2;
			done;
			echo "❌ OSRM not ready";
			exit 1;`,
		})

	osrmOutput, err := osrmTester.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	println("OSRM:", osrmOutput)

	// === STEP 2: Optimizer con OSRM binding ===
	println("=== STEP 2: Configuring VROOM Optimizer ===")
	optimizerContainer := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithMountedDirectory("/conf", client.Host().Directory("./conf/optimizer")).
		WithExposedPort(3000).
		WithServiceBinding("osrm", osrm)
	optimizer := optimizerContainer.AsService()

	// Validar VROOM Optimizer inmediatamente
	optimizerTester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("osrm", osrm).
		WithServiceBinding("optimizer", optimizer).
		WithExec([]string{
			"sh", "-c",
			`echo "Validating VROOM Optimizer...";
			for i in $(seq 1 10); do
				res=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://optimizer:3000/" -H "Content-Type: application/json" -d '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}');
				if [ "$res" = "200" ]; then 
					echo "✅ VROOM Optimizer healthy"; 
					exit 0; 
				fi;
				echo "⏳ waiting for VROOM Optimizer... ($i/10)";
				sleep 2;
			done;
			echo "❌ VROOM Optimizer not ready";
			exit 1;`,
		})

	optimizerOutput, err := optimizerTester.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	println("VROOM Optimizer:", optimizerOutput)

	// === STEP 3: Planner con OSRM binding ===
	println("=== STEP 3: Configuring VROOM Planner ===")
	plannerContainer := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithMountedDirectory("/conf", client.Host().Directory("./conf/planner")).
		WithExposedPort(3000).
		WithServiceBinding("osrm", osrm)
	planner := plannerContainer.AsService()

	// Validar VROOM Planner inmediatamente
	plannerTester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("osrm", osrm).
		WithServiceBinding("planner", planner).
		WithExec([]string{
			"sh", "-c",
			`echo "Validating VROOM Planner...";
			for i in $(seq 1 10); do
				res=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://planner:3000/" -H "Content-Type: application/json" -d '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}');
				if [ "$res" = "200" ]; then 
					echo "✅ VROOM Planner healthy"; 
					exit 0; 
				fi;
				echo "⏳ waiting for VROOM Planner... ($i/10)";
				sleep 2;
			done;
			echo "❌ VROOM Planner not ready";
			exit 1;`,
		})

	plannerOutput, err := plannerTester.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	println("VROOM Planner:", plannerOutput)

	println("🎉 All services are operational!")

	// === PUSH VROOM IMAGES TO REGISTRY ===
	println("=== STEP 4: Pushing VROOM containers to registry ===")

	imageTag := os.Getenv("IMAGE_TAG")
	if imageTag == "" {
		imageTag = uuid.New().String()
	}

	println("Pushing VROOM Optimizer...")
	_, err = optimizerContainer.Publish(ctx, "ghcr.io/ignaciojeria/transport-app/vroom-optimizer:"+imageTag)
	if err != nil {
		panic(err)
	}
	println("✅ Successfully pushed VROOM Optimizer to registry:")
	println("   Image: ghcr.io/ignaciojeria/transport-app/vroom-optimizer:" + imageTag)

	println("Pushing VROOM Planner...")
	_, err = plannerContainer.Publish(ctx, "ghcr.io/ignaciojeria/transport-app/vroom-planner:"+imageTag)
	if err != nil {
		panic(err)
	}
	println("✅ Successfully pushed VROOM Planner to registry:")
	println("   Image: ghcr.io/ignaciojeria/transport-app/vroom-planner:" + imageTag)

	println("🎉 All VROOM containers pushed successfully!")
}

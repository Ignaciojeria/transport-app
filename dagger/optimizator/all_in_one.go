package main

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// ✅ Leer el secreto desde variable de entorno
	natsCreds := client.SetSecret("nats-creds", os.Getenv("NATS_CONNECTION_CREDS_FILECONTENT"))

	// ✅ OSRM preprocesado
	osrm := client.Container().
		From("ghcr.io/ignaciojeria/transport-app/osrm-chile:latest").
		WithExposedPort(5000).
		AsService()

	// ✅ Optimizer con binding a OSRM
	optimizer := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithMountedDirectory("/conf", client.Host().Directory("./conf/optimizer")).
		WithExposedPort(3000).
		WithServiceBinding("osrm", osrm).
		AsService()

	// ✅ Planner con binding a OSRM
	planner := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithMountedDirectory("/conf", client.Host().Directory("./conf/planner")).
		WithExposedPort(3000).
		WithServiceBinding("osrm", osrm).
		AsService()

	// ✅ Transport App con bindings y variables
	appService := client.Container().
		From("ghcr.io/ignaciojeria/transport-app/transport-app-d0a6ffdd2b5a22c2c0423e7b340b3900@sha256:e1ee7fd378916720caaa961352cf287bede6bfc117638cbf64e47dcd0876abed").
		WithServiceBinding("planner", planner).
		WithServiceBinding("optimizer", optimizer).
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
		WithExposedPort(8080).
		AsService()

	// Sidecar Alpine para testear el healthcheck
	// Sidecar Alpine para testear el healthcheck con espera
	tester := client.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("app", appService).
		WithServiceBinding("osrm", osrm).
		WithServiceBinding("optimizer", optimizer).
		WithServiceBinding("planner", planner).
		WithExec([]string{
			"sh", "-c",
			`echo "=== Validando OSRM ===";
			for i in $(seq 1 5); do
				res=$(curl -s -o /dev/null -w "%{http_code}" "http://osrm:5000/route/v1/driving/-70.65,-33.45;-70.66,-33.46");
				if [ "$res" = "200" ]; then echo "OSRM healthy"; break; fi;
				echo "waiting for OSRM... ($i)";
				sleep 1;
			done;
			
			echo "=== Validando VROOM Optimizer ===";
			for i in $(seq 1 5); do
				res=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://optimizer:3000/" -H "Content-Type: application/json" -d '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}');
				if [ "$res" = "200" ]; then echo "VROOM Optimizer healthy"; break; fi;
				echo "waiting for VROOM Optimizer... ($i)";
				sleep 1;
			done;
			
			echo "=== Validando VROOM Planner ===";
			for i in $(seq 1 5); do
				res=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://planner:3000/" -H "Content-Type: application/json" -d '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}');
				if [ "$res" = "200" ]; then echo "VROOM Planner healthy"; break; fi;
				echo "waiting for VROOM Planner... ($i)";
				sleep 1;
			done;
			
			echo "=== Validando Transport App ===";
			for i in $(seq 1 10); do
				res=$(curl -s -o /dev/null -w "%{http_code}" http://app:8080/health);
				if [ "$res" = "200" ]; then echo "Transport App healthy"; exit 0; fi;
				echo "waiting for app... ($i)";
				sleep 1;
			done;
			echo "app not ready"; exit 1;`,
		})

	output, err := tester.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	println("Respuesta de /health:", output)

}

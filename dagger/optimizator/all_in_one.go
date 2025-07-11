package main

import (
	"context"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

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
	app := client.Container().
		From("ghcr.io/ignaciojeria/transport-app/transport-app-d0a6ffdd2b5a22c2c0423e7b340b3900:latest").
		WithServiceBinding("planner", planner).
		WithServiceBinding("optimizer", optimizer).
		WithEnvVariable("ENVIRONMENT", "prod").
		WithEnvVariable("DB_STRATEGY", "disabled").
		WithEnvVariable("OPTIMIZATION_REQUESTED_SUBSCRIPTION", "transport-app-events-optimization-requested").
		WithEnvVariable("MASTER_NODE_URL", "https://einar-main-f0820bc.d2.zuplo.dev").
		WithEnvVariable("VROOM_PLANNER_URL", "http://planner:3000").
		WithEnvVariable("VROOM_OPTIMIZER_URL", "http://optimizer:3000").
		WithEnvVariable("NATS_CONNECTION_URL", "connect.ngs.global").
		WithSecretVariable("NATS_CONNECTION_CREDS_FILECONTENT", natsCreds).
		WithExec([]string{"sleep", "10"}) // ⏳ mantenerlo vivo lo justo para probar

	// 🧪 Ejecutar y mostrar salida (debug/logs)
	_, err = app.Stdout(ctx)
	if err != nil {
		panic(err)
	}
}

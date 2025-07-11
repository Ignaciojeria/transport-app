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

	// OSRM backend con mapa de Chile
	osrm := client.Container().
		From("ghcr.io/project-osrm/osrm-backend").
		WithWorkdir("/data").
		WithExec([]string{
			"wget", "-O", "chile-latest.osm.pbf",
			"https://download.geofabrik.de/south-america/chile-latest.osm.pbf",
		}).
		WithExec([]string{
			"osrm-extract", "-p", "/opt/car.lua", "chile-latest.osm.pbf",
		}).
		WithExec([]string{
			"osrm-partition", "chile-latest.osrm",
		}).
		WithExec([]string{
			"osrm-customize", "chile-latest.osrm",
		})

	// También podrías componer la imagen final aquí si quieres:
	final := client.Container().
		From("ghcr.io/project-osrm/osrm-backend").
		WithDirectory("/data", osrm.Directory("/data")).
		WithExposedPort(5000).
		WithEntrypoint([]string{
			"osrm-routed",
			"--algorithm", "mld",
			"--max-table-size", "10000",
			"--max-viaroute-size", "2000",
			"--max-trip-size", "2000",
			"/data/chile-latest.osrm",
		})

	_, err = final.Publish(ctx, "ghcr.io/ignaciojeria/osrm-chile:latest")
	if err != nil {
		panic(err)
	}

}

package main

import (
	"context"
	"os"

	"dagger.io/dagger"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/sdk"
	_ "go.opentelemetry.io/otel/trace"
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

	// Extraer solo los archivos de datos procesados
	osrmFiles := osrm.Directory("/data")
	_, err = osrmFiles.Export(ctx, "./osrm-data")
	if err != nil {
		panic(err)
	}

	// Crear binario estático usando el mismo contenedor base de OSRM
	staticOsrm := osrm.
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "--no-cache", "python3", "py3-pip", "binutils", "patchelf", "scons", "gcc", "musl-dev", "python3-dev"}).
		WithExec([]string{"python3", "-m", "venv", "/staticx-env"}).
		WithExec([]string{"/staticx-env/bin/pip", "install", "setuptools"}).
		WithExec([]string{"/staticx-env/bin/pip", "install", "staticx"}).
		WithWorkdir("/usr/local/bin").
		WithExec([]string{
			"/staticx-env/bin/staticx", "--strip", "osrm-routed", "osrm-routed-static",
		}).
		WithExec([]string{
			"/staticx-env/bin/staticx", "--strip", "osrm-extract", "osrm-extract-static",
		}).
		WithExec([]string{
			"/staticx-env/bin/staticx", "--strip", "osrm-partition", "osrm-partition-static",
		}).
		WithExec([]string{
			"/staticx-env/bin/staticx", "--strip", "osrm-customize", "osrm-customize-static",
		})

	// Extraer solo los binarios estáticos
	staticBinaries := staticOsrm.Directory("/usr/local/bin")
	_, err = staticBinaries.Export(ctx, "./osrm-static")
	if err != nil {
		panic(err)
	}

	// Imagen final
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

	// Sincronizar el contenedor final
	_, err = final.Sync(ctx)
	if err != nil {
		panic(err)
	}

	println("Archivos OSRM extraídos a ./osrm-data")
	println("Binarios estáticos OSRM extraídos a ./osrm-static")
	println("Contenedor final sincronizado")
}

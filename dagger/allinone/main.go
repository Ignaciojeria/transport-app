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

	// VROOM Optimizer - Crear binario estático en carpeta
	println("Creando binarios estáticos de VROOM Optimizer...")
	vroomOptimizer := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "python3", "python3-pip", "python3-venv", "binutils", "patchelf", "build-essential", "gcc", "g++", "make"}).
		WithExec([]string{"python3", "-m", "venv", "/staticx-env"}).
		WithExec([]string{"/staticx-env/bin/pip", "install", "setuptools"}).
		WithExec([]string{"/staticx-env/bin/pip", "install", "staticx"}).
		WithWorkdir("/usr/local/bin").
		WithExec([]string{
			"/staticx-env/bin/staticx", "--strip", "vroom", "vroom-optimizer-static",
		})

	// Crear directorio para VROOM Optimizer con solo el binario
	vroomOptimizerDir := client.Directory().
		WithFile("vroom-optimizer-static", vroomOptimizer.File("/usr/local/bin/vroom-optimizer-static"))

	// Exportar directorio de VROOM Optimizer
	_, err = vroomOptimizerDir.Export(ctx, "./vroom-optimizer")
	if err != nil {
		panic(err)
	}

	// VROOM Planner - Crear binario estático en carpeta
	println("Creando binarios estáticos de VROOM Planner...")
	vroomPlanner := client.Container().
		From("ghcr.io/vroom-project/vroom-docker:v1.14.0").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "python3", "python3-pip", "python3-venv", "binutils", "patchelf", "build-essential", "gcc", "g++", "make"}).
		WithExec([]string{"python3", "-m", "venv", "/staticx-env"}).
		WithExec([]string{"/staticx-env/bin/pip", "install", "setuptools"}).
		WithExec([]string{"/staticx-env/bin/pip", "install", "staticx"}).
		WithWorkdir("/usr/local/bin").
		WithExec([]string{
			"/staticx-env/bin/staticx", "--strip", "vroom", "vroom-planner-static",
		})

	// Crear directorio para VROOM Planner con solo el binario
	vroomPlannerDir := client.Directory().
		WithFile("vroom-planner-static", vroomPlanner.File("/usr/local/bin/vroom-planner-static"))

	// Exportar directorio de VROOM Planner
	_, err = vroomPlannerDir.Export(ctx, "./vroom-planner")
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
	println("VROOM Optimizer extraído a ./vroom-optimizer/")
	println("VROOM Planner extraído a ./vroom-planner/")
	// Transport App - Extraer binario desde imagen de ko
	println("Extrayendo binario de Transport App desde imagen de ko...")
	koImage := client.Container().
		From("ghcr.io/ignaciojeria/transport-app/transport-app-d0a6ffdd2b5a22c2c0423e7b340b3900:latest")

	// Extraer todo el directorio /ko-app
	koAppDir := koImage.Directory("/ko-app")

	// Exportar directorio completo de Transport App
	_, err = koAppDir.Export(ctx, "./transport-app")
	if err != nil {
		panic(err)
	}

	println("Transport App extraída a ./transport-app/")

	// Clonar vroom-express desde GitHub
	println("Clonando vroom-express desde GitHub...")
	vroomExpressRepo := client.Git("https://github.com/Ignaciojeria/vroom-express.git").
		Branch("master").
		Tree()

	// Exportar el repositorio clonado a la carpeta vroom-express
	_, err = vroomExpressRepo.Export(ctx, "./vroom-express")
	if err != nil {
		panic(err)
	}

	println("vroom-express clonado a ./vroom-express/")
	println("Contenedor final sincronizado")
}

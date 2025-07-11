package main

import (
	_ "embed"
	"log"
	"os"
	_ "transport-app/app/shared/configuration"
	"transport-app/app/shared/constants"

	_ "jwks/app/adapter/in/fuegoapi"
	_ "transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}

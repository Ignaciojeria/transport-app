package main

import (
	_ "micartapro/app/shared/configuration"
	"micartapro/app/shared/constants"
	_ "embed"
	"log"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	_ "micartapro/app/shared/infrastructure/httpserver"
	_ "micartapro/app/adapter/in/fuegoapi"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	_ "embed"
	"log"
	"os"
	_ "transport-app/app/adapter/in/fuegoapi"
	_ "transport-app/app/adapter/out/tidbrepository"
	_ "transport-app/app/adapter/out/tidbrepository/table"
	_ "transport-app/app/shared/configuration"
	"transport-app/app/shared/constants"
	_ "transport-app/app/shared/infrastructure/httpserver"

	_ "transport-app/app/shared/infrastructure/database"

	_ "transport-app/app/usecase"

	_ "transport-app/app/adapter/out/fuegoapiclient"
	_ "transport-app/app/adapter/out/gcppublisher"
	_ "transport-app/app/adapter/out/restyclient/locationiq"
	_ "transport-app/app/onload"
	_ "transport-app/app/shared/infrastructure/gcppubsub"
	_ "transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	_ "transport-app/app/shared/infrastructure/httpresty"
	_ "transport-app/app/shared/infrastructure/observability"
	_ "transport-app/app/shared/infrastructure/observability/strategy"
	_ "transport-app/mocks"

	_ "transport-app/app/adapter/out/geocoding"

	_ "transport-app/app/adapter/in/natsconsumer"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	_ "transport-app/app/shared/infrastructure/storj"
	_ "transport-app/app/adapter/out/storjbucket"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}

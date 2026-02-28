package main

import (
	_ "embed"
	"log"
	_ "micartapro/app/shared/configuration"
	"micartapro/app/shared/constants"
	"os"

	_ "micartapro/app/adapter/in/fuegoapi"
	_ "micartapro/app/shared/infrastructure/gcs"
	_ "micartapro/app/shared/infrastructure/httpserver"
	_ "micartapro/app/shared/infrastructure/ngrok"

	_ "micartapro/app/adapter/in/subscriber"
	_ "micartapro/app/shared/infrastructure/ai"
	_ "micartapro/app/shared/infrastructure/eventprocessing"
	_ "micartapro/app/shared/infrastructure/eventprocessing/gcp"
	_ "micartapro/app/shared/infrastructure/observability"
	_ "micartapro/app/shared/infrastructure/observability/strategy"

	_ "micartapro/app/shared/infrastructure/auth"

	_ "micartapro/app/adapter/out/agents"
	_ "micartapro/app/adapter/out/imagegenerator"
	_ "micartapro/app/adapter/out/imageuploader"
	_ "micartapro/app/adapter/out/speechtotext"
	_ "micartapro/app/adapter/out/storage"
	_ "micartapro/app/usecase/billing"
	_ "micartapro/app/usecase/creem"
	_ "micartapro/app/usecase/menu"

	_ "micartapro/app/adapter/out/restyclient"
	_ "micartapro/app/shared/infrastructure/httpresty"
	_ "micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/ioc"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}

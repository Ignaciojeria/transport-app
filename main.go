package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	_ "transport-app/app/adapter/in/fuegoapi"
	_ "transport-app/app/adapter/out/cacherepository"
	_ "transport-app/app/adapter/out/tidbrepository"
	_ "transport-app/app/adapter/out/tidbrepository/table"
	_ "transport-app/app/shared/configuration"
	"transport-app/app/shared/constants"
	_ "transport-app/app/shared/infrastructure/httpserver"
	_ "transport-app/app/usecase/normalization/mexico"

	_ "transport-app/app/shared/infrastructure/database"

	_ "transport-app/app/shared/infrastructure/cache"

	_ "transport-app/app/usecase"

	_ "transport-app/app/adapter/in/gcpsubscription"
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
	_ "transport-app/app/shared/infrastructure/gemini"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

//go:embed .version
var version string

func main() {

	projectRoot := "./" // Asumiendo que estás en el root del proyecto
	var totalSize int64

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar carpetas ocultas como .git, .idea, etc.
		if info.IsDir() && (info.Name() == ".git" || info.Name() == ".idea" || info.Name() == "node_modules" || info.Name() == "vendor") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Tamaño total del proyecto: %.2f MB\n", float64(totalSize)/(1024*1024))

	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}

}

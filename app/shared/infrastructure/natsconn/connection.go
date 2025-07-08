package natsconn

import (
	"os"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go"
)

func init() {
	ioc.Registry(NewConn, configuration.NewNatsConfiguration)
}

func NewConn(conf configuration.NatsConfiguration) (*nats.Conn, error) {
	var credsPath string

	// Si viene por ENV, lo escribimos a un temp file
	if conf.NATS_CONNECTION_CREDS_FILECONTENT != "" {
		tmpFile, err := os.CreateTemp("", "nats-creds-*.creds")
		if err != nil {
			return nil, err
		}
		defer tmpFile.Close()

		if _, err := tmpFile.WriteString(conf.NATS_CONNECTION_CREDS_FILECONTENT); err != nil {
			return nil, err
		}
		credsPath = tmpFile.Name()
	} else {
		// Fallback a ruta de archivo
		credsPath = conf.NATS_CONNECTION_CREDS_FILEPATH
	}

	return nats.Connect(
		conf.NATS_CONNECTION_URL,
		nats.UserCredentials(credsPath),
	)
}

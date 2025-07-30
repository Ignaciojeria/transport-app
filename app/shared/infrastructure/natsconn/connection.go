package natsconn

import (
	"fmt"
	"os"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go"
)

func init() {
	ioc.Registry(NewConn, configuration.NewNatsConfiguration)
}

func NewConn(conf configuration.NatsConfiguration) (*nats.Conn, error) {
	var opts []nats.Option

	if conf.NATS_CONNECTION_CREDS_FILECONTENT != "" {
		tmpFile, err := os.CreateTemp("", "nats-creds-*.creds")
		if err != nil {
			return nil, err
		}
		defer tmpFile.Close()

		if _, err := tmpFile.WriteString(conf.NATS_CONNECTION_CREDS_FILECONTENT); err != nil {
			return nil, err
		}
		opts = append(opts, nats.UserCredentials(tmpFile.Name()))
	} else if conf.NATS_CONNECTION_CREDS_FILEPATH != "" {
		opts = append(opts, nats.UserCredentials(conf.NATS_CONNECTION_CREDS_FILEPATH))
	}

	nc, err := nats.Connect(conf.NATS_CONNECTION_URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("error conectando a NATS: %w", err)
	}

	return nc, nil
}

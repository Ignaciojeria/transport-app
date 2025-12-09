package ngrok

import (
	"context"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/httpserver"

	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	ioc.Registry(
		newTunnel,
		httpserver.New,
		configuration.NewConf)
}
func newTunnel(s httpserver.Server, conf configuration.Conf) error {
	if conf.NGROK_AUTHTOKEN == "" {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	var success bool
	go func() {
		time.Sleep(10 * time.Second)
		if !success {
			cancel()
		}
	}()
	tunel, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err == nil {
		success = true
	}
	s.SetListener(tunel)
	return err
}

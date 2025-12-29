package ngrok

import (
	"context"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	ioc.Registry(
		newTunnel,
		httpserver.New,
		configuration.NewConf,
		observability.NewObservability)
}
func newTunnel(s httpserver.Server, conf configuration.Conf, obs observability.Observability) error {
	if conf.NGROK_AUTHTOKEN == "" {
		return nil
	}
	tunel, err := ngrok.Listen(context.Background(),
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err == nil && tunel != nil {
		ngrokURL := tunel.URL()
		if ngrokURL != "" {
			obs.Logger.Info("ngrok tunnel established", "ngrok_url", ngrokURL)
		}
	}
	s.SetListener(tunel)
	return err
}

package resendcli

import (
	"fmt"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/resend/resend-go/v2"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

type ResendClient struct {
	*resend.Client
}

func NewClient(conf configuration.Conf) ResendClient {
	return ResendClient{
		Client: resend.NewClient(conf.RESEND_API_KEY),
	}
}

func (c *ResendClient) HealthCheck() error {
	domains, err := c.Domains.List()
	if err != nil {
		return fmt.Errorf("resend API no disponible: %w", err)
	}
	if len(domains.Data) == 0 {
		return fmt.Errorf("resend API respondió pero sin dominios (posible error de configuración)")
	}
	return nil
}

package resendcli

import (
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/resend/resend-go/v2"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

func NewClient(conf configuration.Conf) *resend.Client {
	return resend.NewClient(conf.RESEND_API_KEY)
}

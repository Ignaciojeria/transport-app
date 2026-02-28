package httpresty

import (
	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-resty/resty/v2"
)

func init() {
	ioc.Register(NewClient)
}
func NewClient() *resty.Client {
	cli := resty.New()
	return cli
}

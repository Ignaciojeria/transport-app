package googlemapsdk

import (
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"googlemaps.github.io/maps"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}
func NewClient(conf configuration.Conf) (*maps.Client, error) {
	return maps.NewClient(maps.WithAPIKey(conf.GOOGLE_MAPS_API_KEY))
}

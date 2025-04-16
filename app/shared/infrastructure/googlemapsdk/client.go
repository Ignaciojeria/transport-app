package googlemapsdk

import (
	"transport-app/app/shared/configuration"

	"googlemaps.github.io/maps"
)

func NewClient(conf configuration.Conf) (*maps.Client, error) {
	return maps.NewClient(maps.WithAPIKey(conf.GOOGLE_MAPS_API_KEY))
}

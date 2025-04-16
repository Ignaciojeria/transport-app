package geocoding

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/paulmach/orb"
)

type GeocodingStrategy func(context.Context, domain.AddressInfo) (orb.Point, error)

func init() {
	ioc.Registry(NewGeocodingStrategy, newGoogleGeocoding)
}

func NewGeocodingStrategy(
	conf configuration.Conf,
	google GeocodingStrategy,
	// locationIQ GeocodingStrategy, // futuro
) GeocodingStrategy {
	return func(ctx context.Context, ai domain.AddressInfo) (orb.Point, error) {
		switch conf.GEOCODING_STRATEGY {
		case "googlemaps":
			return google(ctx, ai)
		// case "libpostal":
		// 	return libpostal(ctx, ai)
		default:
			return google(ctx, ai) // fallback
		}
	}
}

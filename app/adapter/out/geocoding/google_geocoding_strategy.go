package geocoding

import (
	"context"
	"strings"
	"transport-app/app/adapter/out/cacherepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/googlemapsdk"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/paulmach/orb"
	"googlemaps.github.io/maps"
)

func init() {
	ioc.Registry(
		newGoogleGeocoding,
		googlemapsdk.NewClient,
		cacherepository.NewGeocodingCacheStrategy)
}
func newGoogleGeocoding(
	c *maps.Client,
	cache cacherepository.GeocodingCacheStrategy) GeocodingStrategy {
	return func(ctx context.Context, ai domain.AddressInfo) (orb.Point, error) {
		// Construir dirección formateada
		fullAddress := buildFullAddress(ctx, ai)

		// Construir request para geocoding
		req := &maps.GeocodingRequest{
			Address: fullAddress,
			Region:  strings.ToLower(sharedcontext.TenantCountryFromContext(ctx)),
		}

		// Llamar a Google Maps
		results, err := c.Geocode(ctx, req)
		if err != nil {
			return orb.Point{}, err
		}
		if len(results) == 0 {
			return orb.Point{}, nil
		}

		// Extraer coordenadas
		lat := results[0].Geometry.Location.Lat
		lng := results[0].Geometry.Location.Lng

		return orb.Point{lng, lat}, nil
	}
}

func buildFullAddress(ctx context.Context, ai domain.AddressInfo) string {
	full := ai.FullAddress()
	country := countries.ByName(sharedcontext.TenantCountryFromContext(ctx)).String()
	if full != "" {
		return full + ", " + country
	}
	return country
}

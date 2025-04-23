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
	cache cacherepository.GeocodingCacheStrategy,
) GeocodingStrategy {
	return func(ctx context.Context, ai domain.AddressInfo) (orb.Point, error) {
		// 1. Intentar obtener desde caché
		if cachedPoint, err := cache.Get(ctx, ai); err == nil && (cachedPoint[0] != 0 || cachedPoint[1] != 0) {
			return cachedPoint, nil
		}

		// 2. Armar dirección para geocodificar
		fullAddress := buildFullAddress(ctx, ai)

		req := &maps.GeocodingRequest{
			Address: fullAddress,
			Region:  strings.ToLower(sharedcontext.TenantCountryFromContext(ctx)),
		}

		// 3. Geocodificar usando Google Maps
		results, err := c.Geocode(ctx, req)
		if err != nil {
			return orb.Point{}, err
		}
		if len(results) == 0 {
			return orb.Point{}, nil
		}

		// 4. Guardar en caché
		lat := results[0].Geometry.Location.Lat
		lng := results[0].Geometry.Location.Lng
		point := orb.Point{lng, lat}

		ai.UpdatePoint(point)   // importante para que Save use Location
		_ = cache.Save(ctx, ai) // ignoramos error de cacheo

		return point, nil
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

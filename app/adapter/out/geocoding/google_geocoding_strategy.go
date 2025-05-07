package geocoding

import (
	"context"
	"strings"
	"transport-app/app/adapter/out/tidbrepository"
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
		tidbrepository.NewGetStoredCoordinates)
}

func newGoogleGeocoding(
	c *maps.Client,
	getCoords tidbrepository.GetStoredCoordinates,
) GeocodingStrategy {
	return func(ctx context.Context, ai domain.AddressInfo) (orb.Point, error) {
		// 1. Consultar si ya fue geocodificada previamente
		if point, err := getCoords(ctx, ai); err == nil && !point.Equal(orb.Point{}) {
			return point, nil
		}

		// 2. Geocodificar si no está en la BD
		fullAddress := buildFullAddress(ctx, ai)
		req := &maps.GeocodingRequest{
			Address: fullAddress,
			Region:  strings.ToLower(sharedcontext.TenantCountryFromContext(ctx)),
		}

		results, err := c.Geocode(ctx, req)
		if err != nil || len(results) == 0 {
			return orb.Point{}, err
		}

		point := orb.Point{
			results[0].Geometry.Location.Lng,
			results[0].Geometry.Location.Lat,
		}

		// 3. No persistimos acá — eso se hará en otra capa cuando se guarde el AddressInfo
		ai.UpdatePoint(point)
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

package cacherepository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"transport-app/app/domain"

	"github.com/paulmach/orb"
	"github.com/valkey-io/valkey-go"
)

type valkeyGeocodingCacheStrategy struct {
	c valkey.Client
}

func newValkeyGeocodingCacheStrategy(c valkey.Client) GeocodingCacheStrategy {
	return valkeyGeocodingCacheStrategy{c}
}

func (v valkeyGeocodingCacheStrategy) Save(ctx context.Context, adi domain.AddressInfo) error {
	lat := adi.Coordinates.Point.Lat()
	lng := adi.Coordinates.Point.Lon()
	concat := fmt.Sprintf("%.6f,%.6f", lat, lng)

	cmd := v.c.B().Set().Key(adi.DocID(ctx).String()).Value(concat).Build()
	return v.c.Do(ctx, cmd).Error()
}

func (v valkeyGeocodingCacheStrategy) Get(ctx context.Context, adi domain.AddressInfo) (orb.Point, error) {
	cmd := v.c.B().Get().Key(adi.DocID(ctx).String()).Build()
	resp := v.c.Do(ctx, cmd)

	val, err := resp.ToString()
	if err != nil {
		// Si la key no existe, devolvemos orb.Point{} sin error
		if strings.Contains(err.Error(), "nil") {
			return orb.Point{}, nil
		}
		return orb.Point{}, err
	}

	parts := strings.Split(val, ",")
	if len(parts) != 2 {
		return orb.Point{}, fmt.Errorf("invalid coordinate format")
	}

	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return orb.Point{}, fmt.Errorf("invalid latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return orb.Point{}, fmt.Errorf("invalid longitude: %w", err)
	}

	return orb.Point{lng, lat}, nil
}

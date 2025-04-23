package cacherepository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"transport-app/app/domain"

	"github.com/paulmach/orb"
	"github.com/redis/go-redis/v9"
)

type redisGeocodingCacheStrategy struct {
	c *redis.Client
}

func newRedisGeocodingCacheStrategy(c *redis.Client) GeocodingCacheStrategy {
	return redisGeocodingCacheStrategy{c}
}

func (r redisGeocodingCacheStrategy) Save(ctx context.Context, adi domain.AddressInfo) error {
	lat := adi.Location.Lat()
	lng := adi.Location.Lon()
	concat := fmt.Sprintf("%.6f,%.6f", lat, lng)
	return r.c.Set(ctx, adi.DocID(ctx).String(), concat, 0).Err()
}

func (r redisGeocodingCacheStrategy) Get(ctx context.Context, adi domain.AddressInfo) (orb.Point, error) {
	val, err := r.c.Get(ctx, adi.DocID(ctx).String()).Result()
	if err == redis.Nil {
		// No hay valor en cache, devolver punto vacío sin error
		return orb.Point{}, nil
	}
	if err != nil {
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

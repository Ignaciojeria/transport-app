package storj

import (
	"context"
	"time"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// Interface simplificada - Solo S3
type UplinkManager interface {
	GeneratePreSignedURL(ctx context.Context, objectKey string, ttl time.Duration) (string, error)
	GeneratePreSignedURLsBatch(ctx context.Context, objectKeys []string, ttl time.Duration) ([]string, error)
	GeneratePublicDownloadURL(ctx context.Context, objectKey string, ttl time.Duration) (string, error)
}

// Struct simplificado - Sin uplink nativo
type Uplink struct {
	// Solo necesitamos el config, ya no usamos uplink nativo
}

func init() {
	ioc.Registry(NewUplink, configuration.NewStorjConfiguration)
}

func NewUplink(env configuration.StorjConfiguration) (*Uplink, error) {
	// Validación de credenciales S3 se hace en TransportAppBucket
	return &Uplink{}, nil
}


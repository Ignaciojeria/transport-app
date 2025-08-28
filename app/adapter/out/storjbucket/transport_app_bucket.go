package storjbucket

import (
	"context"
	"fmt"
	"time"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/storj"
	"transport-app/app/shared/sharedcontext"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type TransportAppBucket struct {
	bucketName string
	s3Client   *s3.S3 // Cliente S3 para pre-signed URLs rápidas
}

func init() {
	ioc.Registry(
		NewTransportAppBucket,
		storj.NewUplink,
		configuration.NewStorjConfiguration)
}

func NewTransportAppBucket(ul *storj.Uplink, config configuration.StorjConfiguration) (storj.UplinkManager, error) {
	bucketName := "transport-app-bucket"

	// Configurar cliente S3 - REQUERIDO
	if config.STORJ_S3_ACCESS_KEY_ID == "" || config.STORJ_S3_SECRET_ACCESS_KEY == "" {
		return nil, fmt.Errorf("STORJ_S3_ACCESS_KEY_ID y STORJ_S3_SECRET_ACCESS_KEY son requeridas")
	}

	fmt.Printf("[INFO] Configurando cliente S3 para pre-signed URLs instantáneas\n")
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			config.STORJ_S3_ACCESS_KEY_ID,
			config.STORJ_S3_SECRET_ACCESS_KEY,
			"",
		),
		Endpoint:         aws.String(config.STORJ_S3_ENDPOINT),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true), // Importante para Storj
	})
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear sesión S3: %w", err)
	}

	s3Client := s3.New(sess)
	fmt.Printf("[INFO] Cliente S3 configurado exitosamente\n")

	return &TransportAppBucket{
		bucketName: bucketName,
		s3Client:   s3Client,
	}, nil
}


// ⚡ MÉTODOS S3 INSTANTÁNEOS

func (b TransportAppBucket) GeneratePreSignedURL(ctx context.Context, objectKey string, ttl time.Duration) (string, error) {
	// Crear prefijo específico por tenant
	tenantID := sharedcontext.TenantIDFromContext(ctx)
	tenantCountry := sharedcontext.TenantCountryFromContext(ctx)

	var prefixedKey string
	if tenantID != uuid.Nil && tenantCountry != "" {
		prefixedKey = fmt.Sprintf("%s-%s/%s", tenantID.String(), tenantCountry, objectKey)
	} else {
		prefixedKey = fmt.Sprintf("default/%s", objectKey)
	}

	// Generar pre-signed URL usando AWS SDK - ⚡ INSTANTÁNEO
	req, _ := b.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(prefixedKey),
	})

	urlStr, err := req.Presign(ttl)
	if err != nil {
		return "", fmt.Errorf("could not presign URL: %w", err)
	}

	return urlStr, nil
}

func (b TransportAppBucket) GeneratePreSignedURLsBatch(ctx context.Context, objectKeys []string, ttl time.Duration) ([]string, error) {
	if len(objectKeys) == 0 {
		return []string{}, nil
	}

	// Crear prefijo específico por tenant
	tenantID := sharedcontext.TenantIDFromContext(ctx)
	tenantCountry := sharedcontext.TenantCountryFromContext(ctx)

	var prefix string
	if tenantID != uuid.Nil && tenantCountry != "" {
		prefix = fmt.Sprintf("%s-%s/", tenantID.String(), tenantCountry)
	} else {
		prefix = "default/"
	}

	// Generar todas las URLs - ⚡ SUPER INSTANTÁNEO
	urls := make([]string, len(objectKeys))
	for i, objectKey := range objectKeys {
		prefixedKey := prefix + objectKey
		req, _ := b.s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(b.bucketName),
			Key:    aws.String(prefixedKey),
		})

		urlStr, err := req.Presign(ttl)
		if err != nil {
			return nil, fmt.Errorf("could not presign URL for %s: %w", objectKey, err)
		}
		urls[i] = urlStr
	}

	return urls, nil
}

func (b TransportAppBucket) GeneratePublicDownloadURL(ctx context.Context, objectKey string, ttl time.Duration) (string, error) {
	// Crear prefijo específico por tenant
	tenantID := sharedcontext.TenantIDFromContext(ctx)
	tenantCountry := sharedcontext.TenantCountryFromContext(ctx)

	var prefixedKey string
	if tenantID != uuid.Nil && tenantCountry != "" {
		prefixedKey = fmt.Sprintf("%s-%s/%s", tenantID.String(), tenantCountry, objectKey)
	} else {
		prefixedKey = fmt.Sprintf("default/%s", objectKey)
	}

	// Generar URL pública de descarga usando AWS SDK - ⚡ INSTANTÁNEO
	req, _ := b.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(prefixedKey),
	})

	urlStr, err := req.Presign(ttl)
	if err != nil {
		return "", fmt.Errorf("could not presign download URL: %w", err)
	}

	return urlStr, nil
}


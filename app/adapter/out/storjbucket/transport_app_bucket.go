package storjbucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	"transport-app/app/shared/infrastructure/storj"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"storj.io/uplink"
	"storj.io/uplink/edge"
)

type TransportAppBucket struct {
	fileExpiration  time.Duration
	sharedLinkCreds *edge.Credentials
	bucketName      string
	upLink          *storj.Uplink
}

func init() {
	ioc.Registry(
		NewTransportAppBucket,
		storj.NewUplink)
}

func NewTransportAppBucket(ul *storj.Uplink) (storj.UplinkManager, error) {
	ctx := context.Background()
	bucketName := "transport-app-bucket"

	// 🚫 Modo Edge: sin access grant, no se puede manipular bucket ni generar credenciales
	if ul.Access == nil {
		return &TransportAppBucket{
			upLink:     ul,
			bucketName: bucketName,
		}, nil
	}

	// ✅ Modo Centralizado: abrir project de forma efímera
	project, err := uplink.OpenProject(ctx, ul.Access)
	if err != nil {
		return nil, fmt.Errorf("could not open project: %w", err)
	}
	defer project.Close()

	// Crear y asegurar bucket
	_, err = project.CreateBucket(ctx, bucketName)
	if err != nil && !errors.Is(err, uplink.ErrBucketAlreadyExists) {
		return nil, fmt.Errorf("error creating bucket: %w", err)
	}
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("could not ensure bucket: %v", err)
	}

	// Crear credenciales para linksharing
	sharedLinkExpiration := 10 * time.Minute
	sharedAccess, err := ul.Access.Share(
		uplink.Permission{
			AllowDownload: true,
			NotAfter:      time.Now().Add(sharedLinkExpiration),
		},
		uplink.SharePrefix{
			Bucket: bucketName,
			Prefix: "",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not restrict access grant: %w", err)
	}

	// 🔁 Reintento de RegisterAccess
	credentials, err := retryRegisterAccess(ctx, ul.Config, sharedAccess, &edge.RegisterAccessOptions{Public: false}, 3)
	if err != nil {
		return nil, fmt.Errorf("could not register access: %w", err)
	}

	return &TransportAppBucket{
		sharedLinkCreds: credentials,
		bucketName:      bucketName,
		upLink:          ul,
	}, nil
}

func retryRegisterAccess(ctx context.Context, cfg edge.Config, access *uplink.Access, opts *edge.RegisterAccessOptions, maxRetries int) (*edge.Credentials, error) {
	var creds *edge.Credentials
	var err error
	delay := time.Second

	for i := 0; i < maxRetries; i++ {
		creds, err = cfg.RegisterAccess(ctx, access, opts)
		if err == nil {
			return creds, nil
		}

		// Solo reintentar errores de red
		if !isNetworkError(err) {
			break
		}

		time.Sleep(delay)
		delay *= 2
	}
	return nil, err
}

func isNetworkError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) ||
		strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "wsarecv") ||
		strings.Contains(err.Error(), "EOF")
}

func (b TransportAppBucket) CreatePublicSharedLink(ctx context.Context, objectKey string) (string, error) {
	// Create a public link that is served by linksharing service.
	url, err := edge.JoinShareURL("https://link.storjshare.io",
		b.sharedLinkCreds.AccessKeyID,
		b.bucketName, objectKey, nil)
	if err != nil {
		return "", fmt.Errorf("could not create a shared link: %w", err)
	}
	return url, nil
}

func (b TransportAppBucket) UploadWithToken(ctx context.Context, token string, objectKey string, data []byte) error {
	project, err := b.upLink.FromEphemeralToken(ctx, token)
	if err != nil {
		return err
	}
	defer project.Close()

	// Asegurar que el objectKey use el prefijo del tenant
	tenantID := sharedcontext.TenantIDFromContext(ctx)
	tenantCountry := sharedcontext.TenantCountryFromContext(ctx)

	var prefixedKey string
	if tenantID != uuid.Nil && tenantCountry != "" {
		prefixedKey = fmt.Sprintf("%s-%s/%s", tenantID.String(), tenantCountry, objectKey)
	} else {
		prefixedKey = fmt.Sprintf("default/%s", objectKey)
	}

	upload, err := project.UploadObject(ctx, b.bucketName, prefixedKey, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}

	_, err = io.Copy(upload, bytes.NewReader(data))
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("upload failed: %v", err)
	}

	return upload.Commit()
}

func (b TransportAppBucket) DownloadWithToken(ctx context.Context, token string, objectKey string) ([]byte, error) {
	project, err := b.upLink.FromEphemeralToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	defer project.Close()

	// Asegurar que el objectKey use el prefijo del tenant
	tenantID := sharedcontext.TenantIDFromContext(ctx)
	tenantCountry := sharedcontext.TenantCountryFromContext(ctx)

	var prefixedKey string
	if tenantID != uuid.Nil && tenantCountry != "" {
		prefixedKey = fmt.Sprintf("%s-%s/%s", tenantID.String(), tenantCountry, objectKey)
	} else {
		prefixedKey = fmt.Sprintf("default/%s", objectKey)
	}

	download, err := project.DownloadObject(ctx, b.bucketName, prefixedKey, nil)
	if err != nil {
		return nil, fmt.Errorf("could not initiate download with token: %w", err)
	}
	defer download.Close()

	var data bytes.Buffer
	_, err = io.Copy(&data, download)
	if err != nil {
		return nil, fmt.Errorf("could not read data: %w", err)
	}

	return data.Bytes(), nil
}

func (b TransportAppBucket) HasAccessGrant() bool {
	return b.upLink.Access != nil
}

func (b TransportAppBucket) GenerateEphemeralToken(ctx context.Context, ttl time.Duration, perm uplink.Permission) (string, error) {
	// Crear prefijo específico por tenant para organización y seguridad
	tenantID := sharedcontext.TenantIDFromContext(ctx)
	tenantCountry := sharedcontext.TenantCountryFromContext(ctx)

	var prefix string
	if tenantID != uuid.Nil && tenantCountry != "" {
		prefix = fmt.Sprintf("%s-%s/", tenantID.String(), tenantCountry)
	} else {
		prefix = "default/"
	}

	return b.upLink.GenerateEphemeralToken(b.bucketName, prefix, ttl, uplink.Permission{
		AllowDownload: true,
		AllowUpload:   true,
		NotAfter:      time.Now().Add(ttl),
	})
}

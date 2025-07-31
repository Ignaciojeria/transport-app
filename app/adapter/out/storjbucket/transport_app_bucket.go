package storjbucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"
	"transport-app/app/shared/infrastructure/storj"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
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

	// 🚫 Modo Edge: no puedes hacer share ni manipular buckets directamente
	if ul.Access == nil || ul.Project == nil {
		return &TransportAppBucket{
			upLink:     ul,
			bucketName: bucketName,
		}, nil
	}

	// ✅ Modo Centralizado
	sharedLinkExpiration := 10 * time.Minute
	bucketFolderName := ""

	sharedLinkRestrictedAccess, err := ul.Access.Share(
		uplink.Permission{
			AllowDownload: true,
			NotAfter:      time.Now().Add(sharedLinkExpiration),
		},
		uplink.SharePrefix{
			Bucket: bucketName,
			Prefix: bucketFolderName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not restrict access grant: %w", err)
	}

	credentials, err := ul.Config.RegisterAccess(ctx, sharedLinkRestrictedAccess, &edge.RegisterAccessOptions{Public: false})
	if err != nil {
		return nil, fmt.Errorf("could not register access: %w", err)
	}

	_, err = ul.Project.CreateBucket(ctx, bucketName)
	if err != nil && !errors.Is(err, uplink.ErrBucketAlreadyExists) {
		return nil, fmt.Errorf("error creating bucket: %w", err)
	}

	_, err = ul.Project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("could not ensure bucket: %v", err)
	}

	return &TransportAppBucket{
		sharedLinkCreds: credentials,
		bucketName:      bucketName,
		upLink:          ul,
	}, nil
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

	upload, err := project.UploadObject(ctx, b.bucketName, objectKey, nil)
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

	download, err := project.DownloadObject(ctx, b.bucketName, objectKey, nil)
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

func (b TransportAppBucket) GenerateEphemeralToken(ctx context.Context, ttl time.Duration, perm uplink.Permission) (string, error) {
	prefix := fmt.Sprintf("%s-%s",
		sharedcontext.TenantIDFromContext(ctx),
		sharedcontext.TenantCountryFromContext(ctx))
	return b.upLink.GenerateEphemeralToken(b.bucketName, prefix, ttl, perm)
}

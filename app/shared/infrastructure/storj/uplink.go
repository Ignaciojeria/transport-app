package storj

import (
	"context"
	"fmt"
	"time"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"storj.io/uplink"
	"storj.io/uplink/edge"
)

type UplinkManager interface {
	CreatePublicSharedLink(ctx context.Context, objectKey string) (string, error)
	UploadWithToken(ctx context.Context, token string, objectKey string, data []byte) error
	DownloadWithToken(ctx context.Context, token string, objectKey string) ([]byte, error)
}

type Uplink struct {
	Access  *uplink.Access
	Project *uplink.Project
	Config  edge.Config
}

func init() {
	ioc.Registry(NewUplink, configuration.NewStorjConfiguration)
}

func NewUplink(env configuration.StorjConfiguration) (*Uplink, error) {
	ctx := context.Background()

	// Modo Edge (sin access grant)
	if env.STORJ_ACCESS_GRANT == "" {
		return &Uplink{
			Access:  nil,
			Project: nil,
			Config: edge.Config{
				AuthServiceAddress: "auth.storjshare.io:7777",
			},
		}, nil
	}

	// Modo Centralizado (con access grant)
	access, err := uplink.ParseAccess(env.STORJ_ACCESS_GRANT)
	if err != nil {
		return nil, fmt.Errorf("invalid STORJ_ACCESS_GRANT: %w", err)
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, fmt.Errorf("could not open project: %v", err)
	}

	return &Uplink{
		Access:  access,
		Project: project,
		Config: edge.Config{
			AuthServiceAddress: "auth.storjshare.io:7777",
		},
	}, nil
}

func (u *Uplink) FromEphemeralToken(ctx context.Context, ephemeralToken string) (*uplink.Project, error) {
	access, err := uplink.ParseAccess(ephemeralToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access grant: %w", err)
	}
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, fmt.Errorf("could not open project from token: %w", err)
	}
	return project, nil
}

func (u *Uplink) GenerateEphemeralToken(bucket, prefix string, ttl time.Duration, perm uplink.Permission) (string, error) {
	if u.Access == nil {
		return "", fmt.Errorf("uplink Access not initialized")
	}

	// Define permisos con expiración
	if perm.NotAfter.IsZero() {
		perm.NotAfter = time.Now().Add(ttl)
	}

	sharedAccess, err := u.Access.Share(perm, uplink.SharePrefix{
		Bucket: bucket,
		Prefix: prefix,
	})
	if err != nil {
		return "", fmt.Errorf("could not generate ephemeral access: %w", err)
	}

	grant, err := sharedAccess.Serialize()
	if err != nil {
		return "", fmt.Errorf("could not serialize access: %w", err)
	}

	return grant, nil
}

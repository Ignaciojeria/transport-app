package usecase

import (
	"context"
	"transport-app/app/adapter/out/email"
	"transport-app/app/adapter/out/firebaseauth"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/jwt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type Register func(context.Context, domain.TenantAccount) error

func init() {
	ioc.Registry(
		NewRegister,
		firebaseauth.NewRegister,
		tidbrepository.NewUpsertAccount,
		jwt.NewJWTServiceFromConfig,
		NewCreateTenantAccount,
		tidbrepository.NewFindDefaultTenantByEmail,
		NewCreateClientCredentials,
		email.NewSendClientCredentialsEmail)
}
func NewRegister(
	register firebaseauth.Register,
	upsertAccount tidbrepository.UpsertAccount,
	jwtService *jwt.JWTService,
	createTenantAccount CreateTenantAccount,
	findDefaultTenantByEmail tidbrepository.FindDefaultTenantByEmail,
	createClientCredentials CreateClientCredentials,
	sendClientCredentialsEmail email.SendClientCredentialsEmail,
) Register {
	return func(ctx context.Context, input domain.TenantAccount) error {
		err := upsertAccount(ctx, domain.Account{
			Email: input.Account.Email,
		})
		if err != nil {
			return err
		}
		/*
			tenant := ""

				token, err := jwtService.GenerateToken(
					input.Account.Email,
					[]string{"tenants:read"},
					map[string]string{},
					tenant,
					"zuplo-gateway",
					60)
				if err != nil {
					return err
				}
		*/

		existingTenantAccount, err := findDefaultTenantByEmail(ctx, input.Account.Email)
		if err != nil {
			return err
		}

		if (domain.TenantAccount{}) != existingTenantAccount {
			return nil
		}

		tenantID := uuid.New()

		err = createTenantAccount(ctx, domain.TenantAccount{
			Tenant: domain.Tenant{
				ID:      tenantID,
				Name:    "Default Tenant",
				Country: input.Tenant.Country,
			},
			Status:  "pending",
			Account: input.Account,
			Role:    input.Role,
		})

		if err != nil {
			return err
		}

		// Generar client credentials para el tenant
		scopes := []string{
			"orders:read",
			"orders:write",
			"routes:read",
			"routes:write",
			"nodes:read",
			"nodes:write",
			"optimization:read",
			"optimization:write",
		}

		clientCredentials, err := createClientCredentials(ctx, tenantID, input.Tenant.Country, scopes)
		if err != nil {
			return err
		}

		// Enviar las credenciales por email
		err = sendClientCredentialsEmail(ctx, input.Account.Email, clientCredentials)
		if err != nil {
			return err
		}

		return nil
	}
}

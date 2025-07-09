package usecase

import (
	"context"
	"transport-app/app/adapter/out/firebaseauth"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/jwt"
	resendcli "transport-app/app/shared/infrastructure/resendcli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
)

type Register func(context.Context, domain.TenantAccount) error

func init() {
	ioc.Registry(
		NewRegister,
		firebaseauth.NewRegister,
		tidbrepository.NewUpsertAccount,
		jwt.NewJWTServiceFromConfig,
		resendcli.NewClient,
		NewCreateTenantAccount,
		tidbrepository.NewFindDefaultTenantByEmail)
}
func NewRegister(
	register firebaseauth.Register,
	upsertAccount tidbrepository.UpsertAccount,
	jwtService *jwt.JWTService,
	resendClient *resend.Client,
	createTenantAccount CreateTenantAccount,
	findDefaultTenantByEmail tidbrepository.FindDefaultTenantByEmail,
) Register {
	return func(ctx context.Context, input domain.TenantAccount) error {
		err := upsertAccount(ctx, domain.Account{
			Email: input.Account.Email,
		})
		if err != nil {
			return err
		}
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

		existingTenantAccount, err := findDefaultTenantByEmail(ctx, input.Account.Email)
		if err != nil {
			return err
		}

		// Solo crear el tenant si no existe uno por defecto
		if (domain.TenantAccount{}) == existingTenantAccount {
			err = createTenantAccount(ctx, domain.TenantAccount{
				Tenant: domain.Tenant{
					ID:      uuid.New(),
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
		}

		params := &resend.SendEmailRequest{
			From:    "Transport App • Test Client Credentials <onboarding@resend.dev>",
			To:      []string{input.Account.Email},
			Html:    "<strong>" + token + "</strong>",
			Subject: "transport app - test client credentials",
			Cc:      []string{},
			Bcc:     []string{},
			ReplyTo: "",
		}

		_, err = resendClient.Emails.Send(params)
		if err != nil {
			return err
		}
		return nil
	}
}

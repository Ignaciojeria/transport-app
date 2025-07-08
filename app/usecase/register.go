package usecase

import (
	"context"
	"transport-app/app/adapter/out/firebaseauth"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/jwt"
	resendcli "transport-app/app/shared/infrastructure/resendcli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/resend/resend-go/v2"
)

type Register func(ctx context.Context, input domain.UserCredentials) error

func init() {
	ioc.Registry(
		NewRegister,
		firebaseauth.NewRegister,
		tidbrepository.NewUpsertAccount,
		jwt.NewJWTServiceFromConfig,
		resendcli.NewClient)
}
func NewRegister(
	register firebaseauth.Register,
	upsertAccount tidbrepository.UpsertAccount,
	jwtService *jwt.JWTService,
	resendClient *resend.Client,
) Register {
	return func(ctx context.Context, input domain.UserCredentials) error {
		err := upsertAccount(ctx, domain.Account{
			Email: input.Email,
		})
		if err != nil {
			return err
		}
		tenant := ""
		token, err := jwtService.GenerateToken(
			input.Email,
			[]string{"onboarding:access"},
			map[string]string{},
			tenant,
			"zuplo-gateway",
			60)
		if err != nil {
			return err
		}

		params := &resend.SendEmailRequest{
			From:    "Transport App • Magic Link <onboarding@resend.dev>",
			To:      []string{input.Email},
			Html:    "<strong>" + token + "</strong>",
			Subject: "transport app - magic link",
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

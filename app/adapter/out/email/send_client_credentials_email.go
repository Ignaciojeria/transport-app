package email

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"transport-app/app/domain"
	resendcli "transport-app/app/shared/infrastructure/resendcli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/resend/resend-go/v2"
)

//go:embed client_credentials_email.html
var emailTemplates embed.FS

type SendClientCredentialsEmail func(ctx context.Context, email string, credentials domain.ClientCredentials) error

func init() {
	ioc.Registry(
		NewSendClientCredentialsEmail,
		resendcli.NewClient,
	)
}

type EmailData struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func NewSendClientCredentialsEmail(resendClient resendcli.ResendClient) SendClientCredentialsEmail {
	return func(ctx context.Context, email string, credentials domain.ClientCredentials) error {
		// Cargar el template HTML desde el embed
		tmpl, err := template.ParseFS(emailTemplates, "client_credentials_email.html")
		if err != nil {
			return err
		}

		// Preparar los datos para el template
		data := EmailData{
			ClientID:     credentials.ClientID,
			ClientSecret: credentials.ClientSecret, // Desencriptado automáticamente
			Scopes:       credentials.AllowedScopes,
		}

		// Renderizar el template
		var body bytes.Buffer
		err = tmpl.Execute(&body, data)
		if err != nil {
			return err
		}

		// Enviar el email
		params := &resend.SendEmailRequest{
			From:    "Transport App <onboarding@transport-app.com>",
			To:      []string{email},
			Html:    body.String(),
			Subject: "🚚 Transport App - Tus Credenciales de Cliente",
		}

		_, err = resendClient.Emails.Send(params)
		return err
	}
}

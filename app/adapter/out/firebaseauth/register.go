package firebaseauth

import (
	"context"
	"log"

	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/firebaseadminsdk"

	"firebase.google.com/go/v4/auth"

	firebase "firebase.google.com/go/v4"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Register func(context.Context, domain.UserCredentials) error

func init() {
	ioc.Registry(
		NewRegister,
		firebaseadminsdk.NewFirebaseAdmin)
}

func NewRegister(app *firebase.App) Register {
	return func(ctx context.Context, uc domain.UserCredentials) error {
		client, err := app.Auth(ctx)
		if err != nil {
			log.Printf("Error obteniendo el cliente de Auth: %v", err)
			return err
		}

		// Crear usuario en Firebase Authentication
		params := (&auth.UserToCreate{}).
			Email(uc.Email).
			Password(uc.Password).
			Disabled(false)

		userRecord, err := client.CreateUser(ctx, params)
		if err != nil {
			log.Printf("Error al registrar usuario en Firebase: %v", err)
			return err
		}

		log.Printf("Usuario creado con éxito: %s", userRecord.UID)
		return nil
	}
}

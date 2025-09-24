package fuegoapi

import (
	"fmt"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(authGoogleCallback, httpserver.New)
}
func authGoogleCallback(s httpserver.Server) {
	fuego.Get(s.Manager, "/auth/google/callback",
		func(c fuego.ContextNoBody) (any, error) {
			// Obtener parámetros de la URL
			code := c.QueryParam("code")
			state := c.QueryParam("state")
			errorParam := c.QueryParam("error")

			fmt.Printf("Google OAuth callback - Code: %s, State: %s, Error: %s\n", code, state, errorParam)

			if errorParam != "" {
				// Usuario canceló o hubo error
				return map[string]string{
					"error":    errorParam,
					"redirect": "https://auth.transport-app.com/auth/error?error=" + errorParam,
				}, nil
			}

			if code == "" {
				return map[string]string{
					"error":    "missing_code",
					"redirect": "https://auth.transport-app.com/auth/error?error=missing_code",
				}, nil
			}

			// TODO: Implementar intercambio de código por token
			// TODO: Validar usuario y generar JWT
			// TODO: Redirigir a frontend con token

			return map[string]string{
				"status":   "success",
				"code":     code,
				"redirect": "https://auth.transport-app.com/auth/success?token=temp_token",
			}, nil
		}, option.Summary("Google OAuth Callback"))
}

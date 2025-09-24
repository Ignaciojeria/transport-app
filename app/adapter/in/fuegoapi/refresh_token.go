package fuegoapi

import (
	"io"
	"net/url"
	"time"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/jwt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(refreshToken, httpserver.New, jwt.NewJWTServiceFromConfig)
}

func refreshToken(s httpserver.Server, jwtService *jwt.JWTService) {
	fuego.Post(s.Manager, "/refresh",
		func(c fuego.ContextNoBody) (response.TokenResponse, error) {
			// Leer el body de la request
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "error leyendo body",
					Detail: "Error al leer el body de la request",
					Status: 400,
				}
			}

			// Parsear los datos application/x-www-form-urlencoded
			values, err := url.ParseQuery(string(body))
			if err != nil {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "error parseando formulario",
					Detail: "Error al procesar los datos del formulario",
					Status: 400,
				}
			}

			// Obtener valores del formulario
			grantType := values.Get("grant_type")
			refreshTokenValue := values.Get("refresh_token")

			// Validar que el grant_type sea refresh_token
			if grantType != "refresh_token" {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "grant_type inválido",
					Detail: "El grant_type debe ser 'refresh_token'",
					Status: 400,
				}
			}

			// Validar que el refresh_token esté presente
			if refreshTokenValue == "" {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "refresh_token requerido",
					Detail: "El campo 'refresh_token' es obligatorio",
					Status: 400,
				}
			}

			// Validar el refresh token
			claims, err := jwtService.ValidateToken(refreshTokenValue)
			if err != nil {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "refresh_token inválido",
					Detail: "El refresh_token proporcionado no es válido o ha expirado",
					Status: 401,
				}
			}

			// Verificar que el token tenga el scope "refresh"
			hasRefreshScope := false
			for _, scope := range claims.Scopes {
				if scope == "refresh" {
					hasRefreshScope = true
					break
				}
			}
			if !hasRefreshScope {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "token inválido",
					Detail: "El token proporcionado no es un refresh token válido",
					Status: 401,
				}
			}

			// Obtener el audience del token original
			var audience string
			if len(claims.Audience) > 0 {
				audience = claims.Audience[0]
			}

			// Determinar el scope para el nuevo access token (usar el primero disponible que no sea "refresh")
			var accessScope string
			for _, scope := range claims.Scopes {
				if scope != "refresh" {
					accessScope = scope
					break
				}
			}
			if accessScope == "" {
				accessScope = "read" // scope por defecto si no se encuentra otro
			}

			// Generar nuevo access token con 60 minutos de duración
			newAccessToken, err := jwtService.GenerateToken(
				claims.Sub,
				[]string{accessScope},
				claims.Context,
				claims.Tenant,
				audience,
				60, // 60 minutos
			)
			if err != nil {
				return response.TokenResponse{}, err
			}

			// Generar nuevo refresh token con 7 días de duración
			newRefreshToken, err := jwtService.GenerateToken(
				claims.Sub,
				[]string{"refresh"},
				claims.Context,
				claims.Tenant,
				audience,
				10080, // 7 días
			)
			if err != nil {
				return response.TokenResponse{}, err
			}

			// Extraer la expiración del access token para calcular ExpiresIn
			newClaims, err := jwtService.ValidateToken(newAccessToken)
			if err != nil {
				return response.TokenResponse{}, err
			}

			return response.TokenResponse{
				AccessToken:  newAccessToken,
				TokenType:    "Bearer",
				ExpiresIn:    int(newClaims.ExpiresAt.Unix() - time.Now().Unix()),
				RefreshToken: newRefreshToken,
				Scope:        accessScope,
			}, nil
		}, option.Summary("refresh token"), option.Tags("jwt"))
}

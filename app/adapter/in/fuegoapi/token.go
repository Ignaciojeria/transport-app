package fuegoapi

import (
	"io"
	"net/url"
	"time"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/jwt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type TokenRequest struct {
	GrantType    string `json:"grant_type" example:"client_credentials"`
	ClientID     string `json:"client_id" example:"abc123"`
	ClientSecret string `json:"client_secret" example:"secret456"`
	Scope        string `json:"scope" example:"read:orders"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
	Scope       string `json:"scope" example:"read:orders"`
}

func init() {
	ioc.Registry(token, httpserver.New, jwt.NewJWTServiceFromConfig)
}

func token(s httpserver.Server, jwtService *jwt.JWTService) {
	fuego.Post(s.Manager, "/token",
		func(c fuego.ContextWithBody[TokenRequest]) (TokenResponse, error) {
			// Leer el body de la request
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return TokenResponse{}, fuego.HTTPError{
					Title:  "error leyendo body",
					Detail: "Error al leer el body de la request",
					Status: 400,
				}
			}

			// Parsear los datos application/x-www-form-urlencoded
			values, err := url.ParseQuery(string(body))
			if err != nil {
				return TokenResponse{}, fuego.HTTPError{
					Title:  "error parseando formulario",
					Detail: "Error al procesar los datos del formulario",
					Status: 400,
				}
			}

			// Obtener valores del formulario
			grantType := values.Get("grant_type")
			clientID := values.Get("client_id")
			clientSecret := values.Get("client_secret")
			scope := values.Get("scope")

			// Validar que el client_id esté presente
			if clientID == "" {
				return TokenResponse{}, fuego.HTTPError{
					Title:  "client_id requerido",
					Detail: "El campo 'client_id' es obligatorio",
					Status: 400,
				}
			}

			// Validar que el client_secret esté presente
			if clientSecret == "" {
				return TokenResponse{}, fuego.HTTPError{
					Title:  "client_secret requerido",
					Detail: "El campo 'client_secret' es obligatorio",
					Status: 400,
				}
			}

			// Validar que el grant_type sea client_credentials
			if grantType != "client_credentials" {
				return TokenResponse{}, fuego.HTTPError{
					Title:  "grant_type inválido",
					Detail: "El grant_type debe ser 'client_credentials'",
					Status: 400,
				}
			}

			// Generar token JWT usando client_id como sub y zuplo-gateway como aud
			token, err := jwtService.GenerateToken(
				clientID,            // sub = client_id
				[]string{scope},     // scopes
				map[string]string{}, // context vacío
				"",                  // tenant vacío
				"zuplo-gateway",     // aud en duro
				60,                  // 60 minutos de expiración
			)
			if err != nil {
				return TokenResponse{}, err
			}

			// Extraer la expiración del token para calcular ExpiresIn
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return TokenResponse{}, err
			}

			return TokenResponse{
				AccessToken: token,
				TokenType:   "Bearer",
				ExpiresIn:   int(claims.ExpiresAt.Unix() - time.Now().Unix()), // Calcular tiempo restante
				Scope:       scope,
			}, nil
		}, option.Summary("token"), option.Tags("jwt"))
}

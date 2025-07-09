package fuegoapi

import (
	"io"
	"net/url"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/jwt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(token, httpserver.New, jwt.NewJWTServiceFromConfig, tidbrepository.NewFindClientCredentialsByClientID)
}

func token(s httpserver.Server, jwtService *jwt.JWTService, findClientCredentials tidbrepository.FindClientCredentialsByClientID) {
	fuego.Post(s.Manager, "/token",
		func(c fuego.ContextWithBody[request.TokenRequest]) (response.TokenResponse, error) {
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
			clientID := values.Get("client_id")
			clientSecret := values.Get("client_secret")
			scope := values.Get("scope")

			// Validar que el client_id esté presente
			if clientID == "" {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "client_id requerido",
					Detail: "El campo 'client_id' es obligatorio",
					Status: 400,
				}
			}

			// Validar que el client_secret esté presente
			if clientSecret == "" {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "client_secret requerido",
					Detail: "El campo 'client_secret' es obligatorio",
					Status: 400,
				}
			}

			// Validar que el grant_type sea client_credentials
			if grantType != "client_credentials" {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "grant_type inválido",
					Detail: "El grant_type debe ser 'client_credentials'",
					Status: 400,
				}
			}

			// Buscar las credenciales del cliente en la base de datos
			credentials, err := findClientCredentials(c.Context(), clientID)
			if err != nil {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "credenciales inválidas",
					Detail: "Las credenciales del cliente no son válidas",
					Status: 401,
				}
			}

			// Validar que las credenciales estén activas
			if !credentials.IsValid() {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "credenciales inactivas",
					Detail: "Las credenciales del cliente no están activas o han expirado",
					Status: 401,
				}
			}

			// Validar que el client_secret coincida
			if credentials.ClientSecret != clientSecret {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "credenciales inválidas",
					Detail: "El client_secret proporcionado no es válido",
					Status: 401,
				}
			}

			// Validar que el scope esté permitido (si se proporciona)
			if scope != "" && !credentials.HasScope(scope) {
				return response.TokenResponse{}, fuego.HTTPError{
					Title:  "scope no autorizado",
					Detail: "El scope solicitado no está autorizado para estas credenciales",
					Status: 403,
				}
			}

			// Si no se proporciona scope, usar el primer scope disponible
			if scope == "" && len(credentials.AllowedScopes) > 0 {
				scope = credentials.AllowedScopes[0]
			}

			// Convertir el tenant ID a string para incluirlo en el token
			tenantID := credentials.TenantID.String()

			// Obtener el país del tenant
			tenantCountry := credentials.TenantCountry.String()

			// Generar token JWT usando client_id como sub, tenant ID, tenant country y zuplo-gateway como aud
			token, err := jwtService.GenerateToken(
				clientID,                   // sub = client_id
				[]string{scope},            // scopes
				map[string]string{},        // context vacío
				tenantID+"-"+tenantCountry, // tenant = tenant ID + country
				"zuplo-gateway",            // aud en duro
				60,                         // 60 minutos de expiración
			)
			if err != nil {
				return response.TokenResponse{}, err
			}

			// Extraer la expiración del token para calcular ExpiresIn
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return response.TokenResponse{}, err
			}

			return response.TokenResponse{
				AccessToken: token,
				TokenType:   "Bearer",
				ExpiresIn:   int(claims.ExpiresAt.Unix() - time.Now().Unix()), // Calcular tiempo restante
				Scope:       scope,
			}, nil
		}, option.Summary("auth token"), option.Tags("jwt"))
}

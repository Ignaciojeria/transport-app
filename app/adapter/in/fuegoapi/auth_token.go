package fuegoapi

/*
import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/jwt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(authToken, httpserver.New, jwt.NewJWTServiceFromConfig)
}

func authToken(s httpserver.Server, jwtService *jwt.JWTService) {
	fuego.Post(s.Manager, "/auth/token",
		func(c fuego.ContextWithBody[request.GenerateTokenRequest]) (response.GenerateTokenResponse, error) {
			// Obtener el body de la request
			req, err := c.Body()
			if err != nil {
				return response.GenerateTokenResponse{}, err
			}

			// Validar que el audience esté presente
			if req.Audience == "" {
				return response.GenerateTokenResponse{}, fuego.HTTPError{
					Title:  "audience requerido",
					Detail: "El campo 'aud' es obligatorio",
					Status: 400,
				}
			}

			// Obtener el tenant del header
			tenant := c.Request().Header.Get("tenant")
			if tenant == "" {
				return response.GenerateTokenResponse{}, fuego.HTTPError{
					Title:  "tenant header requerido",
					Detail: "El header 'tenant' es obligatorio",
					Status: 400,
				}
			}

			// Generar token JWT
			token, err := jwtService.GenerateToken(
				req.Sub,
				req.Scopes,
				req.Context,
				tenant,
				req.Audience,
				60, // 60 minutos de expiración
			)
			if err != nil {
				return response.GenerateTokenResponse{}, err
			}

			// Extraer la expiración del token para la respuesta
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return response.GenerateTokenResponse{}, err
			}

			return response.GenerateTokenResponse{
				Token:     token,
				ExpiresAt: claims.ExpiresAt.Unix(),
				TokenType: "Bearer",
			}, nil
		},
		option.Header("tenant", "api tenant (required only for local development)", param.Required()),
		option.Header("X-Access-Token", "api access token"),
		option.Summary("auth token"),
		option.Tags("jwt"))
}
*/

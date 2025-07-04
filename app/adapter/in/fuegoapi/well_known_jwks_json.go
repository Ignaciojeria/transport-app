package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/jwt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(validateToken, httpserver.New, jwt.NewJWTServiceFromConfig)
}

func validateToken(s httpserver.Server, jwtService *jwt.JWTService) {
	// JWKS endpoint para OpenID Connect - ESENCIAL para Zuplo
	fuego.Get(s.Manager, "/.well-known/jwks.json",
		func(c fuego.ContextNoBody) (*jwt.JWKS, error) {
			jwks, err := jwtService.GetJWKS()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "JWKS not available",
					Detail: "No public key configured for JWKS",
					Status: 500,
				}
			}
			return jwks, nil
		},
		option.Summary("/.well-known/jwks.json"),
		option.Tags("openid"))
}

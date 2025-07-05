package fuegoapi

import (
	"encoding/json"
	"net/http"
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
	fuego.GetStd(s.Manager, "/.well-known/jwks.json",
		func(w http.ResponseWriter, r *http.Request) {
			jwks, err := jwtService.GetJWKS()
			if err != nil {
				http.Error(w, "JWKS not available", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jwks)
		},
		option.Summary("/.well-known/jwks.json"),
		option.Tags("openid"))
}

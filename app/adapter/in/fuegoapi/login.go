package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		login,
		httpserver.New,
		usecase.NewLogin)
}
func login(
	s httpserver.Server,
	login usecase.Login) {
	fuego.Post(s.Manager, "/login",
		func(c fuego.ContextWithBody[request.LoginRequest]) (domain.ProviderToken, error) {
			requestBody, err := c.Body()
			if err != nil {
				return domain.ProviderToken{}, err
			}
			return login(c.Context(), domain.UserCredentials{
				Email:    requestBody.Email,
				Password: requestBody.Password,
			})
		},
		option.Tags(tagAuthentication),
		option.Summary("login"))
}

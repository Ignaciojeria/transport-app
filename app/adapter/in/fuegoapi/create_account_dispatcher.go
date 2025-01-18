package fuegoapi

import (
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		createAccount,
		httpserver.New,
		usecase.NewCreateAccount)
}
func createAccount(
	s httpserver.Server,
	createAccount usecase.CreateAccount) {
	fuego.Post(s.Manager, "/account/dispatcher",
		func(c fuego.ContextWithBody[request.CreateDispatcherRequest]) (response.CreateAccountResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.CreateAccountResponse{}, err
			}
			acc := requestBody.Map()
			acc.Organization.Key = c.Header("organization-key")
			acc.Organization.Country = countries.ByName(c.Header("country"))
			_, err = createAccount(c.Context(), acc)
			if err != nil {
				return response.CreateAccountResponse{}, fuego.HTTPError{
					Title:  "error creating account",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			return response.CreateAccountResponse{
				Message: "account created",
			}, nil
		}, option.Summary("createAccountDispatcher"),
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
		option.Tags(tagAccounts),
	)
}

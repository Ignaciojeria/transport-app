package fuegoapi

import (
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
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
	fuego.Post(s.Manager, "/account",
		func(c fuego.ContextWithBody[model.CreateAccountRequest]) (model.CreateAccountResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return model.CreateAccountResponse{}, err
			}

			_, err = createAccount(c.Context(), requestBody.Map())
			if err != nil {
				return model.CreateAccountResponse{}, fuego.HTTPError{
					Title:  "error creating account",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			return model.CreateAccountResponse{
				Message: "account created",
			}, nil
		}, option.Summary("createAccount"))
}

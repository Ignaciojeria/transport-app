package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(register, httpserver.New, usecase.NewRegister)
}
func register(s httpserver.Server, register usecase.Register) {
	fuego.Post(s.Manager, "/register",
		func(c fuego.ContextWithBody[request.RegisterRequest]) (response.RegisterResponse, error) {
			req, err := c.Body()
			if err != nil {
				return response.RegisterResponse{}, err
			}
			err = register(c.Context(), domain.UserCredentials{
				Email:    req.Email,
				Password: req.Password,
			})
			if err != nil {
				return response.RegisterResponse{}, fuego.InternalServerError{
					Err:    err,
					Detail: err.Error(),
				}
			}
			return response.RegisterResponse{
				Message: "user registered",
			}, nil
		},
		option.Tags(tagAuthentication),
		option.Summary("register"))
}

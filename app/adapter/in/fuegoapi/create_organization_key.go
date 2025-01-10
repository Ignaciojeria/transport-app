package fuegoapi

import (
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/domain"
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
		createOrganization,
		httpserver.New,
		usecase.NewCreateOrganizationKey)
}
func createOrganization(s httpserver.Server, createOrg usecase.CreateOrganizationKey) {
	fuego.Post(s.Manager, "/organization-key",
		func(c fuego.ContextWithBody[request.CreateOrganizationKeyRequest]) (response.CreateOrganizationKeyResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.CreateOrganizationKeyResponse{}, err
			}
			org := domain.Organization{
				Name:    requestBody.Email,
				Email:   requestBody.Email,
				Country: countries.ByName(c.Header("country")),
			}
			org, err = createOrg(c.Context(), org)
			if err != nil {
				return response.CreateOrganizationKeyResponse{}, fuego.HTTPError{
					Title:  "error creating organization",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			return response.CreateOrganizationKeyResponse{
				OrganizationKey: org.Key,
				Message:         "Organization key created successfully. Please save it securely as it will be required for API authentication.",
			}, nil
		},
		option.Summary("createOrganizationKey"),
		option.Header("country", "api country", param.Required()))
}

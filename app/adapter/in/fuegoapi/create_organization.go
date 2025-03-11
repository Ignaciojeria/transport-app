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
)

func init() {
	ioc.Registry(
		createOrganization,
		httpserver.New,
		usecase.NewCreateOrganization)
}
func createOrganization(s httpserver.Server, createOrg usecase.CreateOrganization) {
	fuego.Post(s.Manager, "/organizations",
		func(c fuego.ContextWithBody[request.CreateOrganizationRequest]) (response.CreateOrganizationResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.CreateOrganizationResponse{}, err
			}
			org := domain.Organization{
				Name:    requestBody.Name,
				Email:   requestBody.Email,
				Country: countries.ByName(c.Header("country")),
			}
			org, err = createOrg(c.Context(), org)
			if err != nil {
				return response.CreateOrganizationResponse{}, fuego.HTTPError{
					Title:  "error creating organization",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			return response.CreateOrganizationResponse{
				OrganizationKey: org.GetOrgKey(),
				Message:         "Organization created successfully",
			}, nil
		},
		option.Summary("createOrganizationKey"),
		option.Tags(tagOrganizations),
	)
}

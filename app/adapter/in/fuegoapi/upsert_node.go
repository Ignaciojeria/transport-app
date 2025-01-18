package fuegoapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		upsertNode, httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		tidbrepository.NewSaveNodeOutbox,
	)
}
func upsertNode(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	outbox tidbrepository.SaveNodeOutbox) {
	fuego.Post(s.Manager, "/node",
		func(c fuego.ContextWithBody[request.UpsertNodeRequest]) (response.UpsertNodeResponse, error) {

			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertNodeResponse{}, err
			}

			organization := domain.Organization{
				Key:     c.Header("organization-key"),
				Country: countries.ByName(c.Header("country")),
			}

			org, err := ensureOrg(c.Context(), organization)
			if err != nil {
				return response.UpsertNodeResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			requestBodyBytes, _ := json.Marshal(requestBody)
			orgIDString := strconv.FormatInt(org.OrganizationCountryID, 10)
			if _, err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":            "node",
					"eventType":             "nodeSubmitted",
					"country":               countries.ByName(c.Header("country")).Alpha2(),
					"organizationCountryID": orgIDString,
					"consumer":              c.Header("consumer"),
					"commerce":              c.Header("commerce"),
					"referenceID":           requestBody.ReferenceID,
				},
				Status:       "pending",
				Organization: org,
				Payload:      requestBodyBytes,
			}); err != nil {
				return response.UpsertNodeResponse{}, err
			}
			return response.UpsertNodeResponse{
				Message: "upsert node submitted",
			}, nil
		},
		option.Summary("upsertNode"),
		option.Tags(tagNetwork))
}

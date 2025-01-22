package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		searchNodes,
		httpserver.New,
		usecase.NewSearchNodes)
}
func searchNodes(s httpserver.Server, search usecase.SearchNodes) {
	fuego.Get(s.Manager, "/nodes",
		func(c fuego.ContextWithBody[request.SearchNodesRequest]) ([]response.SearchNodesResponse, error) {
			req, err := c.Body()
			if err != nil {
				return nil, err
			}
			res, err := search(c.Context(), req.Map())
			if err != nil {
				return nil, err
			}
			return response.MapSearchNodesResponse(res), nil
		},
		option.Summary("searchNodes"),
		option.Tags(tagNetwork),
	)
}

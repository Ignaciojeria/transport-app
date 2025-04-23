package fuegoapi

/*
import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/gemini"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase/normalization/chile"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		directionsNormalization,
		httpserver.New,
		gemini.NewGemini2Dot0FlashModelWrapper,
		chile.NewSingleInputPrompt)
}

func directionsNormalization(
	s httpserver.Server,
	model gemini.Gemini2Dot0FlashModelWrapper, retrievePrompt chile.SingleInputPrompt) {

	fuego.Post(s.Manager, "/normalization",
		func(c fuego.ContextWithBody[request.
			AddressNormalizationRequest]) (response.AddressNormalizationResponse, error) {
			req, err := c.Body()
			if err != nil {
				return response.AddressNormalizationResponse{}, err
			}
			userInput, providerInput := req.Map()
			prompt := retrievePrompt(c.Context(), userInput, providerInput)
			res, err := model.StartChat(c.Context(), prompt)
			res.Normalize()
			return response.MapAddressNormalizationResponse(res), err
		},
		option.Tags(tagLocations),
		option.Summary("directionsNormalization"))
}
*/

package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/paulmach/orb"
	"google.golang.org/api/iterator"
)

type Gemini2Dot0FlashModelWrapper struct {
	*genai.GenerativeModel
}

func init() {
	ioc.Registry(NewGemini2Dot0FlashModelWrapper, newClient)
}
func NewGemini2Dot0FlashModelWrapper(client *genai.Client) Gemini2Dot0FlashModelWrapper {
	iter := client.ListModels(context.Background())
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(m.Name, m.Description)
	}
	model := client.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = AddressNormalizationSchema
	return Gemini2Dot0FlashModelWrapper{
		GenerativeModel: model,
	}
}

func (s Gemini2Dot0FlashModelWrapper) StartChat(
	ctx context.Context,
	msg string) (domain.AddressInfo, error) {
	res, err := s.GenerativeModel.StartChat().SendMessage(ctx, genai.Text(msg))
	if err != nil {
		fmt.Printf("Error al enviar mensaje a Gemini: %v\n", err.Error())
		return domain.AddressInfo{}, err
	}
	return getResponse(res)
}

func getResponse(resp *genai.GenerateContentResponse) (domain.AddressInfo, error) {
	var output strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {

			for _, part := range cand.Content.Parts {
				text := fmt.Sprintf("%s", part)
				output.WriteString(text)
			}
		}
	}

	type ResponseData struct {
		ProviderAddress string  `json:"providerAddress"`
		AddressLine1    string  `json:"addressLine1"`
		AddressLine2    string  `json:"addressLine2"`
		District        string  `json:"district"`
		State           string  `json:"state"`
		Province        string  `json:"province"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
	}
	var data ResponseData
	err := json.Unmarshal([]byte(output.String()), &data)
	if err != nil {
		return domain.AddressInfo{}, err
	}

	return domain.AddressInfo{
		AddressLine1: data.AddressLine1,
		//	AddressLine2:    data.AddressLine2,
		//	ProviderAddress: data.ProviderAddress,
		District: domain.District(data.District),
		State:    domain.State(data.State),
		Province: domain.Province(data.Province),
		Coordinates: domain.Coordinates{
			Point: orb.Point{
				data.Longitude,
				data.Latitude,
			},
			Source: "gemini",
			Confidence: domain.CoordinatesConfidence{
				Level:   1.0,
				Message: "High confidence from Gemini",
				Reason:  "Direct geocoding",
			},
		},
	}, nil
}

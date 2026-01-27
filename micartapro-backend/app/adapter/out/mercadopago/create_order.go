package mercadopago

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/httpresty"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type CreateMercadoPagoOrder func(ctx context.Context, request CreateOrderRequest) (CreateOrderResponse, error)

type CreateOrderRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	TotalAmount float64                `json:"total_amount"`
	Items       []OrderItem            `json:"items"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ExternalID  string                 `json:"external_reference,omitempty"`
}

type OrderItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type CreateOrderResponse struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	StatusDetail    string `json:"status_detail"`
	PreferenceID    string `json:"preference_id,omitempty"`
	InitPoint       string `json:"init_point,omitempty"`
	SandboxInitPoint string `json:"sandbox_init_point,omitempty"`
	ExternalReference string `json:"external_reference,omitempty"`
}

func init() {
	ioc.Registry(NewCreateMercadoPagoOrder, httpresty.NewClient, configuration.NewConf)
}

func NewCreateMercadoPagoOrder(cli *resty.Client, conf configuration.Conf) CreateMercadoPagoOrder {
	return func(ctx context.Context, request CreateOrderRequest) (CreateOrderResponse, error) {
		// Detectar si el token es de prueba
		isTestToken := strings.HasPrefix(conf.MERCADOPAGO_ACCESS_TOKEN, "TEST-")
		
		// Usar la API de Orders de Mercado Pago
		url := "https://api.mercadopago.com/checkout/preferences"

		// Construir items según la estructura de Checkout Pro
		items := make([]map[string]interface{}, len(request.Items))
		for i, item := range request.Items {
			itemMap := map[string]interface{}{
				"title":      item.Title,
				"quantity":   item.Quantity,
				"unit_price": item.UnitPrice,
			}
			// Agregar description solo si existe
			if item.Description != "" {
				itemMap["description"] = item.Description
			}
			items[i] = itemMap
		}

		// Construir el request según la estructura de Checkout Pro
		mpRequest := map[string]interface{}{
			"items": items,
			"back_urls": map[string]string{
				"success": conf.MERCADOPAGO_SUCCESS_URL,
				"failure": conf.MERCADOPAGO_FAILURE_URL,
				"pending": conf.MERCADOPAGO_PENDING_URL,
			},
			"auto_return": "approved",
		}

		// Agregar external_reference si existe
		if request.ExternalID != "" {
			mpRequest["external_reference"] = request.ExternalID
		}

		// Agregar notification_url para recibir webhooks automáticamente
		if conf.MERCADOPAGO_WEBHOOK_URL != "" {
			mpRequest["notification_url"] = conf.MERCADOPAGO_WEBHOOK_URL
		}

		// Agregar metadata si existe
		if request.Metadata != nil {
			mpRequest["metadata"] = request.Metadata
		}

		resp, err := cli.R().
			SetContext(ctx).
			SetHeader("Authorization", "Bearer "+conf.MERCADOPAGO_ACCESS_TOKEN).
			SetHeader("Content-Type", "application/json").
			SetBody(mpRequest).
			Post(url)

		if err != nil {
			return CreateOrderResponse{}, fmt.Errorf("failed to make HTTP request: %w", err)
		}

		if resp.IsError() {
			return CreateOrderResponse{}, fmt.Errorf("Mercado Pago API error (status %d): %s", resp.StatusCode(), resp.String())
		}

		// Extraer init_point o sandbox_init_point según el entorno
		var responseData map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
			return CreateOrderResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		var result CreateOrderResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return CreateOrderResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Extraer campos específicos del response
		if initPoint, ok := responseData["init_point"].(string); ok && initPoint != "" {
			result.InitPoint = initPoint
		}
		if sandboxInitPoint, ok := responseData["sandbox_init_point"].(string); ok && sandboxInitPoint != "" {
			result.SandboxInitPoint = sandboxInitPoint
		}
		if id, ok := responseData["id"].(string); ok {
			result.PreferenceID = id
		}

		// Validación: Si es token de prueba y no hay sandbox_init_point, es un problema
		if isTestToken && result.SandboxInitPoint == "" {
			return CreateOrderResponse{}, fmt.Errorf("token de prueba detectado pero Mercado Pago no devolvió sandbox_init_point. Verifica que la aplicación esté en modo Test y que uses un usuario Comprador de prueba para pagar")
		}

		return result, nil
	}
}

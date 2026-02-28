package mercadopago

import (
	"context"
	"fmt"
	"micartapro/app/adapter/out/mercadopago"
	"micartapro/app/events"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"strings"

	ioc "github.com/Ignaciojeria/ioc"
)

type CreateMercadoPagoCheckoutResult struct {
	CheckoutURL  string `json:"checkout_url"`
	PreferenceID string `json:"preference_id"`
}

type CreateMercadoPagoCheckout func(ctx context.Context, checkoutRequest events.MercadoPagoCheckoutRequest, externalReference string) (CreateMercadoPagoCheckoutResult, error)

func init() {
	ioc.Register(NewCreateMercadoPagoCheckout)
}

func NewCreateMercadoPagoCheckout(
	obs observability.Observability,
	createOrder mercadopago.CreateMercadoPagoOrder,
	conf configuration.Conf,
) CreateMercadoPagoCheckout {
	return func(ctx context.Context, checkoutRequest events.MercadoPagoCheckoutRequest, externalReference string) (CreateMercadoPagoCheckoutResult, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "create_mercadopago_checkout")
		defer span.End()

		obs.Logger.InfoContext(spanCtx, "creating_mercadopago_checkout", "externalReference", externalReference)

		// Construir items para Mercado Pago
		items := make([]mercadopago.OrderItem, 0, len(checkoutRequest.Items))
		for _, item := range checkoutRequest.Items {
			items = append(items, mercadopago.OrderItem{
				Title:       item.ProductName,
				Description: fmt.Sprintf("Cantidad: %.2f %s", item.Quantity, item.Unit),
				Quantity:    int(item.Quantity),
				UnitPrice:   float64(item.UnitPrice),
			})
		}

		// Agregar delivery fee como item adicional si existe
		if checkoutRequest.Totals.DeliveryFee > 0 {
			items = append(items, mercadopago.OrderItem{
				Title:       "Costo de envío",
				Description: "Costo de entrega a domicilio",
				Quantity:    1,
				UnitPrice:   float64(checkoutRequest.Totals.DeliveryFee),
			})
		}

		// Extraer user_id del contexto si está disponible
		userID, _ := sharedcontext.UserIDFromContext(spanCtx)

		// Construir metadata: combinar metadata del request con metadata por defecto
		metadata := make(map[string]interface{})

		// Agregar metadata por defecto
		metadata["business_name"] = checkoutRequest.BusinessInfo.BusinessName
		metadata["whatsapp"] = checkoutRequest.BusinessInfo.Whatsapp
		metadata["fulfillment_type"] = checkoutRequest.Fulfillment.Type

		// Agregar user_id al metadata si está disponible
		if userID != "" {
			metadata["user_id"] = userID
		}

		// Agregar metadata personalizado del request (sobrescribe los valores por defecto si hay conflictos)
		if checkoutRequest.Metadata != nil {
			for k, v := range checkoutRequest.Metadata {
				metadata[k] = v
			}
		}

		// Crear la orden en Mercado Pago
		mpRequest := mercadopago.CreateOrderRequest{
			Title:       fmt.Sprintf("Pedido - %s", checkoutRequest.BusinessInfo.BusinessName),
			Description: fmt.Sprintf("Pedido de %s", checkoutRequest.BusinessInfo.BusinessName),
			TotalAmount: float64(checkoutRequest.Totals.Total),
			Items:       items,
			ExternalID:  externalReference,
			Metadata:    metadata,
		}

		result, err := createOrder(spanCtx, mpRequest)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_creating_mercadopago_order", "error", err)
			return CreateMercadoPagoCheckoutResult{}, fmt.Errorf("failed to create Mercado Pago order: %w", err)
		}

		// Determinar qué URL usar (init_point o sandbox_init_point)
		// Si el token es de prueba, FORZAR el uso de sandbox_init_point
		isTestToken := strings.HasPrefix(conf.MERCADOPAGO_ACCESS_TOKEN, "TEST-")

		var checkoutURL string
		var isSandbox bool

		if isTestToken {
			// Con token de prueba, SOLO usar sandbox_init_point
			checkoutURL = result.SandboxInitPoint
			isSandbox = true
			if checkoutURL == "" {
				obs.Logger.ErrorContext(spanCtx, "test_token_but_no_sandbox_url",
					"preferenceID", result.PreferenceID,
					"initPoint", result.InitPoint,
					"sandboxInitPoint", result.SandboxInitPoint)
				return CreateMercadoPagoCheckoutResult{}, fmt.Errorf("token de prueba detectado pero no se recibió sandbox_init_point. Asegúrate de usar un usuario Comprador de prueba creado desde el panel de desarrolladores")
			}
		} else {
			// Con token de producción, usar init_point
			checkoutURL = result.InitPoint
			isSandbox = false
			if checkoutURL == "" {
				checkoutURL = result.SandboxInitPoint // Fallback
				isSandbox = checkoutURL != ""
			}
		}

		if checkoutURL == "" {
			obs.Logger.ErrorContext(spanCtx, "no_checkout_url_returned",
				"preferenceID", result.PreferenceID,
				"initPoint", result.InitPoint,
				"sandboxInitPoint", result.SandboxInitPoint,
				"isTestToken", isTestToken)
			return CreateMercadoPagoCheckoutResult{}, fmt.Errorf("no checkout URL returned from Mercado Pago")
		}

		obs.Logger.InfoContext(spanCtx, "mercadopago_checkout_created_successfully",
			"preferenceID", result.PreferenceID,
			"checkoutURL", checkoutURL,
			"isSandbox", isSandbox,
			"initPoint", result.InitPoint,
			"sandboxInitPoint", result.SandboxInitPoint)

		return CreateMercadoPagoCheckoutResult{
			CheckoutURL:  checkoutURL,
			PreferenceID: result.PreferenceID,
		}, nil
	}
}

package mercadopago

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"micartapro/app/adapter/out/mercadopago"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type ProcessMercadoPagoWebhook func(ctx context.Context, webhookData map[string]interface{}) error

type MercadoPagoWebhookData struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Action   string                 `json:"action"`
	Data     map[string]interface{} `json:"data"`
	Date     string                 `json:"date"`
	UserID   interface{}            `json:"user_id,omitempty"` // Puede venir como número o string
	LiveMode bool                   `json:"live_mode"`
}

func init() {
	ioc.Registry(NewProcessMercadoPagoWebhook,
		observability.NewObservability,
		mercadopago.NewGetMercadoPagoPayment,
		supabaserepo.NewSaveBillingEvent,
		supabaserepo.NewGrantCredits,
	)
}

func NewProcessMercadoPagoWebhook(
	obs observability.Observability,
	getPayment mercadopago.GetMercadoPagoPayment,
	saveBillingEvent supabaserepo.SaveBillingEvent,
	grantCredits supabaserepo.GrantCredits,
) ProcessMercadoPagoWebhook {
	return func(ctx context.Context, webhookData map[string]interface{}) error {
		spanCtx, span := obs.Tracer.Start(ctx, "process_mercadopago_webhook")
		defer span.End()

		obs.Logger.InfoContext(spanCtx, "processing_mercadopago_webhook", "webhookData", webhookData)

		// Extraer el tipo de webhook directamente del map
		webhookType, ok := webhookData["type"].(string)
		if !ok {
			obs.Logger.ErrorContext(spanCtx, "webhook_type_not_found", "webhookData", webhookData)
			return fmt.Errorf("webhook type not found or invalid")
		}

		// Mercado Pago envía diferentes tipos de webhooks
		// El más común es "payment" que contiene el ID del pago
		if webhookType == "payment" {
			// Obtener el objeto data
			dataObj, ok := webhookData["data"].(map[string]interface{})
			if !ok {
				obs.Logger.ErrorContext(spanCtx, "webhook_data_not_found", "webhookData", webhookData)
				return fmt.Errorf("webhook data not found or invalid")
			}

			// Obtener el ID del pago desde data.id
			var paymentID string
			if idStr, ok := dataObj["id"].(string); ok {
				paymentID = idStr
			} else if idNum, ok := dataObj["id"].(float64); ok {
				paymentID = fmt.Sprintf("%.0f", idNum)
			} else if idInt, ok := dataObj["id"].(int64); ok {
				paymentID = fmt.Sprintf("%d", idInt)
			} else {
				obs.Logger.ErrorContext(spanCtx, "payment_id_not_found", "data", dataObj)
				return fmt.Errorf("payment ID not found in webhook data")
			}

			// Obtener información completa del pago desde la API
			// Intentar varias veces con retry porque el webhook puede llegar antes de que el pago esté disponible
			var payment mercadopago.PaymentResponse
			var err error
			maxRetries := 3
			retryDelay := time.Second * 2
			
			for attempt := 0; attempt < maxRetries; attempt++ {
				if attempt > 0 {
					obs.Logger.InfoContext(spanCtx, "retrying_get_payment", 
						"paymentID", paymentID,
						"attempt", attempt+1,
						"maxRetries", maxRetries)
					time.Sleep(retryDelay)
					retryDelay *= 2 // Exponential backoff
				}
				
				payment, err = getPayment(spanCtx, paymentID)
				if err == nil {
					// Pago encontrado exitosamente
					break
				}
				
				// Si es un error 404, puede ser que el pago aún no esté disponible
				// Continuar con el retry
				if attempt < maxRetries-1 {
					obs.Logger.WarnContext(spanCtx, "payment_not_available_yet", 
						"paymentID", paymentID,
						"attempt", attempt+1,
						"error", err,
						"will_retry", true)
				}
			}
			
			var eventPayload json.RawMessage
			var eventType string
			var providerEventID string
			var subscriptionID string
			var providerCreatedAt time.Time
			var paymentMetadata map[string]interface{}

			if err != nil {
				// Si el pago no se encuentra (404), puede ser un webhook de prueba o un pago aún no disponible
				// Guardamos el webhook completo para tener registro
				obs.Logger.WarnContext(spanCtx, "payment_not_found_or_error", 
					"paymentID", paymentID, 
					"error", err,
					"note", "saving webhook data instead of payment details")
				
				// Usar el webhook completo como payload
				webhookPayload, marshalErr := json.Marshal(webhookData)
				if marshalErr != nil {
					obs.Logger.ErrorContext(spanCtx, "error_marshaling_webhook", "error", marshalErr)
					return fmt.Errorf("failed to marshal webhook data: %w", marshalErr)
				}
				eventPayload = json.RawMessage(webhookPayload)
				
				// Extraer action si existe, sino usar type
				if action, ok := webhookData["action"].(string); ok && action != "" {
					eventType = action
				} else {
					eventType = "payment.unknown"
				}
				
				// Usar el ID del webhook como provider event ID
				if webhookID, ok := webhookData["id"].(string); ok {
					providerEventID = webhookID
				} else if webhookIDNum, ok := webhookData["id"].(float64); ok {
					providerEventID = fmt.Sprintf("%.0f", webhookIDNum)
				} else {
					providerEventID = paymentID // Fallback al payment ID
				}
				
				// Intentar extraer external_reference del metadata si existe
				if metadata, ok := webhookData["metadata"].(map[string]interface{}); ok {
					if extRef, ok := metadata["external_reference"].(string); ok {
						subscriptionID = extRef
					}
				}
				
				// Parsear fecha del webhook
				if dateStr, ok := webhookData["date_created"].(string); ok {
					if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
						providerCreatedAt = parsed
					} else {
						providerCreatedAt = time.Now()
					}
				} else {
					providerCreatedAt = time.Now()
				}
			} else {
				// Pago encontrado exitosamente
				obs.Logger.InfoContext(spanCtx, "payment_received", 
					"paymentID", paymentID,
					"status", payment.Status,
					"externalReference", payment.ExternalReference,
					"transactionAmount", payment.TransactionAmount,
					"hasTransactionDetails", payment.TransactionDetails != nil)
				
				// Log detallado de transaction_details si está disponible
				if payment.TransactionDetails != nil {
					obs.Logger.InfoContext(spanCtx, "payment_transaction_details",
						"totalPaidAmount", payment.TransactionDetails.TotalPaidAmount,
						"netReceivedAmount", payment.TransactionDetails.NetReceivedAmount)
				} else {
					obs.Logger.WarnContext(spanCtx, "transaction_details_not_parsed",
						"paymentID", paymentID,
						"note", "TransactionDetails is nil, may need to check JSON parsing")
				}
				
				// Log detallado de transaction_details si está disponible
				if payment.TransactionDetails != nil {
					obs.Logger.InfoContext(spanCtx, "payment_transaction_details",
						"totalPaidAmount", payment.TransactionDetails.TotalPaidAmount,
						"netReceivedAmount", payment.TransactionDetails.NetReceivedAmount)
				}

				// Guardar el evento de billing con los datos del pago
				paymentPayload, err := json.Marshal(payment)
				if err != nil {
					obs.Logger.ErrorContext(spanCtx, "error_marshaling_payment", "error", err)
					return fmt.Errorf("failed to marshal payment data: %w", err)
				}
				eventPayload = json.RawMessage(paymentPayload)

				// Determinar el tipo de evento según el status
				eventType = "payment." + payment.Status
				if payment.Status == "approved" {
					eventType = "payment.approved"
				} else if payment.Status == "pending" {
					eventType = "payment.pending"
				} else if payment.Status == "rejected" || payment.Status == "cancelled" {
					eventType = "payment.failed"
				}

				// Parsear la fecha
				providerCreatedAt, err = time.Parse(time.RFC3339, payment.DateLastUpdated)
				if err != nil {
					// Intentar otros formatos comunes
					providerCreatedAt, err = time.Parse("2006-01-02T15:04:05.000-07:00", payment.DateLastUpdated)
					if err != nil {
						providerCreatedAt = time.Now()
					}
				}

				providerEventID = fmt.Sprintf("%d", payment.ID)
				subscriptionID = payment.ExternalReference
				paymentMetadata = payment.Metadata
			}

			// Extraer user_id del metadata o del external_reference
			// Prioridad: 1) metadata.user_id, 2) parsear external_reference
			var userID *uuid.UUID
			
			// Intentar extraer del metadata del pago primero (más confiable)
			if paymentMetadata != nil {
				if userIDStr, ok := paymentMetadata["user_id"].(string); ok && userIDStr != "" {
					if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
						userID = &parsedUserID
						obs.Logger.InfoContext(spanCtx, "extracted_user_id_from_payment_metadata", 
							"userID", parsedUserID.String())
					}
				}
			}
			
			// Si no se encontró en metadata, intentar parsear del external_reference
			// El formato es: {user_id}_{uuid}
			if userID == nil && subscriptionID != "" {
				parts := strings.Split(subscriptionID, "_")
				if len(parts) > 0 {
					// La primera parte es el user_id
					if parsedUserID, err := uuid.Parse(parts[0]); err == nil {
						userID = &parsedUserID
						obs.Logger.InfoContext(spanCtx, "extracted_user_id_from_external_reference", 
							"userID", parsedUserID.String(),
							"externalReference", subscriptionID)
					} else {
						obs.Logger.WarnContext(spanCtx, "failed_to_parse_user_id_from_external_reference", 
							"externalReference", subscriptionID,
							"error", err)
					}
				}
			}

			// Crear el BillingEvent
			billingEvent := billing.BillingEvent{
				Provider:          "mercadopago",
				ProviderEventID:   providerEventID,
				EventType:         eventType,
				SubscriptionID:    subscriptionID,
				Payload:           eventPayload,
				ProviderCreatedAt: providerCreatedAt,
				UserID:            userID,
			}

			// Guardar el evento
			err = saveBillingEvent(spanCtx, billingEvent)

			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_saving_billing_event", "error", err)
				return fmt.Errorf("failed to save billing event: %w", err)
			}

			// Si el pago fue aprobado y tenemos user_id, otorgar créditos
			if eventType == "payment.approved" && userID != nil {
				// Calcular créditos basado en el monto del pago
				// 1 crédito = $140 CLP
				var creditsAmount int
				// Verificar si el pago se obtuvo exitosamente (payment.ID solo se asigna cuando el pago se obtiene)
				if payment.ID > 0 && payment.TransactionAmount > 0 {
					// Usar total_paid_amount si está disponible (monto total pagado por el usuario)
					// Si no, usar transaction_amount (monto de la transacción)
					var amountToUse float64
					if payment.TransactionDetails != nil && payment.TransactionDetails.TotalPaidAmount > 0 {
						amountToUse = payment.TransactionDetails.TotalPaidAmount
						obs.Logger.InfoContext(spanCtx, "using_total_paid_amount_for_credits",
							"totalPaidAmount", payment.TransactionDetails.TotalPaidAmount,
							"transactionAmount", payment.TransactionAmount)
					} else {
						amountToUse = payment.TransactionAmount
					}
					
					// Log del monto recibido para debugging
					obs.Logger.InfoContext(spanCtx, "calculating_credits_from_payment",
						"paymentID", paymentID,
						"transactionAmount", payment.TransactionAmount,
						"amountToUse", amountToUse,
						"currencyID", payment.CurrencyID)
					
					// 1 crédito por cada $140 CLP pagados
					creditsAmount = int(amountToUse / 140)
					
					// Log del cálculo
					obs.Logger.InfoContext(spanCtx, "credits_calculation",
						"amountToUse", amountToUse,
						"creditsAmount", creditsAmount,
						"formula", fmt.Sprintf("%.2f / 140 = %d", amountToUse, creditsAmount))
					
					// Mínimo 1 crédito si el pago es menor a $140
					if creditsAmount < 1 {
						creditsAmount = 1
					}
				} else {
					// Si no pudimos obtener el pago o el monto no está disponible, otorgar un crédito por defecto
					obs.Logger.WarnContext(spanCtx, "payment_data_not_available_for_credits",
						"paymentID", paymentID,
						"paymentIDValue", payment.ID,
						"transactionAmount", payment.TransactionAmount)
					creditsAmount = 1
				}

				paymentIDStr := providerEventID
				var amountForDescription float64
				if payment.TransactionDetails != nil && payment.TransactionDetails.TotalPaidAmount > 0 {
					amountForDescription = payment.TransactionDetails.TotalPaidAmount
				} else {
					amountForDescription = payment.TransactionAmount
				}
				description := fmt.Sprintf("Créditos otorgados por pago aprobado de MercadoPago (ID: %s, Monto: %.2f CLP)", paymentIDStr, amountForDescription)

				_, grantErr := grantCredits(spanCtx, billing.GrantCreditsRequest{
					UserID:      *userID,
					Amount:      creditsAmount,
					Source:      "payment.mercadopago",
					SourceID:    &paymentIDStr,
					Description: &description,
				})

				if grantErr != nil {
					obs.Logger.ErrorContext(spanCtx, "error_granting_credits", 
						"error", grantErr,
						"userID", userID.String(),
						"amount", creditsAmount,
						"paymentID", paymentIDStr)
					// No retornamos error aquí para no bloquear el webhook
					// El pago ya fue registrado, los créditos se pueden otorgar manualmente si es necesario
				} else {
					obs.Logger.InfoContext(spanCtx, "credits_granted_successfully",
						"userID", userID.String(),
						"amount", creditsAmount,
						"paymentID", paymentIDStr)
				}
			}

			obs.Logger.InfoContext(spanCtx, "mercadopago_webhook_processed_successfully", 
				"paymentID", paymentID,
				"eventType", eventType)
		} else {
			obs.Logger.InfoContext(spanCtx, "webhook_type_not_handled", "type", webhookType)
		}

		return nil
	}
}

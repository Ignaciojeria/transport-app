package mapper

import (
	"context"
	"encoding/json"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("Webhook Table Mapping", func() {
	var (
		ctx           context.Context
		tenantID      uuid.UUID
		sampleWebhook domain.Webhook
		sampleTime    time.Time
	)

	BeforeEach(func() {
		// Setup context with tenant ID
		tenantID = uuid.New()
		tID, err := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID.String())
		Expect(err).ToNot(HaveOccurred())

		bag, err := baggage.New(tID)
		Expect(err).ToNot(HaveOccurred())

		ctx = baggage.ContextWithBaggage(context.Background(), bag)

		// Setup sample webhook with all field types
		sampleTime = time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)
		sampleWebhook = domain.Webhook{
			Type: "order.created",
			URL:  "https://example.com/webhook",
			Headers: map[string]string{
				"Authorization": "Bearer token123",
				"Content-Type":  "application/json",
			},
			RetryPolicy: domain.RetryPolicy{
				MaxRetries:     3,
				BackoffSeconds: 60,
			},
			CreatedAt: sampleTime,
			UpdatedAt: sampleTime.Add(time.Hour),
		}
	})

	Describe("MapWebhookToTable", func() {
		Context("when mapping a valid webhook", func() {
			It("should successfully map all fields without data loss", func() {
				// Act
				tableWebhook, err := MapWebhookToTable(ctx, sampleWebhook)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(tableWebhook.TenantID).To(Equal(tenantID))
				Expect(tableWebhook.DocumentID).ToNot(BeEmpty())

				// Verify that complex types are preserved correctly
				payload := map[string]string(tableWebhook.Payload)

				// Check basic string fields
				Expect(payload["type"]).To(Equal(`"order.created"`))
				Expect(payload["url"]).To(Equal(`"https://example.com/webhook"`))

				// Check that headers map is preserved
				var headers map[string]string
				err = json.Unmarshal([]byte(payload["headers"]), &headers)
				Expect(err).ToNot(HaveOccurred())
				Expect(headers["Authorization"]).To(Equal("Bearer token123"))
				Expect(headers["Content-Type"]).To(Equal("application/json"))

				// Check that retry policy nested struct is preserved
				var retryPolicy domain.RetryPolicy
				err = json.Unmarshal([]byte(payload["retryPolicy"]), &retryPolicy)
				Expect(err).ToNot(HaveOccurred())
				Expect(retryPolicy.MaxRetries).To(Equal(3))
				Expect(retryPolicy.BackoffSeconds).To(Equal(60))

				// Check that time fields are preserved
				var createdAt time.Time
				err = json.Unmarshal([]byte(payload["createdAt"]), &createdAt)
				Expect(err).ToNot(HaveOccurred())
				Expect(createdAt.Equal(sampleTime)).To(BeTrue())

				var updatedAt time.Time
				err = json.Unmarshal([]byte(payload["updatedAt"]), &updatedAt)
				Expect(err).ToNot(HaveOccurred())
				Expect(updatedAt.Equal(sampleTime.Add(time.Hour))).To(BeTrue())
			})
		})

		Context("when mapping a webhook with nil/empty values", func() {
			It("should handle empty headers map", func() {
				webhook := sampleWebhook
				webhook.Headers = make(map[string]string)

				tableWebhook, err := MapWebhookToTable(ctx, webhook)

				Expect(err).ToNot(HaveOccurred())
				payload := map[string]string(tableWebhook.Payload)

				var headers map[string]string
				err = json.Unmarshal([]byte(payload["headers"]), &headers)
				Expect(err).ToNot(HaveOccurred())
				Expect(headers).To(BeEmpty())
			})

			It("should handle zero values in retry policy", func() {
				webhook := sampleWebhook
				webhook.RetryPolicy = domain.RetryPolicy{
					MaxRetries:     0,
					BackoffSeconds: 0,
				}

				tableWebhook, err := MapWebhookToTable(ctx, webhook)

				Expect(err).ToNot(HaveOccurred())
				payload := map[string]string(tableWebhook.Payload)

				var retryPolicy domain.RetryPolicy
				err = json.Unmarshal([]byte(payload["retryPolicy"]), &retryPolicy)
				Expect(err).ToNot(HaveOccurred())
				Expect(retryPolicy.MaxRetries).To(Equal(0))
				Expect(retryPolicy.BackoffSeconds).To(Equal(0))
			})
		})

		Context("when mapping a webhook with extreme values", func() {
			It("should handle large retry values", func() {
				webhook := sampleWebhook
				webhook.RetryPolicy = domain.RetryPolicy{
					MaxRetries:     999999,
					BackoffSeconds: 86400, // 1 day
				}

				tableWebhook, err := MapWebhookToTable(ctx, webhook)

				Expect(err).ToNot(HaveOccurred())
				payload := map[string]string(tableWebhook.Payload)

				var retryPolicy domain.RetryPolicy
				err = json.Unmarshal([]byte(payload["retryPolicy"]), &retryPolicy)
				Expect(err).ToNot(HaveOccurred())
				Expect(retryPolicy.MaxRetries).To(Equal(999999))
				Expect(retryPolicy.BackoffSeconds).To(Equal(86400))
			})
		})
	})

	Describe("table.Webhook.Map", func() {
		Context("when mapping a table webhook back to domain", func() {
			It("should successfully reconstruct the original webhook", func() {
				// Arrange - first create a table webhook
				tableWebhook, err := MapWebhookToTable(ctx, sampleWebhook)
				Expect(err).ToNot(HaveOccurred())

				// Act - map back to domain
				reconstructedWebhook, err := tableWebhook.Map()

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(reconstructedWebhook.Type).To(Equal(sampleWebhook.Type))
				Expect(reconstructedWebhook.URL).To(Equal(sampleWebhook.URL))
				Expect(reconstructedWebhook.Headers).To(Equal(sampleWebhook.Headers))
				Expect(reconstructedWebhook.RetryPolicy).To(Equal(sampleWebhook.RetryPolicy))

				// Time comparison needs to be done carefully due to precision
				Expect(reconstructedWebhook.CreatedAt.Equal(sampleWebhook.CreatedAt)).To(BeTrue())
				Expect(reconstructedWebhook.UpdatedAt.Equal(sampleWebhook.UpdatedAt)).To(BeTrue())
			})
		})

		Context("when mapping malformed data", func() {
			It("should return error for invalid JSON in payload", func() {
				// Arrange
				tableWebhook := table.Webhook{
					Payload: table.JSONMap{
						"type": "invalid-json{",
					},
				}

				// Act
				_, err := tableWebhook.Map()

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to unmarshal field type"))
			})

			It("should return error for incompatible types", func() {
				// Arrange
				tableWebhook := table.Webhook{
					Payload: table.JSONMap{
						"type":        `"order.created"`,
						"retryPolicy": `"not-a-struct"`, // Should be an object, not a string
					},
				}

				// Act
				_, err := tableWebhook.Map()

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to unmarshal to domain.Webhook"))
			})
		})
	})

	Describe("Round-trip mapping", func() {
		It("should preserve all data through complete round-trip", func() {
			// Act - round trip: domain -> table -> domain
			tableWebhook, err := MapWebhookToTable(ctx, sampleWebhook)
			Expect(err).ToNot(HaveOccurred())

			reconstructed, err := tableWebhook.Map()
			Expect(err).ToNot(HaveOccurred())

			// Assert - compare all fields
			Expect(reconstructed.Type).To(Equal(sampleWebhook.Type))
			Expect(reconstructed.URL).To(Equal(sampleWebhook.URL))
			Expect(reconstructed.Headers).To(Equal(sampleWebhook.Headers))
			Expect(reconstructed.RetryPolicy.MaxRetries).To(Equal(sampleWebhook.RetryPolicy.MaxRetries))
			Expect(reconstructed.RetryPolicy.BackoffSeconds).To(Equal(sampleWebhook.RetryPolicy.BackoffSeconds))
			Expect(reconstructed.CreatedAt.Equal(sampleWebhook.CreatedAt)).To(BeTrue())
			Expect(reconstructed.UpdatedAt.Equal(sampleWebhook.UpdatedAt)).To(BeTrue())
		})

		It("should work with multiple webhooks", func() {
			webhooks := []domain.Webhook{
				sampleWebhook,
				{
					Type:        "order.updated",
					URL:         "https://api.example.com/hooks",
					Headers:     map[string]string{"X-API-Key": "secret"},
					RetryPolicy: domain.RetryPolicy{MaxRetries: 5, BackoffSeconds: 120},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Type:        "order.cancelled",
					URL:         "https://webhook.test.com",
					Headers:     map[string]string{},
					RetryPolicy: domain.RetryPolicy{MaxRetries: 1, BackoffSeconds: 30},
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			}

			for i, webhook := range webhooks {
				// Map to table and back
				tableWebhook, err := MapWebhookToTable(ctx, webhook)
				Expect(err).ToNot(HaveOccurred(), "Failed for webhook %d", i)

				reconstructed, err := tableWebhook.Map()
				Expect(err).ToNot(HaveOccurred(), "Failed reconstruction for webhook %d", i)

				// Verify preservation
				Expect(reconstructed.Type).To(Equal(webhook.Type), "Type mismatch for webhook %d", i)
				Expect(reconstructed.URL).To(Equal(webhook.URL), "URL mismatch for webhook %d", i)
				Expect(reconstructed.Headers).To(Equal(webhook.Headers), "Headers mismatch for webhook %d", i)
				Expect(reconstructed.RetryPolicy).To(Equal(webhook.RetryPolicy), "RetryPolicy mismatch for webhook %d", i)
			}
		})
	})
})

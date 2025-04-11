package tidbrepository

import (
	"context"
	"strconv"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertOrderType", func() {
	var ctx context.Context

	// Helper function to create context with organization
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		// Create context with organization1
		ctx = createOrgContext(organization1)

		// Limpia la tabla antes de cada test
		err := connection.DB.Exec("DELETE FROM order_types").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert order type if it does not exist", func() {
		ot := domain.OrderType{
			Type:        "retail",
			Description: "Entrega a cliente final",
		}

		upsert := NewUpsertOrderType(connection)
		err := upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		var result table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID(ctx)).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.Type).To(Equal("retail"))
		Expect(result.Description).To(Equal("Entrega a cliente final"))
	})

	It("should not create if document is present", func() {
		ot := domain.OrderType{
			Type:        "express",
			Description: "Entrega en 1 hora",
		}

		upsert := NewUpsertOrderType(connection)
		err := upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		// Get initial record to compare timestamps later
		var initialRecord table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID(ctx)).
			First(&initialRecord).Error
		Expect(err).ToNot(HaveOccurred())
		initialTimestamp := initialRecord.CreatedAt

		// Ejecutamos con los mismos datos
		err = upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1))) // sigue habiendo solo uno

		// Verify timestamp hasn't changed
		var finalRecord table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID(ctx)).
			First(&finalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(finalRecord.CreatedAt).To(Equal(initialTimestamp))
	})

	It("should allow same type for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		ot1 := domain.OrderType{
			Type:        "multi-org",
			Description: "Testing multiple organizations",
		}

		ot2 := domain.OrderType{
			Type:        "multi-org",
			Description: "Testing multiple organizations",
		}

		upsert := NewUpsertOrderType(connection)

		err := upsert(ctx1, ot1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, ot2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("order_types").
			Where("type = ?", "multi-org").
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2))) // One per organization

		// Verify they have different document IDs
		Expect(ot1.DocID(ctx1)).ToNot(Equal(ot2.DocID(ctx2)))
	})

	It("should fail if table does not exist", func() {
		ot := domain.OrderType{
			Type:        "ghost",
			Description: "Fallará",
		}

		upsert := NewUpsertOrderType(noTablesContainerConnection)
		err := upsert(ctx, ot)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_types"))
	})
})

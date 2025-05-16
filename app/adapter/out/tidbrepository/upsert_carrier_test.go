package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertCarrier", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		// Clean the table before each test
		err := connection.DB.Exec("DELETE FROM carriers").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert carrier if not exists", func() {
		ctx := createOrgContext(organization1)

		carrier := domain.Carrier{
			Name:       "Transportes XYZ",
			NationalID: "12345678-9",
		}

		upsert := NewUpsertCarrier(connection)
		err := upsert(ctx, carrier)
		Expect(err).ToNot(HaveOccurred())

		var dbCarrier table.Carrier
		err = connection.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", carrier.DocID(ctx)).
			First(&dbCarrier).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier.Name).To(Equal("Transportes XYZ"))
		Expect(dbCarrier.NationalID).To(Equal("12345678-9"))
	})

	It("should update carrier if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.Carrier{
			Name:       "Original Name",
			NationalID: "11111111-1",
		}

		upsert := NewUpsertCarrier(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Carrier{
			Name:       "Updated Name",
			NationalID: "11111111-1",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbCarrier table.Carrier
		err = connection.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbCarrier).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier.Name).To(Equal("Updated Name"))
		Expect(dbCarrier.NationalID).To(Equal("11111111-1"))
	})

	It("should not update if no fields changed", func() {
		ctx := createOrgContext(organization1)

		carrier := domain.Carrier{
			Name:       "No Change",
			NationalID: "22222222-2",
		}

		upsert := NewUpsertCarrier(connection)
		err := upsert(ctx, carrier)
		Expect(err).ToNot(HaveOccurred())

		// Capture original record to verify timestamp doesn't change
		var originalRecord table.Carrier
		err = connection.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", carrier.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Execute again without changes
		err = upsert(ctx, carrier)
		Expect(err).ToNot(HaveOccurred())

		var dbCarrier table.Carrier
		err = connection.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", carrier.DocID(ctx)).
			First(&dbCarrier).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier.Name).To(Equal("No Change"))
		Expect(dbCarrier.NationalID).To(Equal("22222222-2"))
		Expect(dbCarrier.UpdatedAt).To(Equal(originalRecord.UpdatedAt)) // Verify timestamp didn't change
	})

	It("should allow same carrier for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		carrier1 := domain.Carrier{
			Name:       "Multi Org Carrier",
			NationalID: "33333333-3",
		}

		carrier2 := domain.Carrier{
			Name:       "Multi Org Carrier",
			NationalID: "33333333-3",
		}

		upsert := NewUpsertCarrier(connection)

		err := upsert(ctx1, carrier1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, carrier2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("carriers").
			Where("national_id = ?", carrier1.NationalID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(carrier1.DocID(ctx1)).ToNot(Equal(carrier2.DocID(ctx2)))
	})

	It("should fail if database has no carriers table", func() {
		ctx := createOrgContext(organization1)

		carrier := domain.Carrier{
			Name:       "Error Case",
			NationalID: "44444444-4",
		}

		upsert := NewUpsertCarrier(noTablesContainerConnection)
		err := upsert(ctx, carrier)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("carriers"))
	})
})

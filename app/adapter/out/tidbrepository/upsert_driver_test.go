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

var _ = Describe("UpsertDriver", func() {
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
		err := connection.DB.Exec("DELETE FROM drivers").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert driver if not exists", func() {
		ctx := createOrgContext(organization1)

		driver := domain.Driver{
			Name:       "Juan Pérez",
			NationalID: "12345678-9",
			Email:      "juan@example.com",
		}

		upsert := NewUpsertDriver(connection)
		err := upsert(ctx, driver)
		Expect(err).ToNot(HaveOccurred())

		var dbDriver table.Driver
		err = connection.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", driver.DocID(ctx)).
			First(&dbDriver).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver.Name).To(Equal("Juan Pérez"))
		Expect(dbDriver.NationalID).To(Equal("12345678-9"))
		Expect(dbDriver.Email).To(Equal("juan@example.com"))
	})

	It("should update driver if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.Driver{
			Name:       "Original Name",
			NationalID: "11111111-1",
			Email:      "original@example.com",
		}

		upsert := NewUpsertDriver(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Capture the original document ID
		originalDocID := original.DocID(ctx)

		modified := domain.Driver{
			Name:       "Updated Name",
			NationalID: "11111111-1", // Same NationalID to ensure same DocumentID
			Email:      "updated@example.com",
		}

		// Verify the DocumentID is the same
		Expect(modified.DocID(ctx)).To(Equal(originalDocID))

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbDriver table.Driver
		err = connection.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", originalDocID).
			First(&dbDriver).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver.Name).To(Equal("Updated Name"))
		Expect(dbDriver.NationalID).To(Equal("11111111-1"))
		Expect(dbDriver.Email).To(Equal("updated@example.com"))
	})

	It("should not update if no fields changed", func() {
		ctx := createOrgContext(organization1)

		driver := domain.Driver{
			Name:       "No Change",
			NationalID: "22222222-2",
			Email:      "nochange@example.com",
		}

		upsert := NewUpsertDriver(connection)
		err := upsert(ctx, driver)
		Expect(err).ToNot(HaveOccurred())

		// Capture original record to verify timestamp doesn't change
		var originalRecord table.Driver
		err = connection.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", driver.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Execute again without changes
		err = upsert(ctx, driver)
		Expect(err).ToNot(HaveOccurred())

		var dbDriver table.Driver
		err = connection.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", driver.DocID(ctx)).
			First(&dbDriver).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver.Name).To(Equal("No Change"))
		Expect(dbDriver.NationalID).To(Equal("22222222-2"))
		Expect(dbDriver.Email).To(Equal("nochange@example.com"))
		Expect(dbDriver.UpdatedAt).To(Equal(originalRecord.UpdatedAt)) // Verify timestamp didn't change
	})

	It("should allow same driver for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		driver1 := domain.Driver{
			Name:       "Multi Org Driver",
			NationalID: "33333333-3",
			Email:      "multi@example.com",
		}

		driver2 := domain.Driver{
			Name:       "Multi Org Driver",
			NationalID: "33333333-3",
			Email:      "multi@example.com",
		}

		upsert := NewUpsertDriver(connection)

		err := upsert(ctx1, driver1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, driver2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("drivers").
			Where("national_id = ?", driver1.NationalID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(driver1.DocID(ctx1)).ToNot(Equal(driver2.DocID(ctx2)))
	})

	It("should fail if database has no drivers table", func() {
		ctx := createOrgContext(organization1)

		driver := domain.Driver{
			Name:       "Error Case",
			NationalID: "44444444-4",
			Email:      "error@example.com",
		}

		upsert := NewUpsertDriver(noTablesContainerConnection)
		err := upsert(ctx, driver)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("drivers"))
	})
})

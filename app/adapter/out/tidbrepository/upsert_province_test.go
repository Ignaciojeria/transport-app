package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertProvince", func() {
	var (
		ctx          context.Context
		conn         database.ConnectionFactory
		upsert       UpsertProvince
		testProvince domain.Province
	)

	// Helper function to create context with organization
	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		ctx = createOrgContext(organization1)
		conn = connection
		upsert = NewUpsertProvince(conn)
		testProvince = domain.Province("Test Province")
	})

	AfterEach(func() {
		conn.WithContext(ctx).Exec("DELETE FROM provinces")
	})

	Describe("UpsertProvince", func() {
		It("should insert a new province", func() {
			err := upsert(ctx, testProvince)
			Expect(err).To(BeNil())

			var savedProvince table.Province
			err = conn.WithContext(ctx).
				Table("provinces").
				Where("document_id = ?", testProvince.DocID(ctx).String()).
				First(&savedProvince).Error
			Expect(err).To(BeNil())
			Expect(savedProvince.Name).To(Equal(testProvince.String()))
		})

		It("should create a new record when province name changes", func() {
			// First insert
			err := upsert(ctx, testProvince)
			Expect(err).To(BeNil())

			var firstProvince table.Province
			err = conn.WithContext(ctx).
				Table("provinces").
				Where("document_id = ?", testProvince.DocID(ctx).String()).
				First(&firstProvince).Error
			Expect(err).To(BeNil())
			firstID := firstProvince.ID

			// Insert with different name
			updatedProvince := domain.Province("Updated Province")
			err = upsert(ctx, updatedProvince)
			Expect(err).To(BeNil())

			var secondProvince table.Province
			err = conn.WithContext(ctx).
				Table("provinces").
				Where("document_id = ?", updatedProvince.DocID(ctx).String()).
				First(&secondProvince).Error
			Expect(err).To(BeNil())
			Expect(secondProvince.ID).ToNot(Equal(firstID)) // Should be a new record
		})

		It("should handle multiple provinces with different DocIDs", func() {
			province1 := domain.Province("Province 1")
			province2 := domain.Province("Province 2")

			err := upsert(ctx, province1)
			Expect(err).To(BeNil())
			err = upsert(ctx, province2)
			Expect(err).To(BeNil())

			var count int64
			conn.WithContext(ctx).Table("provinces").Count(&count)
			Expect(count).To(Equal(int64(2)))
		})

		It("should handle database errors gracefully", func() {
			// Create an invalid province with a very long name to trigger a database error
			invalidProvince := domain.Province("This is a very long province name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database tableThis is a very long province name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database table")

			err := upsert(ctx, invalidProvince)
			Expect(err).NotTo(BeNil())
		})
	})
})

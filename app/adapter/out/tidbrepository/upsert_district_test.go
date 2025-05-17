package tidbrepository

import (
	"context"
	"fmt"
	"testing"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

func TestUpsertDistrict(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpsertDistrict Suite")
}

var _ = Describe("UpsertDistrict", func() {
	var (
		ctx          context.Context
		conn         database.ConnectionFactory
		upsert       UpsertDistrict
		testDistrict domain.District
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
		upsert = NewUpsertDistrict(conn)
		testDistrict = domain.District("Test District")
	})

	AfterEach(func() {
		conn.WithContext(ctx).Exec("DELETE FROM districts")
	})

	Describe("UpsertDistrict", func() {
		It("should insert a new district", func() {
			err := upsert(ctx, testDistrict)
			Expect(err).To(BeNil())

			var savedDistrict table.District
			err = conn.WithContext(ctx).
				Table("districts").
				Where("document_id = ?", testDistrict.DocID(ctx).String()).
				First(&savedDistrict).Error
			Expect(err).To(BeNil())
			Expect(savedDistrict.Name).To(Equal(testDistrict.String()))
		})

		It("should create a new record when district name changes", func() {
			// First insert
			err := upsert(ctx, testDistrict)
			Expect(err).To(BeNil())

			var firstDistrict table.District
			err = conn.WithContext(ctx).
				Table("districts").
				Where("document_id = ?", testDistrict.DocID(ctx).String()).
				First(&firstDistrict).Error
			Expect(err).To(BeNil())
			firstID := firstDistrict.ID

			// Insert with different name
			updatedDistrict := domain.District("Updated District")
			err = upsert(ctx, updatedDistrict)
			Expect(err).To(BeNil())

			var secondDistrict table.District
			err = conn.WithContext(ctx).
				Table("districts").
				Where("document_id = ?", updatedDistrict.DocID(ctx).String()).
				First(&secondDistrict).Error
			Expect(err).To(BeNil())
			Expect(secondDistrict.ID).ToNot(Equal(firstID)) // Should be a new record
		})

		It("should handle multiple districts with different DocIDs", func() {
			district1 := domain.District("District 1")
			district2 := domain.District("District 2")

			err := upsert(ctx, district1)
			Expect(err).To(BeNil())
			err = upsert(ctx, district2)
			Expect(err).To(BeNil())

			var count int64
			conn.WithContext(ctx).Table("districts").Count(&count)
			Expect(count).To(Equal(int64(2)))
		})

		It("should handle database errors gracefully", func() {
			// Create an invalid district with a very long name to trigger a database error
			invalidDistrict := domain.District("This is a very long district name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database tableThis is a very long district name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database table")

			err := upsert(ctx, invalidDistrict)
			Expect(err).NotTo(BeNil())
		})

		It("should insert and retrieve an empty record", func() {
			// Insert empty district
			emptyDistrict := domain.District("")
			err := upsert(ctx, emptyDistrict)
			Expect(err).To(BeNil())

			// Get the DocID for debugging
			docID := emptyDistrict.DocID(ctx).String()
			tenantID := sharedcontext.TenantIDFromContext(ctx).String()
			fmt.Printf("\n=== EMPTY DISTRICT INSERTED ===\n")
			fmt.Printf("DocID: %s\n", docID)
			fmt.Printf("TenantID: %s\n", tenantID)

			// Try to retrieve the empty record
			var savedDistrict table.District
			err = conn.WithContext(ctx).
				Table("districts").
				Where("document_id = ?", docID).
				First(&savedDistrict).Error
			Expect(err).To(BeNil())
			Expect(savedDistrict.Name).To(Equal(""))
			Expect(savedDistrict.DocumentID).To(Equal(docID))
			Expect(savedDistrict.TenantID).To(Equal(tenantID))

			// Print the saved record for debugging
			fmt.Printf("\n=== RETRIEVED DISTRICT ===\n")
			fmt.Printf("Name: %s\n", savedDistrict.Name)
			fmt.Printf("DocID: %s\n", savedDistrict.DocumentID)
			fmt.Printf("TenantID: %s\n", savedDistrict.TenantID)
			fmt.Printf("========================\n\n")
		})
	})
})

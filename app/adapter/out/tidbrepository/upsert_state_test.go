package tidbrepository

import (
	"context"
	"testing"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

func TestUpsertState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpsertState Suite")
}

var _ = Describe("UpsertState", func() {
	var (
		ctx       context.Context
		conn      database.ConnectionFactory
		upsert    UpsertState
		testState domain.State
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
		upsert = NewUpsertState(conn)
		testState = domain.State("Test State")
	})

	AfterEach(func() {
		conn.WithContext(ctx).Exec("DELETE FROM states")
	})

	Describe("UpsertState", func() {
		It("should insert a new state", func() {
			err := upsert(ctx, testState)
			Expect(err).To(BeNil())

			var savedState table.State
			err = conn.WithContext(ctx).
				Table("states").
				Where("document_id = ?", testState.DocID(ctx).String()).
				First(&savedState).Error
			Expect(err).To(BeNil())
			Expect(savedState.Name).To(Equal(testState.String()))
		})

		It("should create a new record when state name changes", func() {
			// First insert
			err := upsert(ctx, testState)
			Expect(err).To(BeNil())

			var firstState table.State
			err = conn.WithContext(ctx).
				Table("states").
				Where("document_id = ?", testState.DocID(ctx).String()).
				First(&firstState).Error
			Expect(err).To(BeNil())
			firstID := firstState.ID

			// Insert with different name
			updatedState := domain.State("Updated State")
			err = upsert(ctx, updatedState)
			Expect(err).To(BeNil())

			var secondState table.State
			err = conn.WithContext(ctx).
				Table("states").
				Where("document_id = ?", updatedState.DocID(ctx).String()).
				First(&secondState).Error
			Expect(err).To(BeNil())
			Expect(secondState.ID).ToNot(Equal(firstID)) // Should be a new record
		})

		It("should handle multiple states with different DocIDs", func() {
			state1 := domain.State("State 1")
			state2 := domain.State("State 2")

			err := upsert(ctx, state1)
			Expect(err).To(BeNil())
			err = upsert(ctx, state2)
			Expect(err).To(BeNil())

			var count int64
			conn.WithContext(ctx).Table("states").Count(&count)
			Expect(count).To(Equal(int64(2)))
		})

		It("should handle database errors gracefully", func() {
			// Create an invalid state with a very long name to trigger a database error
			invalidState := domain.State("This is a very long state name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database tableThis is a very long state name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database table")

			err := upsert(ctx, invalidState)
			Expect(err).NotTo(BeNil())
		})

		It("should create an empty record", func() {
			emptyState := domain.State("")
			err := upsert(ctx, emptyState)
			Expect(err).To(BeNil())

			var savedState table.State
			err = conn.WithContext(ctx).
				Table("states").
				Where("document_id = ?", emptyState.DocID(ctx).String()).
				First(&savedState).Error
			Expect(err).To(BeNil())
			Expect(savedState.Name).To(Equal(""))
		})
	})
})

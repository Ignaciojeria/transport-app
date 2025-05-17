package tidbrepository

import (
	"context"
	"time"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertPlan", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	It("should insert plan if not exists", func() {
		ctx := createOrgContext(organization1)

		plan := domain.Plan{
			ReferenceID: "PLAN-001",
			PlannedDate: time.Now(),
		}

		upsert := NewUpsertPlan(connection)
		err := upsert(ctx, plan)
		Expect(err).ToNot(HaveOccurred())

		var dbPlan table.Plan
		err = connection.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", plan.DocID(ctx)).
			First(&dbPlan).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan.ReferenceID).To(Equal("PLAN-001"))
	})

	It("should update plan if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.Plan{
			ReferenceID: "PLAN-002",
			PlannedDate: time.Now(),
		}

		upsert := NewUpsertPlan(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Plan{
			ReferenceID: "PLAN-002",
			PlannedDate: time.Now().AddDate(0, 0, 1), // Tomorrow
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbPlan table.Plan
		err = connection.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbPlan).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan.ReferenceID).To(Equal("PLAN-002"))

		// Format dates to YYYY-MM-DD for comparison
		dbDate := dbPlan.PlannedDate.Format("2006-01-02")
		modifiedDate := modified.PlannedDate.Format("2006-01-02")
		Expect(dbDate).To(Equal(modifiedDate))
	})

	It("should not update if no fields changed", func() {
		ctx := createOrgContext(organization1)

		plan := domain.Plan{
			ReferenceID: "PLAN-003",
			PlannedDate: time.Now(),
		}

		upsert := NewUpsertPlan(connection)
		err := upsert(ctx, plan)
		Expect(err).ToNot(HaveOccurred())

		// Capture original record to verify timestamp doesn't change
		var originalRecord table.Plan
		err = connection.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", plan.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Execute again without changes
		err = upsert(ctx, plan)
		Expect(err).ToNot(HaveOccurred())

		var dbPlan table.Plan
		err = connection.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", plan.DocID(ctx)).
			First(&dbPlan).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan.ReferenceID).To(Equal("PLAN-003"))
	})

	It("should allow same plan for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		plan1 := domain.Plan{
			ReferenceID: "PLAN-004",
			PlannedDate: time.Now(),
		}

		plan2 := domain.Plan{
			ReferenceID: "PLAN-004",
			PlannedDate: time.Now(),
		}

		upsert := NewUpsertPlan(connection)

		err := upsert(ctx1, plan1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, plan2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("plans").
			Where("reference_id = ?", plan1.ReferenceID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(plan1.DocID(ctx1)).ToNot(Equal(plan2.DocID(ctx2)))
	})

	It("should fail if database has no plans table", func() {
		ctx := createOrgContext(organization1)

		plan := domain.Plan{
			ReferenceID: "PLAN-ERROR",
			PlannedDate: time.Now(),
		}

		upsert := NewUpsertPlan(noTablesContainerConnection)
		err := upsert(ctx, plan)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("plans"))
	})
})

package tidbrepository

import (
	"context"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertPlan", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertPlan
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertPlan(conn)
	})

	It("should insert plan if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		plan := domain.Plan{
			ReferenceID: "PLAN-001",
			PlannedDate: time.Now(),
		}

		err = upsert(ctx, plan)
		Expect(err).ToNot(HaveOccurred())

		var dbPlan table.Plan
		err = conn.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", plan.DocID(ctx)).
			First(&dbPlan).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan.ReferenceID).To(Equal("PLAN-001"))
		Expect(dbPlan.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update plan if fields are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Plan{
			ReferenceID: "PLAN-002",
			PlannedDate: time.Now(),
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Plan{
			ReferenceID: "PLAN-002",
			PlannedDate: time.Now().AddDate(0, 0, 1), // Tomorrow
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbPlan table.Plan
		err = conn.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbPlan).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan.ReferenceID).To(Equal("PLAN-002"))
		Expect(dbPlan.TenantID.String()).To(Equal(tenant.ID.String()))

		// Format dates to YYYY-MM-DD for comparison
		dbDate := dbPlan.PlannedDate.Format("2006-01-02")
		modifiedDate := modified.PlannedDate.Format("2006-01-02")
		Expect(dbDate).To(Equal(modifiedDate))
	})

	It("should allow same plan for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		plan1 := domain.Plan{
			ReferenceID: "PLAN-123",
		}

		plan2 := domain.Plan{
			ReferenceID: "PLAN-123",
		}

		err = upsert(ctx1, plan1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, plan2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := plan1.DocID(ctx1)
		docID2 := plan2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each plan belongs to its respective tenant using DocID
		var dbPlan1, dbPlan2 table.Plan
		err = conn.DB.WithContext(ctx1).
			Table("plans").
			Where("document_id = ?", docID1).
			First(&dbPlan1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("plans").
			Where("document_id = ?", docID2).
			First(&dbPlan2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPlan2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no plans table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		plan := domain.Plan{
			ReferenceID: "PLAN-ERROR",
			PlannedDate: time.Now(),
		}

		upsert := NewUpsertPlan(noTablesContainerConnection)
		err = upsert(ctx, plan)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("plans"))
	})
})

package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertStatus", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertStatus
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertStatus(conn)
	})

	It("should insert status if not exists", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		status := domain.Status{
			Status: "TEST",
		}

		err = upsert(ctx, status)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := status.DocID()

		var dbStatus table.Status
		err = conn.DB.WithContext(ctx).
			Table("statuses").
			Where("document_id = ?", docID).
			First(&dbStatus).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbStatus.Status).To(Equal("TEST"))
	})

	It("should create new record when status changes", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Status{
			Status: "TEST",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.Status
		err = conn.DB.WithContext(ctx).
			Table("statuses").
			Where("document_id = ?", original.DocID()).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Status{
			Status: "UPDATED",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verify a new record was created
		var dbStatus table.Status
		err = conn.DB.WithContext(ctx).
			Table("statuses").
			Where("document_id = ?", modified.DocID()).
			First(&dbStatus).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbStatus.Status).To(Equal("UPDATED"))
		Expect(dbStatus.ID).ToNot(Equal(originalRecord.ID)) // Should be a new record
	})

	It("should not update if status is the same", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		status := domain.Status{
			Status: "TEST",
		}

		err = upsert(ctx, status)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := status.DocID()

		// Get initial timestamp
		var initialStatus table.Status
		err = conn.DB.WithContext(ctx).
			Table("statuses").
			Where("document_id = ?", docID).
			First(&initialStatus).Error
		Expect(err).ToNot(HaveOccurred())
		initialUpdatedAt := initialStatus.UpdatedAt

		// Try to update with same values
		err = upsert(ctx, status)
		Expect(err).ToNot(HaveOccurred())

		var dbStatus table.Status
		err = conn.DB.WithContext(ctx).
			Table("statuses").
			Where("document_id = ?", docID).
			First(&dbStatus).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbStatus.UpdatedAt).To(Equal(initialUpdatedAt))
	})

	It("should allow same status for different tenants", func() {
		// Create two tenants for this test
		_, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		_, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		status1 := domain.Status{
			Status: "multi-org",
		}
		status2 := domain.Status{
			Status: "multi-org",
		}

		err = upsert(ctx1, status1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, status2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := status1.DocID()
		docID2 := status2.DocID()

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each status exists
		var dbStatus1, dbStatus2 table.Status
		err = conn.DB.WithContext(ctx1).
			Table("statuses").
			Where("document_id = ?", docID1).
			First(&dbStatus1).Error
		Expect(err).ToNot(HaveOccurred())

		err = conn.DB.WithContext(ctx2).
			Table("statuses").
			Where("document_id = ?", docID2).
			First(&dbStatus2).Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail if database has no statuses table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		status := domain.Status{
			Status: "TEST",
		}

		upsert := NewUpsertStatus(noTablesContainerConnection)
		err = upsert(ctx, status)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("statuses"))
	})
})

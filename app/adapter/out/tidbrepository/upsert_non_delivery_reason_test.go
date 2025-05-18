package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertNonDeliveryReason", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert non delivery reason if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		reason := domain.NonDeliveryReason{
			ReferenceID: "REF-001",
			Reason:      "NO_ONE_HOME",
			Details:     "Cliente no estaba presente",
		}

		upsert := NewUpsertNonDeliveryReason(conn)
		err = upsert(ctx, reason)
		Expect(err).ToNot(HaveOccurred())

		var dbReason table.NonDeliveryReason
		err = conn.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", reason.DocID(ctx)).
			First(&dbReason).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason.Reason).To(Equal("NO_ONE_HOME"))
		Expect(dbReason.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should create new record when fields change", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.NonDeliveryReason{
			ReferenceID: "REF-002",
			Reason:      "BAD_ADDRESS",
			Details:     "Dirección no existe",
		}
		upsert := NewUpsertNonDeliveryReason(conn)
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.NonDeliveryReason
		err = conn.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", original.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalRecord.TenantID.String()).To(Equal(tenant.ID.String()))

		modified := domain.NonDeliveryReason{
			ReferenceID: "REF-002",
			Reason:      "WRONG_PHONE",
			Details:     "Número incorrecto",
		}
		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbReason table.NonDeliveryReason
		err = conn.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbReason).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason.Reason).To(Equal("WRONG_PHONE"))
		Expect(dbReason.TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(dbReason.ID).To(Equal(originalRecord.ID)) // Should be the same record since ReferenceID is the same
	})

	It("should not create new record if no fields changed", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		reason := domain.NonDeliveryReason{
			ReferenceID: "REF-003",
			Reason:      "INCOMPLETE_ADDRESS",
		}

		upsert := NewUpsertNonDeliveryReason(conn)
		err = upsert(ctx, reason)
		Expect(err).ToNot(HaveOccurred())

		var original table.NonDeliveryReason
		err = conn.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", reason.DocID(ctx)).
			First(&original).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(original.TenantID.String()).To(Equal(tenant.ID.String()))

		// Ejecutar nuevamente sin cambios
		err = upsert(ctx, reason)
		Expect(err).ToNot(HaveOccurred())

		var dbReason table.NonDeliveryReason
		err = conn.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", reason.DocID(ctx)).
			First(&dbReason).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason.TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(dbReason.UpdatedAt).To(Equal(original.UpdatedAt))
	})

	It("should allow same reason for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		reason1 := domain.NonDeliveryReason{
			ReferenceID: "REF-004",
			Reason:      "SAME_REASON",
			Details:     "Same reason for different tenants",
		}
		reason2 := domain.NonDeliveryReason{
			ReferenceID: "REF-004",
			Reason:      "SAME_REASON",
			Details:     "Same reason for different tenants",
		}

		upsert := NewUpsertNonDeliveryReason(conn)

		err = upsert(ctx1, reason1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, reason2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := reason1.DocID(ctx1)
		docID2 := reason2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each reason belongs to its respective tenant using DocID
		var dbReason1, dbReason2 table.NonDeliveryReason
		err = conn.DB.WithContext(ctx1).
			Table("non_delivery_reasons").
			Where("document_id = ?", docID1).
			First(&dbReason1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("non_delivery_reasons").
			Where("document_id = ?", docID2).
			First(&dbReason2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no non_delivery_reasons table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		reason := domain.NonDeliveryReason{
			ReferenceID: "REF-999",
			Reason:      "DATABASE_MISSING",
		}
		upsert := NewUpsertNonDeliveryReason(noTablesContainerConnection)
		err = upsert(ctx, reason)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("non_delivery_reasons"))
	})
})

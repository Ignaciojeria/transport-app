package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertState", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertState
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertState(conn)
	})

	It("should insert state if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		state := domain.State("Test State")

		err = upsert(ctx, state)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := state.DocID(ctx)

		var dbState table.State
		err = conn.DB.WithContext(ctx).
			Table("states").
			Where("document_id = ?", docID).
			First(&dbState).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbState.Name).To(Equal("Test State"))
		Expect(dbState.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should create new record when name changes", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.State("Nombre Original")

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.State
		err = conn.DB.WithContext(ctx).
			Table("states").
			Where("document_id = ?", original.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalRecord.TenantID.String()).To(Equal(tenant.ID.String()))

		modified := domain.State("Nombre Modificado")

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verify a new record was created
		var dbState table.State
		err = conn.DB.WithContext(ctx).
			Table("states").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbState).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbState.Name).To(Equal("Nombre Modificado"))
		Expect(dbState.TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(dbState.ID).ToNot(Equal(originalRecord.ID)) // Should be a new record
	})

	It("should not update if state is the same", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		state := domain.State("Test State")

		err = upsert(ctx, state)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := state.DocID(ctx)

		// Get initial timestamp
		var initialState table.State
		err = conn.DB.WithContext(ctx).
			Table("states").
			Where("document_id = ?", docID).
			First(&initialState).Error
		Expect(err).ToNot(HaveOccurred())
		initialUpdatedAt := initialState.UpdatedAt

		// Try to update with same values
		err = upsert(ctx, state)
		Expect(err).ToNot(HaveOccurred())

		var dbState table.State
		err = conn.DB.WithContext(ctx).
			Table("states").
			Where("document_id = ?", docID).
			First(&dbState).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbState.UpdatedAt).To(Equal(initialUpdatedAt))
		Expect(dbState.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same state for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		state1 := domain.State("multi-org")
		state2 := domain.State("multi-org")

		err = upsert(ctx1, state1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, state2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := state1.DocID(ctx1)
		docID2 := state2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each state belongs to its respective tenant using DocID
		var dbState1, dbState2 table.State
		err = conn.DB.WithContext(ctx1).
			Table("states").
			Where("document_id = ?", docID1).
			First(&dbState1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbState1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("states").
			Where("document_id = ?", docID2).
			First(&dbState2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbState2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no states table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		state := domain.State("Test State")

		upsert := NewUpsertState(noTablesContainerConnection)
		err = upsert(ctx, state)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("states"))
	})
})

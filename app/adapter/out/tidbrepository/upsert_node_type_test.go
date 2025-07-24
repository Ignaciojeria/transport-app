package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertNodeType", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert node type if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		nodeType := domain.NodeType{
			Value: "pickup",
		}

		upsert := NewUpsertNodeType(conn, nil)
		err = upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("pickup"))
		Expect(dbNodeType.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should create new record when value changes", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.NodeType{
			Value: "original",
		}

		upsert := NewUpsertNodeType(conn, nil)
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.NodeType
		err = conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", original.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalRecord.TenantID.String()).To(Equal(tenant.ID.String()))

		modified := domain.NodeType{
			Value: "updated",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("updated"))
		Expect(dbNodeType.TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(dbNodeType.ID).ToNot(Equal(originalRecord.ID)) // Should be a new record
	})

	It("should not create new record if no fields changed", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		nodeType := domain.NodeType{
			Value: "no-change",
		}

		upsert := NewUpsertNodeType(conn, nil)
		err = upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		// Capture original record to verify timestamp doesn't change
		var originalRecord table.NodeType
		err = conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalRecord.TenantID.String()).To(Equal(tenant.ID.String()))

		// Ejecutar nuevamente sin cambios
		err = upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("no-change"))
		Expect(dbNodeType.TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(dbNodeType.UpdatedAt).To(Equal(originalRecord.UpdatedAt)) // Verify timestamp didn't change
	})

	It("should allow same node type for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		nodeType1 := domain.NodeType{
			Value: "multi-org",
		}
		nodeType2 := domain.NodeType{
			Value: "multi-org",
		}

		upsert := NewUpsertNodeType(conn, nil)

		err = upsert(ctx1, nodeType1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, nodeType2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := nodeType1.DocID(ctx1)
		docID2 := nodeType2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each node type belongs to its respective tenant using DocID
		var dbNodeType1, dbNodeType2 table.NodeType
		err = conn.DB.WithContext(ctx1).
			Table("node_types").
			Where("document_id = ?", docID1).
			First(&dbNodeType1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("node_types").
			Where("document_id = ?", docID2).
			First(&dbNodeType2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should generate predictable DocID with empty value", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		nodeType := domain.NodeType{
			Value: "",
		}

		upsert := NewUpsertNodeType(conn, nil)
		err = upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(BeEmpty())
		Expect(dbNodeType.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should fail if database has no node_types table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		nodeType := domain.NodeType{
			Value: "error-case",
		}

		upsert := NewUpsertNodeType(noTablesContainerConnection, nil)
		err = upsert(ctx, nodeType)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("node_types"))
	})
})

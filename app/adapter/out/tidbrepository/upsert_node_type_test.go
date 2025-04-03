package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertNodeType", func() {

	It("should insert node type if not exists", func() {
		ctx := context.Background()

		nodeType := domain.NodeType{
			Value:        "pickup",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID()).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("pickup"))
	})

	It("should update node type if fields are different", func() {
		ctx := context.Background()

		original := domain.NodeType{
			Value:        "original",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeType{
			Value:        "updated",
			Organization: organization1,
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", modified.DocID()).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("updated"))
	})

	It("should not update if no fields changed", func() {
		ctx := context.Background()

		nodeType := domain.NodeType{
			Value:        "no-change",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar nuevamente sin cambios
		err = upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID()).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("no-change"))
	})

	It("should allow same node type for different organizations", func() {
		ctx := context.Background()

		node1 := domain.NodeType{
			Value:        "shared-type",
			Organization: organization1,
		}

		node2 := domain.NodeType{
			Value:        "shared-type",
			Organization: organization2,
		}

		upsert := NewUpsertNodeType(connection)

		err := upsert(ctx, node1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx, node2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("value = ?", node1.Value).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should generate predictable DocID with empty value", func() {
		ctx := context.Background()

		nodeType := domain.NodeType{
			Value:        "",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID()).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(BeEmpty())
	})

	It("should fail if database has no node_types table", func() {
		ctx := context.Background()

		nodeType := domain.NodeType{
			Value:        "error-case",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(noTablesContainerConnection)
		err := upsert(ctx, nodeType)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("node_types"))
	})
})

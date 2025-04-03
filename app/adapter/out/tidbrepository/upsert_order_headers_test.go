package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestUpsertNodeType", func() {
	It("should insert node type if not exists", func() {
		ctx := context.Background()

		nt := domain.NodeType{
			Value:         "TestNodeType",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nt)
		Expect(err).ToNot(HaveOccurred())

		var result table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID()).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.DocumentID).To(Equal(string(nt.DocID())))
	})

	It("should update node type if exists and has changes", func() {
		ctx := context.Background()

		// Crear registro inicial
		nt := domain.NodeType{
			Value:         "UpdateableNodeType",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nt)
		Expect(err).ToNot(HaveOccurred())

		// Obtener el ID del registro creado para validación posterior
		var originalRecord table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID()).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		originalID := originalRecord.ID
		originalCreatedAt := originalRecord.CreatedAt

		// Modificar el nodo y actualizarlo
		updatedNt := domain.NodeType{
			Value:         "UpdateableNodeType-Modified",
			Organization: organization1,
		}

		err = upsert(ctx, updatedNt)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se actualizó el registro existente
		var updatedRecord table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID()).
			First(&updatedRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se mantiene el mismo ID y CreatedAt
		Expect(updatedRecord.ID).To(Equal(originalID))
		Expect(updatedRecord.CreatedAt).To(Equal(originalCreatedAt))

		// Verificar que se actualizó el nombre
		Expect(updatedRecord.Value).To(Equal("UpdateableNodeType-Modified"))
	})

	It("should not update node type if there are no changes", func() {
		ctx := context.Background()

		// Crear registro inicial
		nt := domain.NodeType{
			Value:         "UnchangedNodeType",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nt)
		Expect(err).ToNot(HaveOccurred())

		// Obtener el registro original para comparar después
		var originalRecord table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID()).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Intentar actualizar con los mismos datos
		err = upsert(ctx, nt)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que no hubo cambios en el registro
		var updatedRecord table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID()).
			First(&updatedRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Los campos deben mantenerse iguales
		Expect(updatedRecord).To(Equal(originalRecord))
	})

	It("should fail when trying to upsert node type in DB without tables", func() {
		ctx := context.Background()

		nt := domain.NodeType{
			Value:        "NodeTypeNoTables",
			Organization: organization1,
		}

		upsert := NewUpsertNodeType(noTablesContainerConnection)
		err := upsert(ctx, nt)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("node_types"))
	})
})

package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertOrderType", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		// Limpia la tabla antes de cada test
		err := connection.DB.Exec("DELETE FROM order_types").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert order type if it does not exist", func() {
		ot := domain.OrderType{
			Type:         "retail",
			Description:  "Entrega a cliente final",
			Organization: organization1,
		}

		upsert := NewUpsertOrderType(connection)
		err := upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		var result table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID()).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.Type).To(Equal("retail"))
		Expect(result.Description).To(Equal("Entrega a cliente final"))
	})

	It("should update order type if description is different", func() {
		ot := domain.OrderType{
			Type:         "b2b",
			Description:  "Original",
			Organization: organization1,
		}

		upsert := NewUpsertOrderType(connection)
		err := upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		// Actualizamos la descripción
		otUpdated := domain.OrderType{
			Type:         "b2b", // misma referencia
			Description:  "Modificada",
			Organization: organization1,
		}

		err = upsert(ctx, otUpdated)
		Expect(err).ToNot(HaveOccurred())

		var result table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID()).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.Description).To(Equal("Modificada"))
	})

	It("should not update if data has not changed", func() {
		ot := domain.OrderType{
			Type:         "express",
			Description:  "Entrega en 1 hora",
			Organization: organization1,
		}

		upsert := NewUpsertOrderType(connection)
		err := upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		// Ejecutamos con los mismos datos
		err = upsert(ctx, ot)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID()).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1))) // sigue habiendo solo uno
	})

	It("should fail if table does not exist", func() {
		ot := domain.OrderType{
			Type:         "ghost",
			Description:  "Fallará",
			Organization: organization1,
		}

		upsert := NewUpsertOrderType(noTablesContainerConnection)
		err := upsert(ctx, ot)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_types"))
	})
})

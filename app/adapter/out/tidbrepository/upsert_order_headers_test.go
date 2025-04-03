package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestUpsertOrderHeaders", func() {
	It("should insert order header if not exists", func() {
		// Preparar contexto y limpiar la tabla
		ctx := context.Background()

		// Crear datos de prueba
		h := domain.Headers{
			Commerce:     "tienda-123",
			Consumer:     "cliente-xyz",
			Organization: organization1,
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que el registro existe
		var result table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID()).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.DocumentID).To(Equal(h.DocID()))
	})

	It("should fail to insert order header if organization is missing", func() {
		ctx := context.Background()

		// Organización vacía
		h := domain.Headers{
			Commerce:     "tienda-123",
			Consumer:     "cliente-xyz",
			Organization: domain.Organization{}, // ID=0, Country=""
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).To(HaveOccurred()) // Debe fallar
	})

	It("should not insert the same order header twice for the same organization", func() {
		ctx := context.Background()

		// Crear cabecera de prueba
		h := domain.Headers{
			Commerce:     "tienda-repetida",
			Consumer:     "cliente-repetido",
			Organization: organization1,
		}

		upsert := NewUpsertOrderHeaders(connection)

		// Primera inserción
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		// Segunda inserción (debe ser ignorada sin fallar)
		err = upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que solo hay una fila
		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID()).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should insert the same order header in different organizations", func() {
		ctx := context.Background()

		// Crear cabecera con misma referencia, pero diferente organización
		h1 := domain.Headers{
			Commerce:     "tienda-repetida",
			Consumer:     "cliente-repetido",
			Organization: organization1,
		}

		h2 := domain.Headers{
			Commerce:     "tienda-repetida",
			Consumer:     "cliente-repetido",
			Organization: organization2,
		}

		upsert := NewUpsertOrderHeaders(connection)

		// Insertar en la primera organización
		err := upsert(ctx, h1)
		Expect(err).ToNot(HaveOccurred())

		// Insertar en la segunda organización
		err = upsert(ctx, h2)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que ambos registros fueron insertados correctamente
		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id IN (?, ?)", h1.DocID(), h2.DocID()).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should fail when trying to upsert order header in DB without tables", func() {
		ctx := context.Background()

		h := domain.Headers{
			Commerce:     "tienda-sin-tablas",
			Consumer:     "cliente-sin-tablas",
			Organization: organization1,
		}

		upsert := NewUpsertOrderHeaders(noTablesContainerConnection)
		err := upsert(ctx, h)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_headers"))
	})

})

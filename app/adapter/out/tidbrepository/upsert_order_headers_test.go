package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestUpsertOrderHeaders", func() {
	// Scenario 1: Insert order header if not exists (ya cubierto por tu test existente)
	It("should insert order header if not exists", func() {
		ctx := context.Background()

		h := domain.Headers{
			Commerce:     "tienda-123",
			Consumer:     "cliente-xyz",
			Organization: organization1,
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		var result table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID()). 
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.DocumentID).To(Equal(string(h.DocID())))
	})

	// Scenario 2: Do nothing if order header already exists
	It("should do nothing if order header already exists", func() {
		ctx := context.Background()

		h := domain.Headers{
			Commerce:     "tienda-existente",
			Consumer:     "cliente-existente",
			Organization: organization1,
		}

		upsert := NewUpsertOrderHeaders(connection)
		
		// Primera inserción
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())
		
		// Obtener timestamp de creación para comparar después
		var firstRecord table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID()).
			First(&firstRecord).Error
		Expect(err).ToNot(HaveOccurred())
		
		// Guardar el estado inicial
		initialTimestamp := firstRecord.CreatedAt
		initialID := firstRecord.ID
		
		// Intentar insertar el mismo registro de nuevo
		err = upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred()) // No debería dar error
		
		// Verificar que no se creó un nuevo registro y que el existente no fue modificado
		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID()).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1))) // Solo debe existir un registro
		
		// Verificar que los datos no cambiaron
		var updatedRecord table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID()).
			First(&updatedRecord).Error
		Expect(err).ToNot(HaveOccurred())
		
		Expect(updatedRecord.ID).To(Equal(initialID))
		Expect(updatedRecord.CreatedAt).To(Equal(initialTimestamp))
	})

	// Scenario 3: Return error when DB query fails
	// Para este escenario usaremos la conexión sin tablas que ya tienes
	It("should return error when DB query fails", func() {
		ctx := context.Background()

		// Modificamos la conexión para que falle al hacer la consulta
		// Esto simulará un error durante la consulta inicial
		badConnection := connection
		badConnection.DB = noTablesContainerConnection.DB
		
		h := domain.Headers{
			Commerce:     "tienda-error",
			Consumer:     "cliente-error",
			Organization: organization1,
		}
		
		upsert := NewUpsertOrderHeaders(badConnection)
		err := upsert(ctx, h)
		
		// Verificar que retorna un error
		Expect(err).To(HaveOccurred())
	})

	// Escenario adicional: falla si la organización no es válida (ya cubierto por tu test existente)
	It("should fail to insert order header if organization is missing", func() {
		ctx := context.Background()

		h := domain.Headers{
			Commerce:     "tienda-123",
			Consumer:     "cliente-xyz",
			Organization: domain.Organization{}, // ID=0, Country=""
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).To(HaveOccurred())
	})

	// Escenario adicional: prueba de concurrencia con dos organizaciones diferentes (ya cubierto por tu test existente)
	It("should insert the same order header in different organizations", func() {
		ctx := context.Background()

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

		err := upsert(ctx, h1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx, h2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id IN (?, ?)", h1.DocID(), h2.DocID()).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	// Escenario adicional: error con conexión a BD sin tablas (ya cubierto por tu test existente)
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
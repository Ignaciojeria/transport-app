package tidbrepository

import (
	"context"
	"strconv"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("TestUpsertOrderHeaders", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	// Scenario 1: Insert order header if not exists
	It("should insert order header if not exists", func() {
		ctx := createOrgContext(organization1)

		h := domain.Headers{
			Commerce: "tienda-123",
			Consumer: "cliente-xyz",
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		var result table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID(ctx)).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.DocumentID).To(Equal(string(h.DocID(ctx))))
	})

	// Scenario 2: Do nothing if order header already exists
	It("should do nothing if order header already exists", func() {
		ctx := createOrgContext(organization1)

		h := domain.Headers{
			Commerce: "tienda-existente",
			Consumer: "cliente-existente",
		}

		upsert := NewUpsertOrderHeaders(connection)

		// Primera inserción
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		// Obtener timestamp de creación para comparar después
		var firstRecord table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID(ctx)).
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
			Where("document_id = ?", h.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1))) // Solo debe existir un registro

		// Verificar que los datos no cambiaron
		var updatedRecord table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID(ctx)).
			First(&updatedRecord).Error
		Expect(err).ToNot(HaveOccurred())

		Expect(updatedRecord.ID).To(Equal(initialID))
		Expect(updatedRecord.CreatedAt).To(Equal(initialTimestamp))
	})

	// Scenario 3: Return error when DB query fails
	It("should return error when DB query fails", func() {
		ctx := createOrgContext(organization1)

		// Modificamos la conexión para que falle al hacer la consulta
		// Esto simulará un error durante la consulta inicial
		badConnection := connection
		badConnection.DB = noTablesContainerConnection.DB

		h := domain.Headers{
			Commerce: "tienda-error",
			Consumer: "cliente-error",
		}

		upsert := NewUpsertOrderHeaders(badConnection)
		err := upsert(ctx, h)

		// Verificar que retorna un error
		Expect(err).To(HaveOccurred())
	})

	// Escenario adicional: falla si la organización no es válida
	It("should fail to insert order header if organization is missing from context", func() {
		// Use a context without organization information
		ctx := context.Background()

		h := domain.Headers{
			Commerce: "tienda-123",
			Consumer: "cliente-xyz",
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).To(HaveOccurred())
	})

	// Escenario adicional: prueba de concurrencia con dos organizaciones diferentes
	It("should insert the same order header in different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		h1 := domain.Headers{
			Commerce: "tienda-repetida",
			Consumer: "cliente-repetido",
		}

		h2 := domain.Headers{
			Commerce: "tienda-repetida",
			Consumer: "cliente-repetido",
		}

		upsert := NewUpsertOrderHeaders(connection)

		err := upsert(ctx1, h1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, h2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("order_headers").
			Where("document_id IN (?, ?)", h1.DocID(ctx1), h2.DocID(ctx2)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(h1.DocID(ctx1)).ToNot(Equal(h2.DocID(ctx2)))
	})

	// Escenario adicional: error con conexión a BD sin tablas
	It("should fail when trying to upsert order header in DB without tables", func() {
		ctx := createOrgContext(organization1)

		h := domain.Headers{
			Commerce: "tienda-sin-tablas",
			Consumer: "cliente-sin-tablas",
		}

		upsert := NewUpsertOrderHeaders(noTablesContainerConnection)
		err := upsert(ctx, h)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_headers"))
	})
})

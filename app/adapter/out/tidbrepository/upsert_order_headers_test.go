package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestUpsertOrderHeaders", func() {
	It("should insert order header if not exists", func() {
		// Preparar contexto y limpiar la tabla
		ctx := context.Background()

		// Crear datos de prueba
		h := domain.Headers{
			Commerce: "tienda-123",
			Consumer: "cliente-xyz",
			Organization: domain.Organization{
				ID:      1,
				Country: countries.CL,
			},
		}

		upsert := NewUpsertOrderHeaders(connection)
		err := upsert(ctx, h)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que el registro existe
		var result table.OrderHeaders
		err = connection.DB.WithContext(ctx).
			Table("order_headers").
			Where("reference_id = ?", h.ReferenceID()).
			First(&result).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(result.ReferenceID).To(Equal(h.ReferenceID()))
	})
})

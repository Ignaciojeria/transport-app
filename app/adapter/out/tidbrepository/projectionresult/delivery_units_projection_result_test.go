package projectionresult

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProjectionResult(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProjectionResult Suite")
}

var _ = Describe("DeliveryUnitsProjectionResults", func() {
	It("should reverse the order of results when using Reversed()", func() {
		// Crear un slice de resultados de prueba
		results := DeliveryUnitsProjectionResults{
			{
				ID:               1,
				OrderReferenceID: "order1",
			},
			{
				ID:               2,
				OrderReferenceID: "order2",
			},
			{
				ID:               3,
				OrderReferenceID: "order3",
			},
		}

		// Verificar orden original
		Expect(results[0].OrderReferenceID).To(Equal("order1"))
		Expect(results[1].OrderReferenceID).To(Equal("order2"))
		Expect(results[2].OrderReferenceID).To(Equal("order3"))

		// Revertir el orden
		reversed := results.Reversed()

		// Verificar que el orden se invirtió correctamente
		Expect(reversed).To(HaveLen(3))
		Expect(reversed[0].OrderReferenceID).To(Equal("order3"))
		Expect(reversed[1].OrderReferenceID).To(Equal("order2"))
		Expect(reversed[2].OrderReferenceID).To(Equal("order1"))

		// Verificar que el slice original no se modificó
		Expect(results[0].OrderReferenceID).To(Equal("order1"))
		Expect(results[1].OrderReferenceID).To(Equal("order2"))
		Expect(results[2].OrderReferenceID).To(Equal("order3"))
	})
})

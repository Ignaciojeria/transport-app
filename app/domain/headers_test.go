package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Headers", func() {
	var org1 = Organization{ID: 1, Country: countries.CL}
	var org2 = Organization{ID: 2, Country: countries.AR}

	Describe("DocID", func() {
		It("should generate unique ID based on Organization, Commerce and Consumer", func() {
			headers1 := Headers{
				Organization: org1,
				Commerce:     "store-1",
				Consumer:     "customer-1",
			}
			headers2 := Headers{
				Organization: org1,
				Commerce:     "store-2",
				Consumer:     "customer-1",
			}
			headers3 := Headers{
				Organization: org2,
				Commerce:     "store-1",
				Consumer:     "customer-1",
			}

			// Mismo Commerce y Consumer pero diferente organization -> diferente DocID
			Expect(headers1.DocID()).To(Equal(Hash(org1, "store-1", "customer-1")))
			Expect(headers1.DocID()).ToNot(Equal(headers2.DocID()))
			Expect(headers1.DocID()).ToNot(Equal(headers3.DocID()))
		})

		It("should generate different IDs for different Commerce values", func() {
			headers1 := Headers{
				Organization: org1,
				Commerce:     "store-1",
				Consumer:     "customer-1",
			}
			headers2 := Headers{
				Organization: org1,
				Commerce:     "store-2",
				Consumer:     "customer-1",
			}

			Expect(headers1.DocID()).ToNot(Equal(headers2.DocID()))
		})

		It("should generate different IDs for different Consumer values", func() {
			headers1 := Headers{
				Organization: org1,
				Commerce:     "store-1",
				Consumer:     "customer-1",
			}
			headers2 := Headers{
				Organization: org1,
				Commerce:     "store-1",
				Consumer:     "customer-2",
			}

			Expect(headers1.DocID()).ToNot(Equal(headers2.DocID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		var baseHeaders Headers

		BeforeEach(func() {
			// Configurar headers base para cada test
			baseHeaders = Headers{
				Organization: org1,
				Commerce:     "original-store",
				Consumer:     "original-customer",
			}
		})

		It("should update Consumer if new value is not empty", func() {
			newHeaders := Headers{Consumer: "new-customer"}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Consumer).To(Equal("new-customer"))
			Expect(updated.Organization).To(Equal(org1))         // No debería cambiar
			Expect(updated.Commerce).To(Equal("original-store")) // No debería cambiar
		})

		It("should not update Consumer if new value is empty", func() {
			newHeaders := Headers{Consumer: ""}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Consumer).To(Equal("original-customer")) // Debería mantener el valor original
		})

		It("should update Commerce if new value is not empty", func() {
			newHeaders := Headers{Commerce: "new-store"}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Commerce).To(Equal("new-store"))
			Expect(updated.Organization).To(Equal(org1))            // No debería cambiar
			Expect(updated.Consumer).To(Equal("original-customer")) // No debería cambiar
		})

		It("should not update Commerce if new value is empty", func() {
			newHeaders := Headers{Commerce: ""}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Commerce).To(Equal("original-store")) // Debería mantener el valor original
		})

		It("should update Organization if new Organization ID is not zero", func() {
			newHeaders := Headers{Organization: org2}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Organization).To(Equal(org2))
			Expect(updated.Commerce).To(Equal("original-store"))    // No debería cambiar
			Expect(updated.Consumer).To(Equal("original-customer")) // No debería cambiar
		})

		It("should not update Organization if new Organization ID is zero", func() {
			newHeaders := Headers{Organization: Organization{ID: 0}}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Organization).To(Equal(org1)) // Debería mantener el valor original
		})

		It("should update multiple fields at once", func() {
			newHeaders := Headers{
				Organization: org2,
				Commerce:     "new-store",
				Consumer:     "new-customer",
			}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Organization).To(Equal(org2))
			Expect(updated.Commerce).To(Equal("new-store"))
			Expect(updated.Consumer).To(Equal("new-customer"))
		})

		It("should update only non-zero and non-empty fields", func() {
			newHeaders := Headers{
				Organization: Organization{ID: 0}, // No debería actualizar
				Commerce:     "",                  // No debería actualizar
				Consumer:     "new-customer",
			}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Organization).To(Equal(org1))         // No debería cambiar
			Expect(updated.Commerce).To(Equal("original-store")) // No debería cambiar
			Expect(updated.Consumer).To(Equal("new-customer"))
		})

		It("should not change anything if all new values are empty or zero", func() {
			newHeaders := Headers{
				Organization: Organization{ID: 0},
				Commerce:     "",
				Consumer:     "",
			}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated).To(Equal(baseHeaders))
		})
	})
})

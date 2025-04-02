package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrderType", func() {

	var (
		org1 = Organization{ID: 1, Country: countries.CL}
		org2 = Organization{ID: 2, Country: countries.AR}
	)

	Describe("ReferenceID", func() {
		It("should generate different reference IDs for different organizations", func() {
			orderType1 := OrderType{Organization: org1, Type: "retail"}
			orderType2 := OrderType{Organization: org2, Type: "retail"}

			Expect(orderType1.DocID()).ToNot(Equal(orderType2.DocID()))
		})

		It("should generate different reference IDs for different types", func() {
			orderType1 := OrderType{Organization: org1, Type: "retail"}
			orderType2 := OrderType{Organization: org1, Type: "b2b"}

			Expect(orderType1.DocID()).ToNot(Equal(orderType2.DocID()))
		})

		It("should generate the same reference ID for same org and type", func() {
			orderType1 := OrderType{Organization: org1, Type: "express"}
			orderType2 := OrderType{Organization: org1, Type: "express"}

			Expect(orderType1.DocID()).To(Equal(orderType2.DocID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should return updated order type if type is different", func() {
			original := OrderType{
				Type:         "retail",
				Description:  "Entrega normal",
				Organization: org1,
			}
			input := OrderType{
				Type:        "express",
				Description: "Entrega normal", // igual
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("express"))
			Expect(updated.Description).To(Equal("Entrega normal"))
		})

		It("should return updated order type if description is different", func() {
			original := OrderType{
				Type:         "retail",
				Description:  "Descripción vieja",
				Organization: org1,
			}
			input := OrderType{
				Type:        "retail",
				Description: "Descripción nueva",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Description).To(Equal("Descripción nueva"))
		})

		It("should not mark as changed if values are the same", func() {
			original := OrderType{
				Type:         "retail",
				Description:  "Estándar",
				Organization: org1,
			}
			input := OrderType{
				Type:        "retail",
				Description: "Estándar",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore empty fields in input", func() {
			original := OrderType{
				Type:         "retail",
				Description:  "Estándar",
				Organization: org1,
			}
			input := OrderType{
				Type:        "", // no cambia
				Description: "",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should not update organization", func() {
			original := OrderType{
				Type:         "retail",
				Description:  "Estándar",
				Organization: org1,
			}
			input := OrderType{
				Type:         "retail",
				Description:  "Estándar",
				Organization: org2, // distinto pero no se debería aplicar
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated.Organization).To(Equal(org1))
		})
	})
})

package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrderType", func() {
	var (
		ctx1 context.Context
		ctx2 context.Context
	)

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate different reference IDs for different contexts", func() {
			orderType := OrderType{Type: "retail"}

			Expect(orderType.DocID(ctx1)).ToNot(Equal(orderType.DocID(ctx2)))
		})

		It("should generate different reference IDs for different types", func() {
			orderType1 := OrderType{Type: "retail"}
			orderType2 := OrderType{Type: "b2b"}

			Expect(orderType1.DocID(ctx1)).ToNot(Equal(orderType2.DocID(ctx1)))
		})

		It("should generate the same reference ID for same context and type", func() {
			orderType1 := OrderType{Type: "express"}
			orderType2 := OrderType{Type: "express"}

			Expect(orderType1.DocID(ctx1)).To(Equal(orderType2.DocID(ctx1)))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should return updated order type if description is different", func() {
			original := OrderType{
				Type:        "retail",
				Description: "Descripción vieja",
			}
			input := OrderType{
				Type:        "retail",
				Description: "Descripción nueva",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Description).To(Equal("Descripción nueva"))
			Expect(updated.Type).To(Equal("retail")) // Type no cambia
		})

		It("should not mark as changed if values are the same", func() {
			original := OrderType{
				Type:        "retail",
				Description: "Estándar",
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
				Type:        "retail",
				Description: "Estándar",
			}
			input := OrderType{
				Type:        "",
				Description: "",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore Type changes as per implementation", func() {
			original := OrderType{
				Type:        "retail",
				Description: "Estándar",
			}
			input := OrderType{
				Type:        "express", // Este cambio debe ser ignorado
				Description: "Estándar",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated.Type).To(Equal("retail")) // Type no cambia
		})

		It("should only update Description field", func() {
			original := OrderType{
				Type:        "retail",
				Description: "Descripción vieja",
			}
			input := OrderType{
				Type:        "express", // Este cambio debe ser ignorado
				Description: "Descripción nueva",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("retail")) // Type no cambia
			Expect(updated.Description).To(Equal("Descripción nueva"))
		})
	})
})

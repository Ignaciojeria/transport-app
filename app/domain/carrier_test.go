package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Carrier", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	Describe("UpdateIfChanged", func() {
		It("should update name if different and not empty", func() {
			original := Carrier{
				Name: "Old Carrier",
			}
			newCarrier := Carrier{
				Name: "New Carrier",
			}

			updated, changed := original.UpdateIfChanged(newCarrier)

			Expect(changed).To(BeTrue())
			Expect(updated.Name).To(Equal("New Carrier"))
		})

		It("should not update name if it's the same", func() {
			original := Carrier{
				Name: "Carrier X",
			}
			newCarrier := Carrier{
				Name: "Carrier X",
			}

			updated, changed := original.UpdateIfChanged(newCarrier)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should not update name if new name is empty", func() {
			original := Carrier{
				Name: "Carrier Z",
			}
			newCarrier := Carrier{
				Name: "",
			}

			updated, changed := original.UpdateIfChanged(newCarrier)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})
	})

	Describe("DocID", func() {
		It("should return hashed ID using tenant and national ID", func() {
			carrier := Carrier{
				NationalID: "12345678-9",
			}

			expectedHash := Hash(ctx, "12345678-9")
			Expect(carrier.DocID(ctx)).To(Equal(expectedHash))
		})
	})
})

package domain_test

import (
	. "transport-app/app/domain"

	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Carrier", func() {

	org := Organization{ID: 1, Country: countries.CL}

	Describe("UpdateIfChanged", func() {
		It("should update name if different and not empty", func() {
			original := Carrier{
				Name:         "Old Carrier",
				Organization: org,
			}
			newCarrier := Carrier{
				Name: "New Carrier",
			}

			updated, changed := original.UpdateIfChanged(newCarrier)

			Expect(changed).To(BeTrue())
			Expect(updated.Name).To(Equal("New Carrier"))
			Expect(updated.Organization).To(Equal(original.Organization)) // No cambia
		})

		It("should not update name if it's the same", func() {
			original := Carrier{
				Name:         "Carrier X",
				Organization: org,
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
				Name:         "Carrier Z",
				Organization: org,
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
		It("should return hashed ID using organization and national ID", func() {
			carrier := Carrier{
				Organization: org,
				NationalID:   "12345678-9",
			}

			expectedHash := Hash(org, "12345678-9")
			Expect(carrier.DocID()).To(Equal(expectedHash))
		})
	})
})

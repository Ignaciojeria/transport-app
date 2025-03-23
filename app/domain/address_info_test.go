package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressInfo", func() {

	var org1 = Organization{ID: 1, Country: countries.CL}

	Describe("ReferenceID", func() {
		It("should generate different reference IDs for different address lines", func() {
			addr1 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				AddressLine2: "Dpto 1402",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}
			addr2 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 5678", // distinto
				AddressLine2: "Dpto 1402",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}

			Expect(addr1.ReferenceID()).ToNot(Equal(addr2.ReferenceID()))
		})

		It("should generate same ReferenceID for same input", func() {
			addr1 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				AddressLine2: "Dpto 1402",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}
			addr2 := addr1

			Expect(addr1.ReferenceID()).To(Equal(addr2.ReferenceID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should mark as changed when district is updated", func() {
			original := AddressInfo{District: "Ñuñoa"}
			input := AddressInfo{District: "Providencia"}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.District).To(Equal("Providencia"))
		})

		It("should not change when all fields are the same", func() {
			original := AddressInfo{
				AddressLine1: "calle 123",
				District:     "Ñuñoa",
				State:        "Metropolitana",
			}
			input := AddressInfo{
				AddressLine1: "calle 123",
				District:     "Ñuñoa",
				State:        "Metropolitana",
			}

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore empty fields in input", func() {
			original := AddressInfo{
				AddressLine1: "calle 123",
				District:     "Ñuñoa",
			}
			input := AddressInfo{} // no hay cambios

			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})
	})

	Describe("Normalize", func() {
		It("should normalize casing and spacing", func() {
			addr := AddressInfo{
				AddressLine1:    "   AVENIDA   PROVIDENCIA   1234  ",
				AddressLine2:    "  Dpto  1402 ",
				AddressLine3:    "  Torre Norte ",
				ProviderAddress: "   AV PROVIDENCIA 1234 ",
				State:           "   Metropolitana ",
				Province:        " SANTIAGO  ",
				District:        "PROVIDENCIA   ",
			}

			addr.Normalize()

			Expect(addr.AddressLine1).To(Equal("avenida providencia 1234"))
			Expect(addr.AddressLine2).To(Equal("dpto 1402"))
			Expect(addr.AddressLine3).To(Equal("torre norte"))
			Expect(addr.ProviderAddress).To(Equal("av providencia 1234"))
			Expect(addr.State).To(Equal("metropolitana"))
			Expect(addr.Province).To(Equal("santiago"))
			Expect(addr.District).To(Equal("providencia"))
		})
	})
})

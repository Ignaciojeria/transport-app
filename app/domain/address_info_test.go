package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("AddressInfo", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	Describe("ReferenceID / DocID", func() {
		It("should generate different reference IDs for different address lines", func() {
			addr1 := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}
			addr2 := AddressInfo{
				AddressLine1: "Av Providencia 5678", // distinto
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}

			Expect(addr1.DocID(ctx)).ToNot(Equal(addr2.DocID(ctx)))
		})

		It("should generate same ReferenceID for same input", func() {
			addr1 := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}
			addr2 := addr1

			Expect(addr1.DocID(ctx)).To(Equal(addr2.DocID(ctx)))
		})

		It("should confirm hash only includes specified fields", func() {
			addr1 := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
				ZipCode:      "7500000",
			}
			addr2 := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
				ZipCode:      "7550000", // No debería afectar el hash
			}

			Expect(addr1.DocID(ctx)).To(Equal(addr2.DocID(ctx)))
		})
	})

	Describe("UpdateIfChanged", func() {
		var original AddressInfo

		BeforeEach(func() {
			original = AddressInfo{
				AddressLine1: "Av Providencia 1234",
				State:        "Metropolitana",
				Province:     "Santiago",
				District:     "Providencia",
				ZipCode:      "7500000",
				TimeZone:     "America/Santiago",
				Location:     orb.Point{-33.4372, -70.6506},
			}
		})

		It("should mark as changed when AddressLine1 is updated", func() {
			input := AddressInfo{AddressLine1: "Av Las Condes 5678"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine1).To(Equal("Av Las Condes 5678"))
		})

		It("should mark as changed when Location is updated", func() {
			newLocation := orb.Point{-33.4400, -70.6600}
			input := AddressInfo{Location: newLocation}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Location).To(Equal(newLocation))
		})

		It("should mark as changed when State is updated", func() {
			input := AddressInfo{State: "Valparaíso"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.State.String()).To(Equal("Valparaíso"))
		})

		It("should mark as changed when Province is updated", func() {
			input := AddressInfo{Province: "Valparaíso"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Province.String()).To(Equal("Valparaíso"))
		})

		It("should mark as changed when District is updated", func() {
			input := AddressInfo{District: "Las Condes"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.District.String()).To(Equal("Las Condes"))
		})

		It("should mark as changed when ZipCode is updated", func() {
			input := AddressInfo{ZipCode: "7550000"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.ZipCode).To(Equal("7550000"))
		})

		It("should mark as changed when TimeZone is updated", func() {
			input := AddressInfo{TimeZone: "America/Argentina/Buenos_Aires"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.TimeZone).To(Equal("America/Argentina/Buenos_Aires"))
		})

		It("should not change when all fields are the same", func() {
			input := original
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore empty fields in input", func() {
			input := AddressInfo{}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should update multiple fields at once", func() {
			input := AddressInfo{
				AddressLine1: "Av Las Condes 5678",
				District:     "Las Condes",
				ZipCode:      "7550000",
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine1).To(Equal("Av Las Condes 5678"))
			Expect(updated.District.String()).To(Equal("Las Condes"))
			Expect(updated.ZipCode).To(Equal("7550000"))
			Expect(updated.State).To(Equal(original.State))
		})
	})

	Describe("Normalize", func() {
		It("should normalize casing and spacing", func() {
			addr := AddressInfo{
				AddressLine1: "   AVENIDA   PROVIDENCIA   1234  ",
				State:        "   Metropolitana ",
				Province:     " SANTIAGO  ",
				District:     "PROVIDENCIA   ",
			}

			addr.ToLowerAndRemovePunctuation()

			Expect(addr.AddressLine1).To(Equal("avenida providencia 1234"))
			Expect(addr.State.String()).To(Equal("metropolitana"))
			Expect(addr.Province.String()).To(Equal("santiago"))
			Expect(addr.District.String()).To(Equal("providencia"))
		})
	})

	Describe("FullAddress", func() {
		It("should concatenate address fields correctly", func() {
			addr := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
				ZipCode:      "7500000",
			}

			fullAddress := addr.FullAddress()

			Expect(fullAddress).To(ContainSubstring("Av Providencia 1234"))
			Expect(fullAddress).To(ContainSubstring("Providencia"))
			Expect(fullAddress).To(ContainSubstring("Santiago"))
			Expect(fullAddress).To(ContainSubstring("Metropolitana"))
			Expect(fullAddress).To(ContainSubstring("7500000"))
		})

		It("should only include specified fields in the address format", func() {
			addr := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
				ZipCode:      "7500000",
			}

			fullAddress := addr.FullAddress()
			Expect(fullAddress).To(Equal("Av Providencia 1234, Providencia, Santiago, Metropolitana, 7500000"))
		})

		It("should correctly handle empty fields", func() {
			addr := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				State:        "Metropolitana",
			}

			fullAddress := addr.FullAddress()
			Expect(fullAddress).To(Equal("Av Providencia 1234, Providencia, Metropolitana"))
		})

		It("should return just AddressLine1 if other fields are empty", func() {
			addr := AddressInfo{AddressLine1: "Av Providencia 1234"}
			Expect(addr.FullAddress()).To(Equal("Av Providencia 1234"))
		})
	})

	Describe("concatenateWithCommas", func() {
		It("should concatenate multiple non-empty strings with commas", func() {
			Expect(concatenateWithCommas("uno", "dos", "tres")).To(Equal("uno, dos, tres"))
		})

		It("should skip empty strings", func() {
			Expect(concatenateWithCommas("uno", "", "tres")).To(Equal("uno, tres"))
		})

		It("should return a single value without comma if only one non-empty string", func() {
			Expect(concatenateWithCommas("", "", "único")).To(Equal("único"))
		})

		It("should return empty string if all inputs are empty", func() {
			Expect(concatenateWithCommas("", "", "")).To(BeEmpty())
		})

		It("should return empty string if no values are passed", func() {
			Expect(concatenateWithCommas()).To(BeEmpty())
		})
	})
})

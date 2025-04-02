package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("AddressInfo", func() {

	var org1 = Organization{ID: 1, Country: countries.CL}

	Describe("ReferenceID", func() {
		It("should generate different reference IDs for different address lines", func() {
			addr1 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}
			addr2 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 5678", // distinto
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}

			Expect(addr1.DocID()).ToNot(Equal(addr2.DocID()))
		})

		It("should generate same ReferenceID for same input", func() {
			addr1 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
			}
			addr2 := addr1

			Expect(addr1.DocID()).To(Equal(addr2.DocID()))
		})

		// Aunque AddressLine2 ya no es parte de la estructura, este test demuestra
		// que el diseño es correcto y solo se incluyen los campos relevantes en el hash
		It("should confirm hash only includes specified fields", func() {
			addr1 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
				ZipCode:      "7500000", // No debería afectar el hash
			}
			addr2 := AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				Province:     "Santiago",
				State:        "Metropolitana",
				ZipCode:      "7550000", // Diferente pero no debería afectar el hash
			}

			Expect(addr1.DocID()).To(Equal(addr2.DocID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		var original AddressInfo

		BeforeEach(func() {
			original = AddressInfo{
				Organization: org1,
				AddressLine1: "Av Providencia 1234",
				State:        "Metropolitana",
				Province:     "Santiago",
				District:     "Providencia",
				ZipCode:      "7500000",
				TimeZone:     "America/Santiago",
				Location:     orb.Point{-33.4372, -70.6506}, // Coordenadas de ejemplo
			}
		})

		It("should mark as changed when AddressLine1 is updated", func() {
			input := AddressInfo{AddressLine1: "Av Las Condes 5678"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine1).To(Equal("Av Las Condes 5678"))
		})

		It("should mark as changed when Location is updated", func() {
			newLocation := orb.Point{-33.4400, -70.6600} // Coordenadas diferentes
			input := AddressInfo{Location: newLocation}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Location).To(Equal(newLocation))
		})

		It("should mark as changed when State is updated", func() {
			input := AddressInfo{State: "Valparaíso"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.State).To(Equal("Valparaíso"))
		})

		It("should mark as changed when Province is updated", func() {
			input := AddressInfo{Province: "Valparaíso"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Province).To(Equal("Valparaíso"))
		})

		It("should mark as changed when District is updated", func() {
			input := AddressInfo{District: "Las Condes"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.District).To(Equal("Las Condes"))
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
			input := AddressInfo{} // Sin campos
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
			Expect(updated.District).To(Equal("Las Condes"))
			Expect(updated.ZipCode).To(Equal("7550000"))
			// Otros campos deben permanecer igual
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

			addr.Normalize()

			Expect(addr.AddressLine1).To(Equal("avenida providencia 1234"))
			Expect(addr.State).To(Equal("metropolitana"))
			Expect(addr.Province).To(Equal("santiago"))
			Expect(addr.District).To(Equal("providencia"))
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

			// Comprueba que el resultado contiene los campos correctos
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

			// Verifica que solo se incluyen los campos especificados
			Expect(fullAddress).To(Equal("Av Providencia 1234, Providencia, Santiago, Metropolitana, 7500000"))
			Expect(fullAddress).NotTo(ContainSubstring("AV PROVIDENCIA 1234"))
		})

		It("should correctly handle empty fields", func() {
			addr := AddressInfo{
				AddressLine1: "Av Providencia 1234",
				District:     "Providencia",
				// Province vacío
				State:   "Metropolitana",
				ZipCode: "", // Vacío
			}

			fullAddress := addr.FullAddress()

			// Verifica que los campos vacíos se manejen correctamente
			Expect(fullAddress).To(ContainSubstring("Av Providencia 1234"))
			Expect(fullAddress).To(ContainSubstring("Providencia"))
			Expect(fullAddress).To(ContainSubstring("Metropolitana"))
			Expect(fullAddress).NotTo(ContainSubstring(",,")) // No debe haber comas consecutivas
		})

		It("should return just AddressLine1 if other fields are empty", func() {
			addr := AddressInfo{
				AddressLine1: "Av Providencia 1234",
			}

			fullAddress := addr.FullAddress()
			Expect(fullAddress).To(Equal("Av Providencia 1234"))
		})
	})
})

package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contact DocID", func() {
	org := Organization{ID: 1, Country: countries.CL}

	It("should prioritize PrimaryEmail when all fields are present", func() {
		contact := Contact{
			NationalID:   "12345678-9",
			PrimaryEmail: "test@example.com",
			PrimaryPhone: "+56912345678",
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).ToNot(BeEmpty())
		Expect(refID).To(Equal(Hash(org, "test@example.com")))
	})

	It("should use PrimaryPhone when PrimaryEmail is missing", func() {
		contact := Contact{
			NationalID:   "12345678-9",
			PrimaryEmail: "", // Campo vacío
			PrimaryPhone: "+56912345678",
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).ToNot(BeEmpty())
		Expect(refID).To(Equal(Hash(org, "+56912345678")))
	})

	It("should use NationalID when PrimaryEmail and PrimaryPhone are missing", func() {
		contact := Contact{
			NationalID:   "12345678-9",
			PrimaryEmail: "", // Campo vacío
			PrimaryPhone: "", // Campo vacío
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).ToNot(BeEmpty())
		Expect(refID).To(Equal(Hash(org, "12345678-9")))
	})

	It("should generate a UUID when all identifier fields are missing", func() {
		contact := Contact{
			NationalID:   "", // Campo vacío
			PrimaryEmail: "", // Campo vacío
			PrimaryPhone: "", // Campo vacío
			Organization: org,
		}

		refID1 := contact.DocID()
		refID2 := contact.DocID()

		Expect(refID1).ToNot(BeEmpty())
		Expect(refID2).ToNot(BeEmpty())
		Expect(refID1).ToNot(Equal(refID2)) // Deberían ser UUIDs diferentes
	})
})

var _ = Describe("Contact UpdateIfChanged Additional Cases", func() {
	var original Contact
	org := Organization{ID: 1, Country: countries.CL}

	BeforeEach(func() {
		original = Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@correo.com",
			PrimaryPhone: "+56900000000",
			NationalID:   "12345678-9",
			Documents: []Document{
				{Type: "RUT", Value: "12345678-9"},
				{Type: "DRIVER_LICENSE", Value: "A-123456"},
			},
			Organization: org,
			AdditionalContactMethods: []ContactMethod{
				{Type: "work_email", Value: "juan.perez@empresa.com"},
				{Type: "whatsapp", Value: "+56900000001"},
			},
		}
	})

	// Test específico para la condición if newContact.Type != "" && newContact.Type != c.Type
	It("should update ContactMethod Type when it's different and not empty", func() {
		// Crear un contacto con un método que tiene el mismo Type pero diferente Value
		newContact := original
		newContact.AdditionalContactMethods = []ContactMethod{
			{Type: "UPDATED_TYPE", Value: "juan.perez@empresa.com"}, // Type diferente, mismo Value
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())

		// Verificar que se creó un nuevo método con el Type actualizado
		var found bool
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "UPDATED_TYPE" && method.Value == "juan.perez@empresa.com" {
				found = true
				break
			}
		}

		Expect(found).To(BeTrue(), "Debería existir un método con el Type actualizado")

		// Asegurarse que el método original sigue existiendo (porque se busca por Type)
		var originalMethodExists bool
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "work_email" {
				originalMethodExists = true
				break
			}
		}

		Expect(originalMethodExists).To(BeTrue(), "El método original debería mantenerse")
	})

	// Test específico para cuando Type es igual pero Value es diferente
	It("should update ContactMethod Value when Type is the same but Value is different", func() {
		newContact := original
		newContact.AdditionalContactMethods = []ContactMethod{
			{Type: "work_email", Value: "nuevo.email@empresa.com"}, // Mismo Type, Value diferente
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())

		// Verificar que se actualizó el Value pero se mantuvo el mismo Type
		var workEmail ContactMethod
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "work_email" {
				workEmail = method
				break
			}
		}

		Expect(workEmail.Value).To(Equal("nuevo.email@empresa.com"))
	})

	// Test específico para cuando Type y Value son iguales
	It("should not update ContactMethod when both Type and Value are the same", func() {
		newContact := original
		newContact.AdditionalContactMethods = []ContactMethod{
			{Type: "work_email", Value: "juan.perez@empresa.com"}, // Mismo Type, mismo Value
		}

		updated, changed := original.UpdateIfChanged(newContact)

		// No debería haber cambios
		Expect(changed).To(BeFalse())

		// Los métodos deberían ser iguales a los originales
		Expect(updated.AdditionalContactMethods).To(Equal(original.AdditionalContactMethods))
	})

	// Test específico para la condición if newContact.PrimaryEmail != "" && newContact.PrimaryEmail != c.PrimaryEmail
	It("should update PrimaryEmail when it's different and not empty", func() {
		newContact := Contact{
			PrimaryEmail: "nuevo@correo.com", // Email diferente
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.PrimaryEmail).To(Equal("nuevo@correo.com"))

		// Verificar que otros campos se mantienen
		Expect(updated.FullName).To(Equal(original.FullName))
		Expect(updated.PrimaryPhone).To(Equal(original.PrimaryPhone))
		Expect(updated.NationalID).To(Equal(original.NationalID))
	})

	It("should not update PrimaryEmail when it's the same", func() {
		newContact := Contact{
			PrimaryEmail: "juan@correo.com", // Mismo email
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.PrimaryEmail).To(Equal("juan@correo.com"))
	})

	It("should not update PrimaryEmail when new value is empty", func() {
		newContact := Contact{
			PrimaryEmail: "", // Email vacío
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.PrimaryEmail).To(Equal("juan@correo.com")) // Mantiene el valor original
	})

	// Test para verificar que se cumplan ambas condiciones
	It("should only update fields that are different AND not empty", func() {
		newContact := Contact{
			// Campo diferente y no vacío -> debe actualizarse
			FullName: "Nombre Nuevo",

			// Campo diferente pero vacío -> no debe actualizarse
			PrimaryPhone: "",

			// Campo igual -> no debe actualizarse
			PrimaryEmail: "juan@correo.com",

			// Campo vacío -> no debe actualizarse
			NationalID: "",
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue()) // Hubo al menos un cambio

		// Verificar campos que deben y no deben cambiar
		Expect(updated.FullName).To(Equal("Nombre Nuevo"))        // Debe cambiar
		Expect(updated.PrimaryPhone).To(Equal("+56900000000"))    // No debe cambiar
		Expect(updated.PrimaryEmail).To(Equal("juan@correo.com")) // No debe cambiar
		Expect(updated.NationalID).To(Equal("12345678-9"))        // No debe cambiar
	})
})

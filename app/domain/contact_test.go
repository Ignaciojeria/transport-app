package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContactMethod UpdateIfChange", func() {
	It("should update Type if different and not empty", func() {
		method := ContactMethod{Type: "email", Value: "test@example.com"}
		newMethod := ContactMethod{Type: "work_email", Value: ""}

		updated, changed := method.UpdateIfChange(newMethod)

		Expect(changed).To(BeTrue())
		Expect(updated.Type).To(Equal("work_email"))
		Expect(updated.Value).To(Equal("test@example.com")) // No debería cambiar
	})

	It("should not update Type if same", func() {
		method := ContactMethod{Type: "email", Value: "test@example.com"}
		newMethod := ContactMethod{Type: "email", Value: "new@example.com"}

		updated, changed := method.UpdateIfChange(newMethod)

		Expect(changed).To(BeTrue())            // Solo cambió el Value
		Expect(updated.Type).To(Equal("email")) // No debería cambiar
		Expect(updated.Value).To(Equal("new@example.com"))
	})

	It("should not update Type if empty", func() {
		method := ContactMethod{Type: "email", Value: "test@example.com"}
		newMethod := ContactMethod{Type: "", Value: "new@example.com"}

		updated, changed := method.UpdateIfChange(newMethod)

		Expect(changed).To(BeTrue())            // Solo cambió el Value
		Expect(updated.Type).To(Equal("email")) // No debería cambiar porque el nuevo está vacío
		Expect(updated.Value).To(Equal("new@example.com"))
	})
})

var _ = Describe("Contact UpdateIfChanged NationalID", func() {
	var original Contact
	org := Organization{ID: 1, Country: countries.CL}

	BeforeEach(func() {
		original = Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@correo.com",
			PrimaryPhone: "+56900000000",
			NationalID:   "12345678-9",
			Organization: org,
		}
	})

	It("should update NationalID if different and not empty", func() {
		newContact := Contact{
			NationalID: "98765432-1", // Diferente y no vacío
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.NationalID).To(Equal("98765432-1"))
	})

	It("should not update NationalID if same", func() {
		newContact := Contact{
			NationalID: "12345678-9", // Igual al original
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.NationalID).To(Equal("12345678-9"))
	})

	It("should not update NationalID if empty", func() {
		newContact := Contact{
			NationalID: "", // Vacío
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.NationalID).To(Equal("12345678-9")) // Mantiene el valor original
	})
})

var _ = Describe("Contact Documents Continue", func() {
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
		}
	})

	It("should skip completely empty documents with continue", func() {
		newContact := original
		// Incluir varios documentos vacíos entre válidos para probar el continue
		newContact.Documents = []Document{
			{Type: "", Value: ""},                 // Vacío 1 - debe ser ignorado por continue
			{Type: "PASSPORT", Value: "AB123456"}, // Válido - debe agregarse
			{Type: "", Value: ""},                 // Vacío 2 - debe ser ignorado por continue
			{Type: "RUT", Value: "98765432-1"},    // Válido - debe actualizar existente
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		// Debe haber 3 documentos: los 2 originales (uno actualizado) + 1 nuevo
		Expect(updated.Documents).To(HaveLen(3))

		// Verificar que el documento RUT se actualizó correctamente
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}
		Expect(rutDoc.Value).To(Equal("98765432-1"))

		// Verificar que se agregó el documento PASSPORT
		var passportFound bool
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" && doc.Value == "AB123456" {
				passportFound = true
				break
			}
		}
		Expect(passportFound).To(BeTrue())

		// Verificar que no hay documentos vacíos (fueron ignorados por continue)
		emptyCount := 0
		for _, doc := range updated.Documents {
			if doc.Type == "" && doc.Value == "" {
				emptyCount++
			}
		}
		Expect(emptyCount).To(Equal(0))
	})

	It("should handle updates with only empty documents", func() {
		newContact := original
		// Lista con solo documentos vacíos que deben ser ignorados por continue
		newContact.Documents = []Document{
			{Type: "", Value: ""},
			{Type: "", Value: ""},
		}

		updated, changed := original.UpdateIfChanged(newContact)

		// No debe haber cambios porque todos los documentos se ignoran por continue
		Expect(changed).To(BeFalse())
		Expect(updated.Documents).To(Equal(original.Documents))
	})

	It("should correctly process a mix of update, add, and continue operations", func() {
		newContact := original
		// Mezcla completa de operaciones
		newContact.Documents = []Document{
			{Type: "", Value: ""},                       // Debe ignorarse por continue
			{Type: "RUT", Value: "98765432-1"},          // Actualizar existente
			{Type: "PASSPORT", Value: "AB123456"},       // Añadir nuevo
			{Type: "", Value: ""},                       // Debe ignorarse por continue
			{Type: "DRIVER_LICENSE", Value: "A-123456"}, // Sin cambio (igual al original)
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(HaveLen(3)) // 2 originales + 1 nuevo

		// Verificar que se actualizó RUT
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}
		Expect(rutDoc.Value).To(Equal("98765432-1"))

		// Verificar que se agregó PASSPORT
		var passportFound bool
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" {
				passportFound = true
				break
			}
		}
		Expect(passportFound).To(BeTrue())
	})
})

var _ = Describe("Contact AdditionalContactMethods Continue", func() {
	var original Contact
	org := Organization{ID: 1, Country: countries.CL}

	BeforeEach(func() {
		original = Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@correo.com",
			PrimaryPhone: "+56900000000",
			NationalID:   "12345678-9",
			AdditionalContactMethods: []ContactMethod{
				{Type: "work_email", Value: "juan.perez@empresa.com"},
				{Type: "whatsapp", Value: "+56900000001"},
			},
			Organization: org,
		}
	})

	It("should skip completely empty contact methods with continue", func() {
		newContact := original
		// Incluir varios métodos vacíos entre válidos para probar el continue
		newContact.AdditionalContactMethods = []ContactMethod{
			{Type: "", Value: ""},                            // Vacío 1 - debe ser ignorado por continue
			{Type: "telegram", Value: "@juanperez"},          // Válido - debe agregarse
			{Type: "", Value: ""},                            // Vacío 2 - debe ser ignorado por continue
			{Type: "work_email", Value: "nuevo@empresa.com"}, // Válido - debe actualizar existente
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		// Debe haber 3 métodos: los 2 originales (uno actualizado) + 1 nuevo
		Expect(updated.AdditionalContactMethods).To(HaveLen(3))

		// Verificar que el método work_email se actualizó correctamente
		var workEmail ContactMethod
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "work_email" {
				workEmail = method
				break
			}
		}
		Expect(workEmail.Value).To(Equal("nuevo@empresa.com"))

		// Verificar que se agregó el método telegram
		var telegramFound bool
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "telegram" && method.Value == "@juanperez" {
				telegramFound = true
				break
			}
		}
		Expect(telegramFound).To(BeTrue())

		// Verificar que no hay métodos vacíos (fueron ignorados por continue)
		emptyCount := 0
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "" && method.Value == "" {
				emptyCount++
			}
		}
		Expect(emptyCount).To(Equal(0))
	})

	It("should handle updates with only empty contact methods", func() {
		newContact := original
		// Lista con solo métodos vacíos que deben ser ignorados por continue
		newContact.AdditionalContactMethods = []ContactMethod{
			{Type: "", Value: ""},
			{Type: "", Value: ""},
		}

		updated, changed := original.UpdateIfChanged(newContact)

		// No debe haber cambios porque todos los métodos se ignoran por continue
		Expect(changed).To(BeFalse())
		Expect(updated.AdditionalContactMethods).To(Equal(original.AdditionalContactMethods))
	})

	It("should correctly process a mix of update, add, and continue operations for contact methods", func() {
		newContact := original
		// Mezcla completa de operaciones
		newContact.AdditionalContactMethods = []ContactMethod{
			{Type: "", Value: ""},                            // Debe ignorarse por continue
			{Type: "work_email", Value: "nuevo@empresa.com"}, // Actualizar existente
			{Type: "telegram", Value: "@juanperez"},          // Añadir nuevo
			{Type: "", Value: ""},                            // Debe ignorarse por continue
			{Type: "whatsapp", Value: "+56900000001"},        // Sin cambio (igual al original)
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.AdditionalContactMethods).To(HaveLen(3)) // 2 originales + 1 nuevo

		// Verificar que se actualizó work_email
		var workEmail ContactMethod
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "work_email" {
				workEmail = method
				break
			}
		}
		Expect(workEmail.Value).To(Equal("nuevo@empresa.com"))

		// Verificar que se agregó telegram
		var telegramFound bool
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "telegram" {
				telegramFound = true
				break
			}
		}
		Expect(telegramFound).To(BeTrue())
	})
})

var _ = Describe("Contact DocID Function", func() {
	org := Organization{ID: 1, Country: countries.CL}

	It("should use PrimaryEmail as key when present", func() {
		contact := Contact{
			PrimaryEmail: "test@example.com",
			PrimaryPhone: "+56912345678",
			NationalID:   "12345678-9",
			Organization: org,
		}

		docID := contact.DocID()
		Expect(docID).To(Equal(Hash(org, "test@example.com")))
	})

	It("should fallback to PrimaryPhone when PrimaryEmail is missing", func() {
		contact := Contact{
			PrimaryEmail: "", // Vacío
			PrimaryPhone: "+56912345678",
			NationalID:   "12345678-9",
			Organization: org,
		}

		docID := contact.DocID()
		Expect(docID).To(Equal(Hash(org, "+56912345678")))
	})

	It("should fallback to NationalID when PrimaryEmail and PrimaryPhone are missing", func() {
		contact := Contact{
			PrimaryEmail: "", // Vacío
			PrimaryPhone: "", // Vacío
			NationalID:   "12345678-9",
			Organization: org,
		}

		docID := contact.DocID()
		Expect(docID).To(Equal(Hash(org, "12345678-9")))
	})

	It("should generate a UUID when all identifiers are missing", func() {
		contact := Contact{
			PrimaryEmail: "", // Vacío
			PrimaryPhone: "", // Vacío
			NationalID:   "", // Vacío
			Organization: org,
		}

		// Generar dos IDs para verificar que son diferentes (UUIDs)
		docID1 := contact.DocID()
		docID2 := contact.DocID()

		Expect(docID1).ToNot(BeEmpty())
		Expect(docID2).ToNot(BeEmpty())
		Expect(docID1).ToNot(Equal(docID2)) // Deberían ser diferentes UUIDs
	})
})

var _ = Describe("Contact PrimaryEmail UpdateIfChanged", func() {
	var original Contact
	org := Organization{ID: 1, Country: countries.CL}

	BeforeEach(func() {
		original = Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@correo.com",
			PrimaryPhone: "+56900000000",
			NationalID:   "12345678-9",
			Organization: org,
		}
	})

	It("should update PrimaryEmail when different and not empty", func() {
		newContact := Contact{
			PrimaryEmail: "nuevo@correo.com", // Diferente y no vacío
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.PrimaryEmail).To(Equal("nuevo@correo.com"))
	})

	It("should not update PrimaryEmail when same", func() {
		newContact := Contact{
			PrimaryEmail: "juan@correo.com", // Igual al original
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.PrimaryEmail).To(Equal("juan@correo.com"))
	})

	It("should not update PrimaryEmail when empty", func() {
		newContact := Contact{
			PrimaryEmail: "", // Vacío
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.PrimaryEmail).To(Equal("juan@correo.com")) // Mantiene el valor original
	})

	It("should verify both conditions together for PrimaryEmail update", func() {
		// Crear tres contactos para probar las tres combinaciones posibles
		contact1 := Contact{PrimaryEmail: "diferente@correo.com"} // Diferente y no vacío -> debe actualizar
		contact2 := Contact{PrimaryEmail: "juan@correo.com"}      // Igual -> no debe actualizar
		contact3 := Contact{PrimaryEmail: ""}                     // Vacío -> no debe actualizar

		updated1, changed1 := original.UpdateIfChanged(contact1)
		updated2, changed2 := original.UpdateIfChanged(contact2)
		updated3, changed3 := original.UpdateIfChanged(contact3)

		// Verificar resultados
		Expect(changed1).To(BeTrue())
		Expect(updated1.PrimaryEmail).To(Equal("diferente@correo.com"))

		Expect(changed2).To(BeFalse())
		Expect(updated2.PrimaryEmail).To(Equal("juan@correo.com"))

		Expect(changed3).To(BeFalse())
		Expect(updated3.PrimaryEmail).To(Equal("juan@correo.com"))
	})
})

var _ = Describe("Contact UpdateIfChanged (missing fields)", func() {
	var original Contact
	org := Organization{ID: 1}

	BeforeEach(func() {
		original = Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@correo.com",
			PrimaryPhone: "+56900000000",
			NationalID:   "12345678-9",
			Organization: org,
		}
	})

	It("should update FullName if different and not empty", func() {
		newContact := Contact{
			FullName: "Juan Pedro Pérez",
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.FullName).To(Equal("Juan Pedro Pérez"))
	})

	It("should update PrimaryPhone if different and not empty", func() {
		newContact := Contact{
			PrimaryPhone: "+56911111111",
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.PrimaryPhone).To(Equal("+56911111111"))
	})
})

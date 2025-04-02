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

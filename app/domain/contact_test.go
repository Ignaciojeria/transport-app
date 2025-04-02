package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Document UpdateIfChange", func() {
	It("should update Type if different and not empty", func() {
		doc := Document{Type: "RUT", Value: "12345678-9"}
		newDoc := Document{Type: "PASSPORT", Value: ""}

		updated, changed := doc.UpdateIfChange(newDoc)

		Expect(changed).To(BeTrue())
		Expect(updated.Type).To(Equal("PASSPORT"))
		Expect(updated.Value).To(Equal("12345678-9")) // No debería cambiar
	})

	It("should update Value if different and not empty", func() {
		doc := Document{Type: "RUT", Value: "12345678-9"}
		newDoc := Document{Type: "", Value: "98765432-1"}

		updated, changed := doc.UpdateIfChange(newDoc)

		Expect(changed).To(BeTrue())
		Expect(updated.Type).To(Equal("RUT")) // No debería cambiar
		Expect(updated.Value).To(Equal("98765432-1"))
	})

	It("should update both Type and Value if both are different and not empty", func() {
		doc := Document{Type: "RUT", Value: "12345678-9"}
		newDoc := Document{Type: "PASSPORT", Value: "AB123456"}

		updated, changed := doc.UpdateIfChange(newDoc)

		Expect(changed).To(BeTrue())
		Expect(updated.Type).To(Equal("PASSPORT"))
		Expect(updated.Value).To(Equal("AB123456"))
	})

	It("should not update if values are the same", func() {
		doc := Document{Type: "RUT", Value: "12345678-9"}
		newDoc := Document{Type: "RUT", Value: "12345678-9"}

		updated, changed := doc.UpdateIfChange(newDoc)

		Expect(changed).To(BeFalse())
		Expect(updated).To(Equal(doc))
	})

	It("should not update if new values are empty", func() {
		doc := Document{Type: "RUT", Value: "12345678-9"}
		newDoc := Document{Type: "", Value: ""}

		updated, changed := doc.UpdateIfChange(newDoc)

		Expect(changed).To(BeFalse())
		Expect(updated).To(Equal(doc))
	})
})

var _ = Describe("Contact Documents UpdateIfChanged", func() {
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

	It("should update existing document by type", func() {
		newContact := original
		// Actualizar el valor del documento con Type "RUT"
		newContact.Documents = []Document{
			{Type: "RUT", Value: "98765432-1"},
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(HaveLen(2)) // Se mantienen los 2 documentos

		// Buscar el documento actualizado por Type
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}

		Expect(rutDoc.Value).To(Equal("98765432-1"))

		// Verificar que el otro documento no cambió
		var licenseDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "DRIVER_LICENSE" {
				licenseDoc = doc
				break
			}
		}

		Expect(licenseDoc.Value).To(Equal("A-123456"))
	})

	It("should add new document when type doesn't exist", func() {
		newContact := original
		// Agregar un nuevo documento con Type "PASSPORT"
		newContact.Documents = []Document{
			{Type: "PASSPORT", Value: "AB123456"},
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(HaveLen(3)) // Ahora hay 3 documentos

		// Buscar el nuevo documento
		var passportDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" {
				passportDoc = doc
				break
			}
		}

		Expect(passportDoc.Value).To(Equal("AB123456"))
	})

	It("should both update existing and add new documents", func() {
		newContact := original
		// Actualizar uno existente y agregar uno nuevo
		newContact.Documents = []Document{
			{Type: "RUT", Value: "98765432-1"},    // Actualizar
			{Type: "PASSPORT", Value: "AB123456"}, // Agregar
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(HaveLen(3)) // Ahora hay 3 documentos

		// Verificar el que se actualizó
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}
		Expect(rutDoc.Value).To(Equal("98765432-1"))

		// Verificar el nuevo
		var passportDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" {
				passportDoc = doc
				break
			}
		}
		Expect(passportDoc.Value).To(Equal("AB123456"))

		// Verificar que el otro existente no cambió
		var licenseDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "DRIVER_LICENSE" {
				licenseDoc = doc
				break
			}
		}
		Expect(licenseDoc.Value).To(Equal("A-123456"))
	})

	It("should skip completely empty documents", func() {
		newContact := original
		// Incluir un documento completamente vacío entre otros válidos
		newContact.Documents = []Document{
			{Type: "RUT", Value: "98765432-1"},    // Válido
			{Type: "", Value: ""},                 // Completamente vacío, debería ser ignorado
			{Type: "PASSPORT", Value: "AB123456"}, // Válido nuevo
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(HaveLen(3)) // 2 originales + 1 nuevo (el vacío se ignora)

		// Verificar que no hay ningún documento completamente vacío
		emptyDocsCount := 0
		for _, doc := range updated.Documents {
			if doc.Type == "" && doc.Value == "" {
				emptyDocsCount++
			}
		}
		Expect(emptyDocsCount).To(Equal(0))

		// Verificar que se agregó el nuevo documento
		var passportDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" {
				passportDoc = doc
				break
			}
		}
		Expect(passportDoc.Value).To(Equal("AB123456"))
	})

	It("should ignore documents with only Type and empty Value", func() {
		newContact := original
		// Un documento con Type pero Value vacío
		newContact.Documents = []Document{
			{Type: "RUT", Value: ""}, // Solo Type, sin Value
		}

		updated, changed := original.UpdateIfChanged(newContact)

		// No debería cambiar el valor ya que el nuevo Value está vacío
		Expect(changed).To(BeFalse())

		// Verificar que el valor original se mantiene
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}
		Expect(rutDoc.Value).To(Equal("12345678-9"))
	})

	It("should handle multiple document operations in single update", func() {
		newContact := original
		// En un solo update: modificar uno, agregar otro, dejar otro igual, incluir uno vacío
		newContact.Documents = []Document{
			{Type: "RUT", Value: "98765432-1"},          // Modificar
			{Type: "PASSPORT", Value: "AB123456"},       // Agregar
			{Type: "DRIVER_LICENSE", Value: "A-123456"}, // No cambiar
			{Type: "", Value: ""},                       // Vacío - debería ignorarse
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(HaveLen(3)) // 2 originales + 1 nuevo

		// Verificar que RUT se actualizó
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}
		Expect(rutDoc.Value).To(Equal("98765432-1"))

		// Verificar que se agregó PASSPORT
		var passportDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" {
				passportDoc = doc
				break
			}
		}
		Expect(passportDoc.Value).To(Equal("AB123456"))

		// Verificar que LICENSE no cambió
		var licenseDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "DRIVER_LICENSE" {
				licenseDoc = doc
				break
			}
		}
		Expect(licenseDoc.Value).To(Equal("A-123456"))
	})

	It("should handle a complex update with documents, basic fields and contact methods", func() {
		newContact := Contact{
			FullName:     "Nombre Completo Actualizado",
			PrimaryPhone: "+56911222333",
			Documents: []Document{
				{Type: "RUT", Value: "98765432-1"},
				{Type: "PASSPORT", Value: "AB123456"},
			},
			AdditionalContactMethods: []ContactMethod{
				{Type: "work_email", Value: "nuevo.email@empresa.com"},
				{Type: "instagram", Value: "@juanp"},
			},
		}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeTrue())

		// Verificar campos básicos
		Expect(updated.FullName).To(Equal("Nombre Completo Actualizado"))
		Expect(updated.PrimaryPhone).To(Equal("+56911222333"))

		// Verificar documentos
		Expect(updated.Documents).To(HaveLen(3)) // 2 originales + 1 nuevo

		// Verificar que RUT se actualizó
		var rutDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "RUT" {
				rutDoc = doc
				break
			}
		}
		Expect(rutDoc.Value).To(Equal("98765432-1"))

		// Verificar que se agregó PASSPORT
		var passportDoc Document
		for _, doc := range updated.Documents {
			if doc.Type == "PASSPORT" {
				passportDoc = doc
				break
			}
		}
		Expect(passportDoc.Value).To(Equal("AB123456"))

		// Verificar métodos de contacto
		Expect(updated.AdditionalContactMethods).To(HaveLen(3)) // 2 originales + 1 nuevo

		// Verificar que work_email se actualizó
		var workEmail ContactMethod
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "work_email" {
				workEmail = method
				break
			}
		}
		Expect(workEmail.Value).To(Equal("nuevo.email@empresa.com"))

		// Verificar que se agregó instagram
		var instagram ContactMethod
		for _, method := range updated.AdditionalContactMethods {
			if method.Type == "instagram" {
				instagram = method
				break
			}
		}
		Expect(instagram.Value).To(Equal("@juanp"))
	})

	It("should not update documents if no changes are required", func() {
		newContact := original
		// Mismos documentos que ya existen
		newContact.Documents = []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "DRIVER_LICENSE", Value: "A-123456"},
		}

		updated, changed := original.UpdateIfChanged(newContact)

		// Solo los documentos no deberían causar cambios
		Expect(changed).To(BeFalse())
		Expect(updated.Documents).To(Equal(original.Documents))
	})

	It("should handle update with empty documents array", func() {
		newContact := original
		// Array vacío no debería cambiar nada
		newContact.Documents = []Document{}

		updated, changed := original.UpdateIfChanged(newContact)

		Expect(changed).To(BeFalse())
		Expect(updated.Documents).To(Equal(original.Documents))
	})
})

package domain

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Document", func() {
	Describe("UpdateIfChange", func() {
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

		It("should handle updating Type only", func() {
			doc := Document{Type: "RUT", Value: "12345678-9"}
			newDoc := Document{Type: "ID_CARD", Value: "12345678-9"} // Mismo Value

			updated, changed := doc.UpdateIfChange(newDoc)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("ID_CARD"))
			Expect(updated.Value).To(Equal("12345678-9"))
		})

		It("should handle updating Value only", func() {
			doc := Document{Type: "RUT", Value: "12345678-9"}
			newDoc := Document{Type: "RUT", Value: "NEW-VALUE"} // Mismo Type

			updated, changed := doc.UpdateIfChange(newDoc)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("RUT"))
			Expect(updated.Value).To(Equal("NEW-VALUE"))
		})

		It("should handle empty initial document", func() {
			doc := Document{Type: "", Value: ""}
			newDoc := Document{Type: "RUT", Value: "12345678-9"}

			updated, changed := doc.UpdateIfChange(newDoc)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("RUT"))
			Expect(updated.Value).To(Equal("12345678-9"))
		})
	})
})

var _ = Describe("compareDocuments", func() {
	It("should return true for identical documents", func() {
		docs1 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "AB123456"},
		}
		docs2 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "AB123456"},
		}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeTrue())
	})

	It("should return false for documents with different values", func() {
		docs1 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "AB123456"},
		}
		docs2 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "CD789012"}, // Valor diferente
		}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeFalse())
	})

	It("should return false for documents with different types", func() {
		docs1 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "AB123456"},
		}
		docs2 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "DRIVER_LICENSE", Value: "AB123456"}, // Tipo diferente
		}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeFalse())
	})

	It("should return false for documents with different lengths", func() {
		docs1 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "AB123456"},
		}
		docs2 := []Document{
			{Type: "RUT", Value: "12345678-9"},
		}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeFalse())
	})

	It("should return false for documents with same elements in different order", func() {
		docs1 := []Document{
			{Type: "RUT", Value: "12345678-9"},
			{Type: "PASSPORT", Value: "AB123456"},
		}
		docs2 := []Document{
			{Type: "PASSPORT", Value: "AB123456"},
			{Type: "RUT", Value: "12345678-9"},
		}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeFalse())
	})

	It("should return true for empty document arrays", func() {
		docs1 := []Document{}
		docs2 := []Document{}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeTrue())
	})

	It("should return true for nil document arrays", func() {
		var docs1 []Document = nil
		var docs2 []Document = nil

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeTrue())
	})

	It("should return true when comparing nil with empty array", func() {
		var docs1 []Document = nil
		docs2 := []Document{}

		result := compareDocuments(docs1, docs2)
		Expect(result).To(BeTrue()) // Técnicamente, ambos tienen longitud 0
	})
})

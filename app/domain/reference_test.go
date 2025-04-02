package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reference", func() {
	var org1 = Organization{ID: 1, Country: countries.CL}
	var org2 = Organization{ID: 2, Country: countries.AR}

	Describe("DocID", func() {
		It("should generate consistent IDs for the same Organization, Type and Value", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}

			id1 := DocID(org1, ref)
			id2 := DocID(org1, ref)

			Expect(id1).To(Equal(id2))
		})

		It("should generate different IDs for different Organizations", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}

			id1 := DocID(org1, ref)
			id2 := DocID(org2, ref)

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate different IDs for different Types", func() {
			ref1 := Reference{Type: "CODE", Value: "ABC123"}
			ref2 := Reference{Type: "ALT_CODE", Value: "ABC123"}

			id1 := DocID(org1, ref1)
			id2 := DocID(org1, ref2)

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate different IDs for different Values", func() {
			ref1 := Reference{Type: "CODE", Value: "ABC123"}
			ref2 := Reference{Type: "CODE", Value: "XYZ789"}

			id1 := DocID(org1, ref1)
			id2 := DocID(org1, ref2)

			Expect(id1).ToNot(Equal(id2))
		})
	})

	Describe("UpdateIfChange", func() {
		It("should update Type if different and not empty", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "NEW_CODE", Value: ""}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("NEW_CODE"))
			Expect(updated.Value).To(Equal("ABC123")) // No debería cambiar
		})

		It("should update Value if different and not empty", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "", Value: "XYZ789"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("CODE")) // No debería cambiar
			Expect(updated.Value).To(Equal("XYZ789"))
		})

		It("should update both Type and Value if both are different and not empty", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "NEW_CODE", Value: "XYZ789"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("NEW_CODE"))
			Expect(updated.Value).To(Equal("XYZ789"))
		})

		It("should not update if values are the same", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "CODE", Value: "ABC123"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(ref))
		})

		It("should not update if new values are empty", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "", Value: ""}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(ref))
		})

		It("should only update the Type when Value is the same", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "NEW_CODE", Value: "ABC123"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("NEW_CODE"))
			Expect(updated.Value).To(Equal("ABC123"))
		})

		It("should only update the Value when Type is the same", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "CODE", Value: "NEW_VALUE"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("CODE"))
			Expect(updated.Value).To(Equal("NEW_VALUE"))
		})

		It("should handle empty initial values", func() {
			ref := Reference{Type: "", Value: ""}
			newRef := Reference{Type: "CODE", Value: "ABC123"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("CODE"))
			Expect(updated.Value).To(Equal("ABC123"))
		})

		It("should verify that DocID changes when values change", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "NEW_CODE", Value: "XYZ789"}

			originalID := DocID(org1, ref)
			updated, _ := ref.UpdateIfChange(newRef)
			newID := DocID(org1, updated)

			Expect(newID).ToNot(Equal(originalID))
		})
	})
})

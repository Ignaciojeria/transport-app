package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reference", func() {
	var ctx1, ctx2 context.Context

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate consistent IDs for the same context, Type and Value", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}

			id1 := ref.DocID(ctx1, "ORDER-123")
			id2 := ref.DocID(ctx1, "ORDER-123")

			Expect(id1).To(Equal(id2))
		})

		It("should generate different IDs for different contexts", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}

			id1 := ref.DocID(ctx1, "ORDER-123")
			id2 := ref.DocID(ctx2, "ORDER-123")

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate different IDs for different Types", func() {
			ref1 := Reference{Type: "CODE", Value: "ABC123"}
			ref2 := Reference{Type: "ALT_CODE", Value: "ABC123"}

			id1 := ref1.DocID(ctx1, "ORDER-123")
			id2 := ref2.DocID(ctx1, "ORDER-123")

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate different IDs for different Values", func() {
			ref1 := Reference{Type: "CODE", Value: "ABC123"}
			ref2 := Reference{Type: "CODE", Value: "XYZ789"}

			id1 := ref1.DocID(ctx1, "ORDER-123")
			id2 := ref2.DocID(ctx1, "ORDER-123")

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate different IDs for different OrderReferences", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}

			id1 := ref.DocID(ctx1, "ORDER-123")
			id2 := ref.DocID(ctx1, "ORDER-999")

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
			Expect(updated.Value).To(Equal("ABC123"))
		})

		It("should update Value if different and not empty", func() {
			ref := Reference{Type: "CODE", Value: "ABC123"}
			newRef := Reference{Type: "", Value: "XYZ789"}

			updated, changed := ref.UpdateIfChange(newRef)

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("CODE"))
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

			originalID := ref.DocID(ctx1, "ORDER-999")
			updated, _ := ref.UpdateIfChange(newRef)
			newID := updated.DocID(ctx1, "ORDER-999")

			Expect(newID).ToNot(Equal(originalID))
		})
	})
})

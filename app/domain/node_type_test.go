package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodeType", func() {
	var ctx1, ctx2 context.Context

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate ID based on context and Value", func() {
			node1 := NodeType{Value: "foo"}
			node2 := NodeType{Value: "bar"}

			Expect(node1.DocID(ctx1)).To(Equal(DocumentID(HashByTenant(ctx1, "foo"))))
			Expect(node2.DocID(ctx1)).To(Equal(DocumentID(HashByTenant(ctx1, "bar"))))
			Expect(node1.DocID(ctx2)).To(Equal(DocumentID(HashByTenant(ctx2, "foo"))))

			Expect(node1.DocID(ctx1)).ToNot(Equal(node2.DocID(ctx1)))
			Expect(node1.DocID(ctx1)).ToNot(Equal(node1.DocID(ctx2)))
		})
	})

	Describe("UpdateIfChange", func() {
		var base NodeType

		BeforeEach(func() {
			base = NodeType{
				Value: "initial-type",
			}
		})

		It("should update Value if new Value is not empty and different", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value: "new-type",
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Value).To(Equal("new-type"))
		})

		It("should not update Value if new Value is empty", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value: "",
			})

			Expect(changed).To(BeFalse())
			Expect(updated.Value).To(Equal("initial-type"))
		})

		It("should not update Value if new Value is same", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value: "initial-type",
			})

			Expect(changed).To(BeFalse())
			Expect(updated.Value).To(Equal("initial-type"))
		})

		It("should not change anything if value is equal or empty", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value: "",
			})

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(base))
		})
	})
})

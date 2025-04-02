package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodeType", func() {
	var org1 = Organization{ID: 1, Country: countries.CL}
	var org2 = Organization{ID: 2, Country: countries.AR}

	Describe("DocID", func() {
		It("should generate ID based on Organization and Value", func() {
			node1 := NodeType{Organization: org1, Value: "foo"}
			node2 := NodeType{Organization: org1, Value: "bar"}
			node3 := NodeType{Organization: org2, Value: "foo"}

			Expect(node1.DocID()).To(Equal(Hash(org1, "foo")))
			Expect(node2.DocID()).To(Equal(Hash(org1, "bar")))
			Expect(node3.DocID()).To(Equal(Hash(org2, "foo")))

			Expect(node1.DocID()).ToNot(Equal(node2.DocID()))
			Expect(node1.DocID()).ToNot(Equal(node3.DocID()))
		})

		It("should use empty string as fallback if Value is empty", func() {
			node := NodeType{Organization: org1, Value: ""}
			Expect(node.DocID()).To(Equal(Hash(org1, "")))
		})
	})

	Describe("UpdateIfChange", func() {
		var base NodeType

		BeforeEach(func() {
			base = NodeType{
				Organization: org1,
				Value:        "initial-type",
			}
		})

		It("should update Value if new Value is not empty and different", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value:        "new-type",
				Organization: org1, // mismo para que no afecte
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Value).To(Equal("new-type"))
			Expect(updated.Organization).To(Equal(org1))
		})

		It("should not update Value if new Value is empty", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value:        "",
				Organization: org1,
			})

			Expect(changed).To(BeFalse())
			Expect(updated.Value).To(Equal("initial-type"))
		})

		It("should not update Value if new Value is same", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value:        "initial-type",
				Organization: org1,
			})

			Expect(changed).To(BeFalse())
			Expect(updated.Value).To(Equal("initial-type"))
		})

		It("should update Organization if different", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value:        "initial-type",
				Organization: org2,
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Organization).To(Equal(org2))
		})

		It("should update both fields if both are different and valid", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value:        "new-type",
				Organization: org2,
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Value).To(Equal("new-type"))
			Expect(updated.Organization).To(Equal(org2))
		})

		It("should not change anything if both fields are equal or empty", func() {
			updated, changed := base.UpdateIfChange(NodeType{
				Value:        "",
				Organization: org1,
			})

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(base))
		})
	})
})

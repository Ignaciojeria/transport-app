package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("NodeInfo", func() {
	var ctx1, ctx2 context.Context

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate unique ID based on context and ReferenceID", func() {
			node1 := NodeInfo{
				ReferenceID: "node-1",
			}
			node2 := NodeInfo{
				ReferenceID: "node-2",
			}

			// Mismo nodo con diferentes contextos
			Expect(node1.DocID(ctx1)).To(Equal(HashByTenant(ctx1, "node-1")))
			Expect(node1.DocID(ctx1)).ToNot(Equal(node2.DocID(ctx1)))
			Expect(node1.DocID(ctx1)).ToNot(Equal(node1.DocID(ctx2)))
		})
	})

	Describe("UpdateIfChanged", func() {
		var baseNode NodeInfo

		BeforeEach(func() {
			baseNode = NodeInfo{
				ReferenceID: "node-test",
				Name:        "Nodo Original",
				NodeType:    NodeType{Value: "WAREHOUSE"},
				AddressInfo: AddressInfo{
					AddressLine1: "Av Providencia 1234",
					District:     "Providencia",
					Province:     "Santiago",
					State:        "Metropolitana",
					Coordinates: Coordinates{
						Point:  orb.Point{-70.6506, -33.4372},
						Source: "geocoding",
						Confidence: CoordinatesConfidence{
							Level:   0.9,
							Message: "High confidence",
							Reason:  "Exact match",
						},
					},
				},
				References: []Reference{
					{Type: "CODE", Value: "REF001"},
					{Type: "ALT_CODE", Value: "ALT001"},
				},
			}
		})

		It("should update basic fields", func() {
			newNode := baseNode
			newNode.Name = "Nodo Actualizado"
			newNode.ReferenceID = "node-updated"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.Name).To(Equal("Nodo Actualizado"))
			Expect(updated.ReferenceID).To(Equal(ReferenceID("node-updated")))
		})

		It("should not update node type from domain logic", func() {
			newNode := baseNode
			newNode.NodeType.Value = "STORE"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeFalse())
			Expect(updated.NodeType.Value).To(Equal("WAREHOUSE"))
		})

		It("should not update nested AddressInfo from domain logic", func() {
			newNode := baseNode
			newNode.AddressInfo.District = "Las Condes"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeFalse())
			Expect(updated.AddressInfo.District.String()).To(Equal("Providencia"))
		})

		It("should update references", func() {
			newNode := baseNode
			newNode.References = []Reference{
				{Type: "CODE", Value: "REF001-UPDATED"},
				{Type: "CUSTOM", Value: "CUST001"},
			}

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.References).To(HaveLen(3))
			Expect(updated.References).To(ContainElement(Reference{Type: "CODE", Value: "REF001-UPDATED"}))
			Expect(updated.References).To(ContainElement(Reference{Type: "ALT_CODE", Value: "ALT001"}))
			Expect(updated.References).To(ContainElement(Reference{Type: "CUSTOM", Value: "CUST001"}))
		})

		It("should ignore empty references", func() {
			newNode := baseNode
			newNode.References = []Reference{{Type: "", Value: ""}}

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeFalse())
			Expect(updated.References).To(Equal(baseNode.References))
		})

		It("should not update if no fields changed", func() {
			updated, changed := baseNode.UpdateIfChanged(baseNode)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(baseNode))
		})

		It("should ignore empty values", func() {
			newNode := NodeInfo{
				Name:        "",
				ReferenceID: "",
				References:  []Reference{},
			}

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(baseNode))
		})
	})
})

package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("NodeInfo", func() {
	var org1 = Organization{ID: 1, Country: countries.CL}
	var org2 = Organization{ID: 2, Country: countries.AR}

	Describe("DocID", func() {
		It("should generate unique ID based on Organization and ReferenceID", func() {
			node1 := NodeInfo{
				Organization: org1,
				ReferenceID:  "node-1",
			}
			node2 := NodeInfo{
				Organization: org1,
				ReferenceID:  "node-2",
			}
			node3 := NodeInfo{
				Organization: org2,
				ReferenceID:  "node-1",
			}

			Expect(node1.DocID()).To(Equal(Hash(org1, "node-1")))
			Expect(node1.DocID()).ToNot(Equal(node2.DocID()))
			Expect(node1.DocID()).ToNot(Equal(node3.DocID()))
		})

		It("should return empty string if ReferenceID is empty", func() {
			node := NodeInfo{
				Organization: org1,
				ReferenceID:  "",
			}
			Expect(node.DocID()).To(Equal(""))
		})
	})

	Describe("UpdateIfChanged", func() {
		var baseNode NodeInfo

		BeforeEach(func() {
			baseNode = NodeInfo{
				ReferenceID:  "node-test",
				Name:         "Nodo Original",
				Organization: org1,
				NodeType:     NodeType{Value: "WAREHOUSE"},
				AddressInfo: AddressInfo{
					AddressLine1: "Av Providencia 1234",
					District:     "Providencia",
					Province:     "Santiago",
					State:        "Metropolitana",
					Location:     orb.Point{-70.6506, -33.4372},
				},
				Contact: Contact{
					FullName:     "Juan Pérez",
					PrimaryEmail: "juan@example.com",
					PrimaryPhone: "+56912345678",
				},
				AddressLine2: "Dpto 1402",
				AddressLine3: "Torre Norte",
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
			Expect(updated.AddressInfo.District).To(Equal("Providencia"))
		})

		It("should not update nested Contact from domain logic", func() {
			newNode := baseNode
			newNode.Contact.FullName = "María López"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeFalse())
			Expect(updated.Contact.FullName).To(Equal("Juan Pérez"))
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
			newNode.References = []Reference{
				{Type: "", Value: ""},
			}

			updated, changed := baseNode.UpdateIfChanged(newNode)
			Expect(changed).To(BeFalse())
			Expect(updated.References).To(Equal(baseNode.References))
		})

		It("should update AddressLine2 and AddressLine3", func() {
			newNode := baseNode
			newNode.AddressLine2 = "Nueva dirección 2"
			newNode.AddressLine3 = "Nueva dirección 3"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine2).To(Equal("Nueva dirección 2"))
			Expect(updated.AddressLine3).To(Equal("Nueva dirección 3"))
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

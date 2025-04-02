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

			// Mismo ReferenceID pero diferente organización -> diferente DocID
			Expect(node1.DocID()).To(Equal(Hash(org1, "node-1")))
			Expect(node1.DocID()).ToNot(Equal(node2.DocID()))
			Expect(node1.DocID()).ToNot(Equal(node3.DocID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		var baseNode NodeInfo

		BeforeEach(func() {
			// Configurar un nodo base para cada test
			baseNode = NodeInfo{
				ReferenceID:  "node-test",
				Name:         "Nodo Original",
				Organization: org1,
				NodeType: NodeType{
					ID:    1,
					Value: "WAREHOUSE",
				},
				AddressInfo: AddressInfo{
					AddressLine1: "Av Providencia 1234",
					District:     "Providencia",
					Province:     "Santiago",
					State:        "Metropolitana",
					Location:     orb.Point{-70.6506, -33.4372}, // [lon, lat]
				},
				Contact: Contact{
					FullName:     "Juan Pérez",
					PrimaryEmail: "juan@example.com",
					PrimaryPhone: "+56912345678",
				},
				AddressLine2: "Dpto 1402",
				AddressLine3: "Torre Norte",
				References:   []Reference{},
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

		It("should update node type", func() {
			newNode := baseNode
			newNode.NodeType.Value = "STORE"
			newNode.NodeType.ID = 2

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.NodeType.Value).To(Equal("STORE"))
			Expect(updated.NodeType.ID).To(Equal(int64(2)))
		})

		It("should update references", func() {
			newNode := baseNode
			newNode.References = []Reference{
				{Type: "CODE", Value: "REF001"},
				{Type: "ALT_CODE", Value: "ALT001"},
			}

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.References).To(HaveLen(2))
			Expect(updated.References[0].Value).To(Equal("REF001"))
			Expect(updated.References[1].Value).To(Equal("ALT001"))
		})

		It("should update AddressLine2 and AddressLine3", func() {
			newNode := baseNode
			newNode.AddressLine2 = "Dpto 1403" // Cambiado
			newNode.AddressLine3 = "Torre Sur" // Cambiado

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine2).To(Equal("Dpto 1403"))
			Expect(updated.AddressLine3).To(Equal("Torre Sur"))
		})

		It("should update nested AddressInfo", func() {
			newNode := baseNode
			newNode.AddressInfo.AddressLine1 = "Av Las Condes 5678"
			newNode.AddressInfo.District = "Las Condes"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.AddressInfo.AddressLine1).To(Equal("Av Las Condes 5678"))
			Expect(updated.AddressInfo.District).To(Equal("Las Condes"))
		})

		It("should update nested Contact", func() {
			newNode := baseNode
			newNode.Contact.FullName = "María López"
			newNode.Contact.PrimaryPhone = "+56987654321"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.Contact.FullName).To(Equal("María López"))
			Expect(updated.Contact.PrimaryPhone).To(Equal("+56987654321"))
		})

		It("should not update if no fields changed", func() {
			updated, changed := baseNode.UpdateIfChanged(baseNode)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(baseNode))
		})

		It("should ignore empty string and zero values in updates", func() {
			newNode := NodeInfo{
				// Campos vacíos que no deberían actualizar
				Name:         "",
				ReferenceID:  "",
				AddressLine2: "",
				AddressLine3: "",
				NodeType: NodeType{
					ID:    0,
					Value: "",
				},
			}

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(baseNode))
		})

		It("should handle multiple field updates together", func() {
			newNode := baseNode
			newNode.Name = "Nodo Multi-Update"
			newNode.AddressLine2 = "Oficina 303"
			newNode.Contact.FullName = "Carlos Ramírez"
			newNode.AddressInfo.District = "Vitacura"

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.Name).To(Equal("Nodo Multi-Update"))
			Expect(updated.AddressLine2).To(Equal("Oficina 303"))
			Expect(updated.Contact.FullName).To(Equal("Carlos Ramírez"))
			Expect(updated.AddressInfo.District).To(Equal("Vitacura"))
		})

		It("should detect changes in references even if length is the same", func() {
			baseNode.References = []Reference{
				{Type: "CODE", Value: "OLD001"},
			}

			newNode := baseNode
			newNode.References = []Reference{
				{Type: "CODE", Value: "NEW001"},
			}

			updated, changed := baseNode.UpdateIfChanged(newNode)

			Expect(changed).To(BeTrue())
			Expect(updated.References[0].Value).To(Equal("NEW001"))
		})
	})
})

package tidbrepository

import (
	"context"

	"github.com/paulmach/orb"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertNodeInfo", func() {
	var (
		ctx1, ctx2                 context.Context
		tenant1                    domain.Tenant
		nodeType1, nodeType2       domain.NodeType
		contact1, contact2         domain.Contact
		addressInfo1, addressInfo2 domain.AddressInfo
	)

	BeforeEach(func() {
		var err error
		// Create two tenants for testing
		tenant1, ctx1, err = CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		_, ctx2, err = CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		// Setup test data
		nodeType1 = domain.NodeType{Value: "type1"}
		nodeType2 = domain.NodeType{Value: "type2"}

		contact1 = domain.Contact{
			FullName:     "Contact 1",
			PrimaryEmail: "contact1@test.com",
		}
		contact2 = domain.Contact{
			FullName:     "Contact 2",
			PrimaryEmail: "contact2@test.com",
		}

		addressInfo1 = domain.AddressInfo{
			AddressLine1: "Address 1",
			Contact:      contact1,
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6001, -33.4500},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}
		addressInfo2 = domain.AddressInfo{
			AddressLine1: "Address 2",
			Contact:      contact2,
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6002, -33.4501},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}
	})

	It("should insert node info if not exists", func() {
		nodeInfo := domain.NodeInfo{
			Name:        "Warehouse Alpha",
			ReferenceID: "WH-001", // Ensure we have a valid ReferenceID
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // This now contains contact1

			References: []domain.Reference{
				{Type: "ERP", Value: "WH001"},
			},
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err := upsert(ctx1, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", nodeInfo.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo.Name).To(Equal("Warehouse Alpha"))
		Expect(dbNodeInfo.NodeTypeDoc).To(Equal(nodeType1.DocID(ctx1).String()))
		Expect(dbNodeInfo.AddressInfoDoc).To(Equal(addressInfo1.DocID(ctx1).String()))

		Expect(dbNodeInfo.TenantID.String()).To(Equal(tenant1.ID.String()))

		// Verify that the JSON references were saved correctly
		references := dbNodeInfo.NodeReferences.Map()
		Expect(len(references)).To(Equal(1))
		Expect(references[0].Type).To(Equal("ERP"))
		Expect(references[0].Value).To(Equal("WH001"))
	})

	It("should update node info if fields are different", func() {
		original := domain.NodeInfo{
			Name:        "Original Name",
			ReferenceID: "REF-001",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo

		}

		upsert := NewUpsertNodeInfo(connection, nil)

		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:        "Modified Name",
			ReferenceID: "REF-001", // Keep same ReferenceID to ensure same DocID
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo

		}

		err = upsert(ctx1, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", original.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo.Name).To(Equal("Modified Name"))

		Expect(dbNodeInfo.TenantID.String()).To(Equal(tenant1.ID.String()))
	})

	It("should not update if no fields changed", func() {
		nodeInfo := domain.NodeInfo{
			Name:        "Unchanged Node",
			ReferenceID: "REF-002",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo

		}

		upsert := NewUpsertNodeInfo(connection, nil)

		err := upsert(ctx1, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		// Get the initial update timestamp
		var initialNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", nodeInfo.DocID(ctx1)).
			First(&initialNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		initialUpdatedAt := initialNodeInfo.UpdatedAt

		// Execute again without changes
		err = upsert(ctx1, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		// Verify timestamp hasn't changed, indicating no update occurred
		var finalNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", nodeInfo.DocID(ctx1)).
			First(&finalNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(finalNodeInfo.UpdatedAt).To(Equal(initialUpdatedAt))
		Expect(finalNodeInfo.TenantID.String()).To(Equal(tenant1.ID.String()))
	})

	It("should allow same node info for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		nodeInfo1 := domain.NodeInfo{
			Name: "Multi Org Node",
		}

		nodeInfo2 := domain.NodeInfo{
			Name: "Multi Org Node",
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err = upsert(ctx1, nodeInfo1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, nodeInfo2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := nodeInfo1.DocID(ctx1)
		docID2 := nodeInfo2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each node info belongs to its respective tenant using DocID
		var dbNodeInfo1, dbNodeInfo2 table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", docID1).
			First(&dbNodeInfo1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = connection.DB.WithContext(ctx2).
			Table("node_infos").
			Where("document_id = ?", docID2).
			First(&dbNodeInfo2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should update when related node type changes", func() {
		original := domain.NodeInfo{
			Name:        "Node With Type Change",
			ReferenceID: "REF-TYPE",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:        "Node With Type Change",
			ReferenceID: "REF-TYPE",
			NodeType:    nodeType2,    // Changed node type
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		err = upsert(ctx1, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", original.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo.NodeTypeDoc).To(Equal(nodeType2.DocID(ctx1).String()))
		Expect(dbNodeInfo.TenantID.String()).To(Equal(tenant1.ID.String()))
	})

	It("should update when related address info with different contact changes", func() {
		original := domain.NodeInfo{
			Name:        "Node With Address Change",
			ReferenceID: "REF-ADDRESS",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contains contact1
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:        "Node With Address Change",
			ReferenceID: "REF-ADDRESS",
			NodeType:    nodeType1,
			AddressInfo: addressInfo2, // Contains contact2
		}

		err = upsert(ctx1, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", original.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo.AddressInfoDoc).To(Equal(addressInfo2.DocID(ctx1).String()))
		Expect(dbNodeInfo.TenantID.String()).To(Equal(tenant1.ID.String()))
	})

	It("should update addressLine2 and addressLine3 when they change", func() {
		original := domain.NodeInfo{
			Name:        "Node With Address Lines",
			ReferenceID: "REF-ADDRLINES",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:        "Node With Address Lines",
			ReferenceID: "REF-ADDRLINES",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		err = upsert(ctx1, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", original.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo.TenantID.String()).To(Equal(tenant1.ID.String()))
	})

	It("should update references when they change", func() {
		original := domain.NodeInfo{
			Name:        "Node With References",
			ReferenceID: "REF-REFS",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
			References: []domain.Reference{
				{Type: "ERP", Value: "WH001"},
			},
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:        "Node With References",
			ReferenceID: "REF-REFS",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
			References: []domain.Reference{
				{Type: "ERP", Value: "WH001-UPDATED"},
				{Type: "SAP", Value: "1000002"},
			},
		}

		err = upsert(ctx1, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", original.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())

		// Verify the references were updated
		references := dbNodeInfo.NodeReferences.Map()
		Expect(len(references)).To(Equal(2))

		// Check if the first reference was updated
		erp := findReferenceByType(references, "ERP")
		Expect(erp).ToNot(BeNil())
		Expect(erp.Value).To(Equal("WH001-UPDATED"))

		// Check if the new reference was added
		sap := findReferenceByType(references, "SAP")
		Expect(sap).ToNot(BeNil())
		Expect(sap.Value).To(Equal("1000002"))
	})

	It("should generate consistent DocIDs for same input", func() {
		// Create two identical NodeInfo objects
		nodeInfo1 := domain.NodeInfo{
			Name:        "Hash Test Node",
			ReferenceID: "REF-HASH",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		nodeInfo2 := domain.NodeInfo{
			Name:        "Hash Test Node",
			ReferenceID: "REF-HASH",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		// DocIDs should be identical for identical input in the same context
		Expect(nodeInfo1.DocID(ctx1)).To(Equal(nodeInfo2.DocID(ctx1)))

		// DocIDs should be different in different contexts (different organizations)
		Expect(nodeInfo1.DocID(ctx1)).ToNot(Equal(nodeInfo1.DocID(ctx2)))
	})

	It("should fail if database has no node_infos table", func() {
		nodeInfo := domain.NodeInfo{
			Name:        "Expected Error",
			ReferenceID: "REF-ERROR",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		upsert := NewUpsertNodeInfo(noTablesContainerConnection, nil)
		err := upsert(ctx1, nodeInfo)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("node_infos"))
	})

	It("should handle headers according to the three scenarios", func() {
		// Scenario 1: Create new node with empty headers
		node := domain.NodeInfo{
			Name:        "Node With Headers",
			ReferenceID: "REF-HEADERS",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1,
			Headers:     domain.Headers{}, // Empty headers
		}

		upsert := NewUpsertNodeInfo(connection, nil)
		err := upsert(ctx1, node)
		Expect(err).ToNot(HaveOccurred())

		// Verify headers were persisted as empty
		var dbNode table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", node.DocID(ctx1)).
			First(&dbNode).Error
		Expect(err).ToNot(HaveOccurred())
		emptyHeaders := domain.Headers{}.DocID(ctx1).String()
		Expect(dbNode.NodeInfoHeadersDoc).To(Equal(emptyHeaders))
		Expect(dbNode.TenantID.String()).To(Equal(tenant1.ID.String()))

		// Scenario 2: Update existing node with empty headers
		updatedNode := node
		updatedNode.Headers = domain.Headers{} // Still empty headers

		err = upsert(ctx1, updatedNode)
		Expect(err).ToNot(HaveOccurred())

		// Verify headers were not updated (still empty)
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", node.DocID(ctx1)).
			First(&dbNode).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNode.NodeInfoHeadersDoc).To(Equal(emptyHeaders))
		Expect(dbNode.TenantID.String()).To(Equal(tenant1.ID.String()))

		// Scenario 3: Update existing node with non-empty headers
		updatedNode.Headers = domain.Headers{
			Commerce: "commerce1",
			Consumer: "consumer1",
			Channel:  "channel1",
		}

		err = upsert(ctx1, updatedNode)
		Expect(err).ToNot(HaveOccurred())

		// Verify headers were updated
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", node.DocID(ctx1)).
			First(&dbNode).Error
		Expect(err).ToNot(HaveOccurred())
		newHeaders := updatedNode.Headers.DocID(ctx1).String()
		Expect(dbNode.NodeInfoHeadersDoc).To(Equal(newHeaders))
		Expect(dbNode.TenantID.String()).To(Equal(tenant1.ID.String()))

		// Scenario 4: Update existing node with empty headers
		updatedNodeFinal := node
		updatedNodeFinal.Headers = domain.Headers{} // Still empty headers

		err = upsert(ctx1, updatedNodeFinal)
		Expect(err).ToNot(HaveOccurred())

		// Verify headers are empty
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", node.DocID(ctx1)).
			First(&dbNode).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNode.NodeInfoHeadersDoc).To(Equal(emptyHeaders))
		Expect(dbNode.TenantID.String()).To(Equal(tenant1.ID.String()))
	})

})

// Helper function to find a reference by its type
func findReferenceByType(references []domain.Reference, refType string) *domain.Reference {
	for _, ref := range references {
		if ref.Type == refType {
			return &ref
		}
	}
	return nil
}

package tidbrepository

import (
	"context"
	"strconv"

	"github.com/paulmach/orb"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertNodeInfo", func() {
	var (
		ctx1, ctx2                 context.Context
		nodeType1, nodeType2       domain.NodeType
		contact1, contact2         domain.Contact
		addressInfo1, addressInfo2 domain.AddressInfo
	)

	// Helper function to create context with organization
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		// Create contexts with different organizations
		ctx1 = createOrgContext(organization1)
		ctx2 = createOrgContext(organization2)

		// Create node types
		nodeType1 = domain.NodeType{
			Value: "warehouse",
		}

		nodeType2 = domain.NodeType{
			Value: "distribution-center",
		}

		// Create contacts
		contact1 = domain.Contact{
			FullName:     "John Smith",
			PrimaryEmail: "john@example.com",
			PrimaryPhone: "123456789",
			Documents: []domain.Document{
				{Type: "ID", Value: "123456789"},
			},
			AdditionalContactMethods: []domain.ContactMethod{
				{Type: "whatsapp", Value: "+56987654321"},
			},
		}

		contact2 = domain.Contact{
			FullName:     "Jane Doe",
			PrimaryEmail: "jane@example.com",
			PrimaryPhone: "987654321",
		}

		// Create address info with contacts embedded
		addressInfo1 = domain.AddressInfo{
			Contact:      contact1,
			State:        "State1",
			Province:     "Province1",
			District:     "District1",
			AddressLine1: "123 Main St",
			Location:     orb.Point{-70.123, -33.456}, // longitude, latitude
			ZipCode:      "12345",
			TimeZone:     "America/Santiago",
		}

		addressInfo2 = domain.AddressInfo{
			Contact:      contact2,
			State:        "State2",
			Province:     "Province2",
			District:     "District2",
			AddressLine1: "456 Second Ave",
			Location:     orb.Point{-71.234, -34.567},
			ZipCode:      "54321",
			TimeZone:     "America/Santiago",
		}
	})

	It("should insert node info if not exists", func() {
		nodeInfo := domain.NodeInfo{
			Name:         "Warehouse Alpha",
			ReferenceID:  "WH-001", // Ensure we have a valid ReferenceID
			NodeType:     nodeType1,
			AddressInfo:  addressInfo1, // This now contains contact1
			AddressLine2: "Floor 3",
			References: []domain.Reference{
				{Type: "ERP", Value: "WH001"},
			},
		}

		upsert := NewUpsertNodeInfo(connection)
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
		Expect(dbNodeInfo.AddressLine2).To(Equal("Floor 3"))

		// Verify that the JSON references were saved correctly
		references := dbNodeInfo.NodeReferences.Map()
		Expect(len(references)).To(Equal(1))
		Expect(references[0].Type).To(Equal("ERP"))
		Expect(references[0].Value).To(Equal("WH001"))
	})

	It("should update node info if fields are different", func() {
		original := domain.NodeInfo{
			Name:         "Original Name",
			ReferenceID:  "REF-001",
			NodeType:     nodeType1,
			AddressInfo:  addressInfo1, // Contact is inside addressInfo
			AddressLine2: "Original Floor",
		}

		upsert := NewUpsertNodeInfo(connection)

		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:         "Modified Name",
			ReferenceID:  "REF-001", // Keep same ReferenceID to ensure same DocID
			NodeType:     nodeType1,
			AddressInfo:  addressInfo1, // Contact is inside addressInfo
			AddressLine2: "Modified Floor",
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
		Expect(dbNodeInfo.AddressLine2).To(Equal("Modified Floor"))
	})

	It("should not update if no fields changed", func() {
		nodeInfo := domain.NodeInfo{
			Name:         "Unchanged Node",
			ReferenceID:  "REF-002",
			NodeType:     nodeType1,
			AddressInfo:  addressInfo1, // Contact is inside addressInfo
			AddressLine2: "Same Floor",
		}

		upsert := NewUpsertNodeInfo(connection)

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
	})

	It("should allow same node info for different organizations", func() {
		nodeInfo1 := domain.NodeInfo{
			Name:        "Multi Org Node",
			ReferenceID: "REF-ORG-1",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		// Same node info structure for organization2
		nodeInfo2 := domain.NodeInfo{
			Name:        "Multi Org Node",
			ReferenceID: "REF-ORG-2",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		upsert := NewUpsertNodeInfo(connection)

		err := upsert(ctx1, nodeInfo1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, nodeInfo2)
		Expect(err).ToNot(HaveOccurred())

		// Verify we have two distinct records
		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("node_infos").
			Where("name = ?", nodeInfo1.Name).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(nodeInfo1.DocID(ctx1)).ToNot(Equal(nodeInfo2.DocID(ctx2)))
	})

	It("should update when related node type changes", func() {
		original := domain.NodeInfo{
			Name:        "Node With Type Change",
			ReferenceID: "REF-TYPE",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contact is inside addressInfo
		}

		upsert := NewUpsertNodeInfo(connection)
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
	})

	It("should update when related address info with different contact changes", func() {
		original := domain.NodeInfo{
			Name:        "Node With Address Change",
			ReferenceID: "REF-ADDRESS",
			NodeType:    nodeType1,
			AddressInfo: addressInfo1, // Contains contact1
		}

		upsert := NewUpsertNodeInfo(connection)
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
	})

	It("should update addressLine2 and addressLine3 when they change", func() {
		original := domain.NodeInfo{
			Name:         "Node With Address Lines",
			ReferenceID:  "REF-ADDRLINES",
			NodeType:     nodeType1,
			AddressInfo:  addressInfo1, // Contact is inside addressInfo
			AddressLine2: "Original Line 2",
		}

		upsert := NewUpsertNodeInfo(connection)
		err := upsert(ctx1, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeInfo{
			Name:         "Node With Address Lines",
			ReferenceID:  "REF-ADDRLINES",
			NodeType:     nodeType1,
			AddressInfo:  addressInfo1, // Contact is inside addressInfo
			AddressLine2: "Updated Line 2",
		}

		err = upsert(ctx1, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeInfo table.NodeInfo
		err = connection.DB.WithContext(ctx1).
			Table("node_infos").
			Where("document_id = ?", original.DocID(ctx1)).
			First(&dbNodeInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeInfo.AddressLine2).To(Equal("Updated Line 2"))
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

		upsert := NewUpsertNodeInfo(connection)
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

		upsert := NewUpsertNodeInfo(noTablesContainerConnection)
		err := upsert(ctx1, nodeInfo)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("node_infos"))
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

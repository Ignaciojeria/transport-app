package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/projection/nodes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("FindNodesProjectionResult", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should return empty list when no nodes exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		findNodes := NewFindNodesProjectionResult(
			conn,
			nodes.NewProjection())
		results, hasMore, err := findNodes(ctx, domain.NodesFilter{})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
	})

	It("should find nodes by reference ID", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create test node
		nodeType := domain.NodeType{Value: "WAREHOUSE"}
		addressInfo := domain.AddressInfo{
			AddressLine1: "Test Address",
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
		nodeInfo := domain.NodeInfo{
			Name:        "Test Node",
			ReferenceID: "TEST-001",
			NodeType:    nodeType,
			AddressInfo: addressInfo,
		}

		// Insert test data
		upsert := NewUpsertNodeInfo(conn)
		err = upsert(ctx, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		upsertNodeType := NewUpsertNodeType(conn)
		err = upsertNodeType(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		upsertAddressInfo := NewUpsertAddressInfo(conn, nil)
		err = upsertAddressInfo(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		projection := nodes.NewProjection()

		// Test finding by reference ID
		findNodes := NewFindNodesProjectionResult(
			conn,
			projection)
		results, hasMore, err := findNodes(ctx, domain.NodesFilter{
			ReferenceIds: []string{"TEST-001"},
			RequestedFields: map[string]any{
				projection.AddressInfo().String():  true,
				projection.AddressLine1().String(): true,
				projection.NodeInfo().String():     true,
				projection.NodeName().String():     true,
				projection.NodeType().String():     true,
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).ToNot(BeNil())
		Expect(hasMore).To(BeFalse())
		Expect(results).To(HaveLen(1))
		Expect(results[0].NodeName).To(Equal("Test Node"))
		Expect(results[0].NodeType).To(Equal("WAREHOUSE"))
		Expect(results[0].AddressLine1).To(Equal("Test Address"))
	})

	It("should find nodes by name", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create test node
		nodeType := domain.NodeType{Value: "STORE"}
		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: "State",
				AdminAreaLevel2: "Central",
				AdminAreaLevel3: "Downtown",
				TimeZone:        "America/Santiago",
			},
			AddressLine1: "Store Address",
			AddressLine2: "Suite 100",
			ZipCode:      "12345",
		}
		nodeInfo := domain.NodeInfo{
			Name:        "Test Store",
			ReferenceID: "STORE-001",
			NodeType:    nodeType,
			AddressInfo: addressInfo,
		}

		// Insert test data
		upsert := NewUpsertNodeInfo(conn)
		err = upsert(ctx, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		upsertNodeType := NewUpsertNodeType(conn)
		err = upsertNodeType(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		upsertAddressInfo := NewUpsertAddressInfo(conn, nil)
		err = upsertAddressInfo(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Test finding by name
		name := "Store"
		projection := nodes.NewProjection()
		findNodes := NewFindNodesProjectionResult(
			conn,
			projection)
		results, hasMore, err := findNodes(ctx, domain.NodesFilter{
			Name: &name,
			RequestedFields: map[string]any{
				projection.NodeName().String(): true,
				projection.NodeType().String(): true,
				projection.NodeInfo().String(): true,
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).ToNot(BeNil())
		Expect(hasMore).To(BeFalse())
		Expect(results).To(HaveLen(1))
		Expect(results[0].NodeName).To(Equal("Test Store"))
		Expect(results[0].NodeType).To(Equal("STORE"))
	})

	It("should find nodes by node type", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create test node
		nodeType := domain.NodeType{Value: "DISTRIBUTION"}
		addressInfo := domain.AddressInfo{
			AddressLine1: "Distribution Center",
		}
		nodeInfo := domain.NodeInfo{
			Name:        "Distribution Center 1",
			ReferenceID: "DC-001",
			NodeType:    nodeType,
			AddressInfo: addressInfo,
		}

		// Insert test data
		upsert := NewUpsertNodeInfo(conn)
		err = upsert(ctx, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		upsertNodeType := NewUpsertNodeType(conn)
		err = upsertNodeType(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		upsertAddressInfo := NewUpsertAddressInfo(conn, nil)
		err = upsertAddressInfo(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Test finding by node type
		projection := nodes.NewProjection()
		findNodes := NewFindNodesProjectionResult(
			conn,
			projection)
		results, hasMore, err := findNodes(ctx, domain.NodesFilter{
			NodeType: &domain.NodeTypeFilter{Value: "DISTRIBUTION"},
			RequestedFields: map[string]any{
				projection.NodeName().String(): true,
				projection.NodeType().String(): true,
				projection.NodeInfo().String(): true,
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).ToNot(BeNil())
		Expect(hasMore).To(BeFalse())
		Expect(results).To(HaveLen(1))
		Expect(results[0].NodeName).To(Equal("Distribution Center 1"))
		Expect(results[0].NodeType).To(Equal("DISTRIBUTION"))
	})

	It("should find nodes by references", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create test node with references
		nodeType := domain.NodeType{Value: "WAREHOUSE"}
		addressInfo := domain.AddressInfo{
			AddressLine1: "Warehouse Address",
		}
		nodeInfo := domain.NodeInfo{
			Name:        "Warehouse 1",
			ReferenceID: "WH-001",
			NodeType:    nodeType,
			AddressInfo: addressInfo,
			References: []domain.Reference{
				{Type: "ERP", Value: "WH001"},
				{Type: "WMS", Value: "DC001"},
			},
		}

		// Insert test data
		upsert := NewUpsertNodeInfo(conn)
		err = upsert(ctx, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		upsertNodeType := NewUpsertNodeType(conn)
		err = upsertNodeType(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		upsertAddressInfo := NewUpsertAddressInfo(conn, nil)
		err = upsertAddressInfo(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		upsertNodeReferences := NewUpsertNodeReferences(conn)
		err = upsertNodeReferences(ctx, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		// Test finding by references
		projection := nodes.NewProjection()
		findNodes := NewFindNodesProjectionResult(
			conn,
			projection)
		results, hasMore, err := findNodes(ctx, domain.NodesFilter{
			References: []domain.ReferenceFilter{
				{Type: "ERP", Value: "WH001"},
			},
			RequestedFields: map[string]any{
				projection.NodeName().String():       true,
				projection.NodeType().String():       true,
				projection.NodeInfo().String():       true,
				projection.NodeReferences().String(): true,
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).ToNot(BeNil())
		Expect(hasMore).To(BeFalse())
		Expect(results).To(HaveLen(1))
		Expect(results[0].NodeName).To(Equal("Warehouse 1"))
		Expect(results[0].NodeType).To(Equal("WAREHOUSE"))
	})

	It("should return requested fields only", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create test node
		nodeType := domain.NodeType{Value: "STORE"}
		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: "State",
				AdminAreaLevel2: "Central",
				AdminAreaLevel3: "Downtown",
				TimeZone:        "America/Santiago",
			},
			AddressLine1: "Store Address",
			AddressLine2: "Suite 100",
			ZipCode:      "12345",
		}
		nodeInfo := domain.NodeInfo{
			Name:        "Test Store",
			ReferenceID: "STORE-001",
			NodeType:    nodeType,
			AddressInfo: addressInfo,
		}

		// Insert test data
		upsert := NewUpsertNodeInfo(conn)
		err = upsert(ctx, nodeInfo)
		Expect(err).ToNot(HaveOccurred())

		upsertNodeType := NewUpsertNodeType(conn)
		err = upsertNodeType(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		upsertAddressInfo := NewUpsertAddressInfo(conn, nil)
		err = upsertAddressInfo(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Test finding with specific fields
		projection := nodes.NewProjection()
		findNodes := NewFindNodesProjectionResult(
			conn,
			projection)
		results, hasMore, err := findNodes(ctx, domain.NodesFilter{
			ReferenceIds: []string{"STORE-001"},
			RequestedFields: map[string]any{
				projection.NodeName().String():     true,
				projection.NodeType().String():     true,
				projection.AddressInfo().String():  true,
				projection.AddressLine1().String(): true,
				projection.NodeInfo().String():     true,
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).ToNot(BeNil())
		Expect(hasMore).To(BeFalse())
		Expect(results).To(HaveLen(1))
		Expect(results[0].NodeName).To(Equal("Test Store"))
		Expect(results[0].NodeType).To(Equal("STORE"))
		Expect(results[0].AddressLine1).To(Equal("Store Address"))
		Expect(results[0].AddressLine2).To(BeEmpty()) // Should not be included
		Expect(results[0].District).To(BeEmpty())     // Should not be included
	})
})

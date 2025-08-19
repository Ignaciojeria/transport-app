package tidbrepository

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertRoute", func() {
	var (
		ctx          context.Context
		tenant       domain.Tenant
		conn         database.ConnectionFactory
		upsert       UpsertRoute
		testRoute    domain.Route
		testContract map[string]interface{}
	)

	BeforeEach(func() {
		var err error
		// Create a new tenant for testing
		tenant, ctx, err = CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		conn = connection
		upsert = NewUpsertRoute(conn, nil)

		// Setup test contract data
		testContract = map[string]interface{}{
			"routeId": "route-001",
			"status":  "active",
			"metadata": map[string]interface{}{
				"priority": "high",
				"notes":    "Test route",
			},
		}

		// Setup test data
		testRoute = domain.Route{
			ReferenceID: "route-001",
			Origin: domain.NodeInfo{
				ReferenceID: "origin-001",
				Name:        "Origin Location",
			},
			Destination: domain.NodeInfo{
				ReferenceID: "dest-001",
				Name:        "Destination Location",
			},
			Vehicle: domain.Vehicle{
				Plate: "ABC123",
			},
			Orders: []domain.Order{
				{ReferenceID: "order-001"},
			},
		}
	})

	It("should insert a new route when it doesn't exist", func() {
		err := upsert(ctx, testRoute, testContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		var savedRoute table.Route
		err = conn.DB.Where("document_id = ?", testRoute.DocID(ctx)).First(&savedRoute).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(savedRoute.ReferenceID).To(Equal(testRoute.ReferenceID))
		Expect(savedRoute.DocumentID).To(Equal(testRoute.DocID(ctx).String()))
		Expect(savedRoute.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verify Raw field contains the contract data
		Expect(savedRoute.Raw).ToNot(BeNil())
		var savedContract map[string]interface{}
		rawBytes, err := json.Marshal(savedRoute.Raw)
		Expect(err).ToNot(HaveOccurred())
		err = json.Unmarshal(rawBytes, &savedContract)
		Expect(err).ToNot(HaveOccurred())
		Expect(savedContract["routeId"]).To(Equal("route-001"))
		Expect(savedContract["status"]).To(Equal("active"))
	})

	It("should update an existing route when it exists", func() {
		// First insert
		err := upsert(ctx, testRoute, testContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		// Modify the route and contract
		updatedRoute := testRoute
		updatedRoute.Origin.Name = "Updated Origin"
		updatedRoute.Vehicle.Plate = "XYZ789"

		updatedContract := map[string]interface{}{
			"routeId": "route-001",
			"status":  "updated",
			"metadata": map[string]interface{}{
				"priority": "low",
				"notes":    "Updated test route",
			},
		}

		// Update
		err = upsert(ctx, updatedRoute, updatedContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		// Verify update
		var savedRoute table.Route
		err = conn.DB.Where("document_id = ?", updatedRoute.DocID(ctx)).First(&savedRoute).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(savedRoute.ReferenceID).To(Equal(updatedRoute.ReferenceID))
		Expect(savedRoute.DocumentID).To(Equal(updatedRoute.DocID(ctx).String()))
		Expect(savedRoute.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verify Raw field contains the updated contract data
		Expect(savedRoute.Raw).ToNot(BeNil())
		var savedContract map[string]interface{}
		rawBytes, err := json.Marshal(savedRoute.Raw)
		Expect(err).ToNot(HaveOccurred())
		err = json.Unmarshal(rawBytes, &savedContract)
		Expect(err).ToNot(HaveOccurred())
		Expect(savedContract["status"]).To(Equal("updated"))
		Expect(savedContract["metadata"].(map[string]interface{})["priority"]).To(Equal("low"))
	})

	It("should maintain the same ID when updating an existing route", func() {
		// First insert
		err := upsert(ctx, testRoute, testContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		var firstRoute table.Route
		err = conn.DB.Where("document_id = ?", testRoute.DocID(ctx)).First(&firstRoute).Error
		Expect(err).ToNot(HaveOccurred())
		firstID := firstRoute.ID

		// Update
		updatedRoute := testRoute
		updatedRoute.Origin.Name = "Updated Origin"
		err = upsert(ctx, updatedRoute, testContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		// Verify same ID
		var updatedRouteRecord table.Route
		err = conn.DB.Where("document_id = ?", updatedRoute.DocID(ctx)).First(&updatedRouteRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(updatedRouteRecord.ID).To(Equal(firstID))
		Expect(updatedRouteRecord.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should handle different routes with different DocIDs", func() {
		// Insert first route
		err := upsert(ctx, testRoute, testContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		// Create and insert second route
		secondRoute := domain.Route{
			ReferenceID: "route-002",
			Origin: domain.NodeInfo{
				ReferenceID: "origin-002",
				Name:        "Second Origin",
			},
		}
		secondContract := map[string]interface{}{
			"routeId": "route-002",
			"status":  "pending",
		}
		err = upsert(ctx, secondRoute, secondContract, "plan-doc-001")
		Expect(err).ToNot(HaveOccurred())

		// Verify both routes exist
		var routes []table.Route
		err = conn.DB.Find(&routes).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(routes).To(HaveLen(2))
		Expect(routes[0].TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(routes[1].TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should handle database errors gracefully", func() {
		// Force a database error by inserting invalid data
		invalidRoute := domain.Route{
			ReferenceID: "route-001",
			Origin: domain.NodeInfo{
				ReferenceID: "origin-001",
				Name:        "Origin Location",
			},
			Destination: domain.NodeInfo{
				ReferenceID: "dest-001",
				Name:        "Destination Location",
			},
			Vehicle: domain.Vehicle{
				Plate: "ABC123",
			},
			Orders: []domain.Order{
				{ReferenceID: "order-001"},
			},
		}

		// Create a context with an invalid tenant ID to force a foreign key error
		invalidCtx := context.Background()
		invalidCtx = context.WithValue(invalidCtx, "tenant_id", "invalid-tenant-id")
		err := upsert(invalidCtx, invalidRoute, testContract, "plan-doc-001")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("violates foreign key constraint"))
	})
})

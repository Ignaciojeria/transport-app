package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertVehicleCategory", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert vehicle category if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		vehicleCategory := domain.VehicleCategory{
			Type:                "VAN",
			MaxPackagesQuantity: 100,
		}

		upsert := NewUpsertVehicleCategory(conn, nil)
		err = upsert(ctx, vehicleCategory)
		Expect(err).ToNot(HaveOccurred())

		// Verify using document_id
		var dbVehicleCategory table.VehicleCategory
		err = conn.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ? AND tenant_id = ?", vehicleCategory.DocID(ctx), tenant.ID).
			First(&dbVehicleCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicleCategory.Type).To(Equal("VAN"))
		Expect(dbVehicleCategory.MaxPackagesQuantity).To(Equal(100))
		Expect(dbVehicleCategory.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update vehicle category if fields are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.VehicleCategory{
			Type:                "TRUCK",
			MaxPackagesQuantity: 50,
		}

		upsert := NewUpsertVehicleCategory(conn, nil)
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.VehicleCategory{
			Type:                "TRUCK",
			MaxPackagesQuantity: 75,
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verify using document_id
		var dbVehicleCategory table.VehicleCategory
		err = conn.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ? AND tenant_id = ?", modified.DocID(ctx), tenant.ID).
			First(&dbVehicleCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicleCategory.Type).To(Equal("TRUCK"))
		Expect(dbVehicleCategory.MaxPackagesQuantity).To(Equal(75))
		Expect(dbVehicleCategory.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update vehicle category if fields are the same", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.VehicleCategory{
			Type:                "VAN",
			MaxPackagesQuantity: 100,
		}

		upsert := NewUpsertVehicleCategory(conn, nil)
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Try to update with same values
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Verify using document_id
		var dbVehicleCategory table.VehicleCategory
		err = conn.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ? AND tenant_id = ?", original.DocID(ctx), tenant.ID).
			First(&dbVehicleCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicleCategory.Type).To(Equal("VAN"))
		Expect(dbVehicleCategory.MaxPackagesQuantity).To(Equal(100))
		Expect(dbVehicleCategory.TenantID.String()).To(Equal(tenant.ID.String()))
	})
})

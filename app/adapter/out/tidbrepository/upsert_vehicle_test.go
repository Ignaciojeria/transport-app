package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertVehicle", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertVehicle
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertVehicle(conn, nil)
	})

	It("should insert vehicle if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		vehicle := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "TRUCK",
				MaxPackagesQuantity: 50,
			},
		}

		err = upsert(ctx, vehicle)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := vehicle.DocID(ctx)

		var dbVehicle table.Vehicle
		err = conn.DB.WithContext(ctx).
			Table("vehicles").
			Where("document_id = ?", docID).
			First(&dbVehicle).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle.Plate).To(Equal("ABC123"))
		Expect(dbVehicle.CertificateDate).To(Equal("2024-01-01"))
		Expect(dbVehicle.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update vehicle if it exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "TRUCK",
				MaxPackagesQuantity: 50,
			},
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-02-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 100,
			},
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := modified.DocID(ctx)

		var dbVehicle table.Vehicle
		err = conn.DB.WithContext(ctx).
			Table("vehicles").
			Where("document_id = ?", docID).
			First(&dbVehicle).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle.Plate).To(Equal("ABC123"))
		Expect(dbVehicle.CertificateDate).To(Equal("2024-02-01"))
		Expect(dbVehicle.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same vehicle for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		vehicle1 := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "TRUCK",
				MaxPackagesQuantity: 50,
			},
		}

		vehicle2 := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "TRUCK",
				MaxPackagesQuantity: 50,
			},
		}

		err = upsert(ctx1, vehicle1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, vehicle2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := vehicle1.DocID(ctx1)
		docID2 := vehicle2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each vehicle belongs to its respective tenant using DocID
		var dbVehicle1, dbVehicle2 table.Vehicle
		err = conn.DB.WithContext(ctx1).
			Table("vehicles").
			Where("document_id = ?", docID1).
			First(&dbVehicle1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("vehicles").
			Where("document_id = ?", docID2).
			First(&dbVehicle2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no vehicles table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		vehicle := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
		}

		upsert := NewUpsertVehicle(noTablesContainerConnection, nil)
		err = upsert(ctx, vehicle)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("vehicles"))
	})
})

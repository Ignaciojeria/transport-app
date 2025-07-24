package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertCarrier", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertCarrier
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertCarrier(conn, nil)
	})

	It("should insert carrier if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		carrier := domain.Carrier{
			Name:       "Transportes XYZ",
			NationalID: "12345678-9",
		}

		err = upsert(ctx, carrier)
		Expect(err).ToNot(HaveOccurred())

		var dbCarrier table.Carrier
		err = conn.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", carrier.DocID(ctx)).
			First(&dbCarrier).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier.Name).To(Equal("Transportes XYZ"))
		Expect(dbCarrier.NationalID).To(Equal("12345678-9"))
		Expect(dbCarrier.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update carrier if fields are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Carrier{
			Name:       "Nombre Original",
			NationalID: "11111111-1",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Carrier{
			Name:       "Nombre Modificado",
			NationalID: "11111111-1",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbCarrier table.Carrier
		err = conn.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbCarrier).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier.Name).To(Equal("Nombre Modificado"))
		Expect(dbCarrier.NationalID).To(Equal("11111111-1"))
		Expect(dbCarrier.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update if no fields changed", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Carrier{
			Name:       "Sin Cambios",
			NationalID: "22222222-2",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Try to update with same values
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		var dbCarrier table.Carrier
		err = conn.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbCarrier).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier.Name).To(Equal("Sin Cambios"))
		Expect(dbCarrier.NationalID).To(Equal("22222222-2"))
		Expect(dbCarrier.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same carrier for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		carrier1 := domain.Carrier{
			Name:       "Multi Org Carrier",
			NationalID: "33333333-3",
		}

		carrier2 := domain.Carrier{
			Name:       "Multi Org Carrier",
			NationalID: "33333333-3",
		}

		err = upsert(ctx1, carrier1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, carrier2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(context.Background()).
			Table("carriers").
			Where("document_id IN (?, ?)", carrier1.DocID(ctx1), carrier2.DocID(ctx2)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs and belong to different tenants
		Expect(carrier1.DocID(ctx1)).ToNot(Equal(carrier2.DocID(ctx2)))

		// Verify each carrier belongs to its respective tenant
		var dbCarrier1, dbCarrier2 table.Carrier
		err = conn.DB.WithContext(ctx1).
			Table("carriers").
			Where("document_id = ?", carrier1.DocID(ctx1)).
			First(&dbCarrier1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("carriers").
			Where("document_id = ?", carrier2.DocID(ctx2)).
			First(&dbCarrier2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbCarrier2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no carriers table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		carrier := domain.Carrier{
			Name:       "Error Case",
			NationalID: "44444444-4",
		}

		upsert := NewUpsertCarrier(noTablesContainerConnection, nil)
		err = upsert(ctx, carrier)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("carriers"))
	})
})

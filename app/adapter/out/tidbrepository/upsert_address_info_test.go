package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("UpsertAddressInfo", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertAddressInfo
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertAddressInfo(conn)
	})

	It("should insert addressInfo and its related entities if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		state := domain.State("Metropolitana")
		province := domain.Province("Santiago")
		district := domain.District("Providencia")

		addressInfo := domain.AddressInfo{
			AddressLine1: "Av Providencia 1234",
			State:        state,
			Province:     province,
			District:     district,
			ZipCode:      "7500000",
			Location:     orb.Point{-70.6506, -33.4372}, // [lon, lat]
		}

		err = upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Verify state was created
		var dbState table.State
		err = conn.DB.WithContext(ctx).
			Table("states").
			Where("document_id = ?", state.DocID(ctx).String()).
			First(&dbState).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbState.Name).To(Equal("Metropolitana"))
		Expect(dbState.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verify province was created
		var dbProvince table.Province
		err = conn.DB.WithContext(ctx).
			Table("provinces").
			Where("document_id = ?", province.DocID(ctx).String()).
			First(&dbProvince).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbProvince.Name).To(Equal("Santiago"))
		Expect(dbProvince.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verify district was created
		var dbDistrict table.District
		err = conn.DB.WithContext(ctx).
			Table("districts").
			Where("document_id = ?", district.DocID(ctx).String()).
			First(&dbDistrict).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDistrict.Name).To(Equal("Providencia"))
		Expect(dbDistrict.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verify addressInfo was created with correct references
		var dbAddressInfo table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", addressInfo.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Av Providencia 1234"))
		Expect(dbAddressInfo.StateDoc).To(Equal(state.DocID(ctx).String()))
		Expect(dbAddressInfo.ProvinceDoc).To(Equal(province.DocID(ctx).String()))
		Expect(dbAddressInfo.DistrictDoc).To(Equal(district.DocID(ctx).String()))
		Expect(dbAddressInfo.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should create new address and reuse existing entities when fields change", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.AddressInfo{
			AddressLine1: "Dirección Original",
			District:     "Las Condes",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7550000",
			Location:     orb.Point{-70.5768, -33.4002}, // [lon, lat]
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original DocID
		originalDocID := original.DocID(ctx)

		// Get the IDs of the related entities
		stateID := original.State.DocID(ctx).String()
		provinceID := original.Province.DocID(ctx).String()
		districtID := original.District.DocID(ctx).String()

		modified := domain.AddressInfo{
			AddressLine1: "Dirección Modificada",
			District:     "Las Condes",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7560000", // Cambiado
			Location:     orb.Point{-70.5768, -33.4002},
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Get the new DocID
		modifiedDocID := modified.DocID(ctx)

		// Verify that DocIDs are different
		Expect(modifiedDocID).ToNot(Equal(originalDocID))

		// Verify both addresses exist in the database
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id IN ? AND tenant_id = ?", []string{string(originalDocID), string(modifiedDocID)}, tenant.ID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify the original address still exists
		var originalAddress table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", originalDocID).
			First(&originalAddress).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalAddress.AddressLine1).To(Equal("Dirección Original"))
		Expect(originalAddress.ZipCode).To(Equal("7550000"))

		// Verify the new address exists
		var modifiedAddress table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", modifiedDocID).
			First(&modifiedAddress).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(modifiedAddress.AddressLine1).To(Equal("Dirección Modificada"))
		Expect(modifiedAddress.ZipCode).To(Equal("7560000"))

		// Verify that the related entities were reused (same IDs)
		Expect(modifiedAddress.StateDoc).To(Equal(stateID))
		Expect(modifiedAddress.ProvinceDoc).To(Equal(provinceID))
		Expect(modifiedAddress.DistrictDoc).To(Equal(districtID))
	})

	It("should allow same address info for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo1 := domain.AddressInfo{
			AddressLine1: "Test Street",
			State:        domain.State("Test State"),
			Province:     domain.Province("Test Province"),
			District:     domain.District("Test District"),
		}

		addressInfo2 := domain.AddressInfo{
			AddressLine1: "Test Street",
			State:        domain.State("Test State"),
			Province:     domain.Province("Test Province"),
			District:     domain.District("Test District"),
		}

		err = upsert(ctx1, addressInfo1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, addressInfo2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := addressInfo1.DocID(ctx1)
		docID2 := addressInfo2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each address info belongs to its respective tenant using DocID
		var dbAddressInfo1, dbAddressInfo2 table.AddressInfo
		err = conn.DB.WithContext(ctx1).
			Table("address_infos").
			Where("document_id = ?", docID1).
			First(&dbAddressInfo1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("address_infos").
			Where("document_id = ?", docID2).
			First(&dbAddressInfo2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo2.TenantID.String()).To(Equal(tenant2.ID.String()))

		// Verify that related entities were created for each tenant
		state1ID := addressInfo1.State.DocID(ctx1).String()
		state2ID := addressInfo2.State.DocID(ctx2).String()
		Expect(state1ID).ToNot(Equal(state2ID))

		province1ID := addressInfo1.Province.DocID(ctx1).String()
		province2ID := addressInfo2.Province.DocID(ctx2).String()
		Expect(province1ID).ToNot(Equal(province2ID))

		district1ID := addressInfo1.District.DocID(ctx1).String()
		district2ID := addressInfo2.District.DocID(ctx2).String()
		Expect(district1ID).ToNot(Equal(district2ID))
	})

	It("should update location coordinates correctly", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.AddressInfo{
			AddressLine1: "Coordenadas Test",
			District:     "Ñuñoa",
			Province:     "Santiago",
			State:        "Metropolitana",
			Location:     orb.Point{-70.5975, -33.4566}, // [lon, lat]
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := original.DocID(ctx)

		// Actualizar con nuevas coordenadas
		modified := original
		modified.Location = orb.Point{-70.6000, -33.4600} // Nuevas coordenadas

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verify the coordinates were updated correctly
		var dbAddressInfo table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", docID).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())

		// Comparar coordenadas con una pequeña tolerancia
		Expect(dbAddressInfo.Longitude).To(BeNumerically("~", -70.6000, 0.0001))
		Expect(dbAddressInfo.Latitude).To(BeNumerically("~", -33.4600, 0.0001))
		Expect(dbAddressInfo.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should correctly handle address without coordinates", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			AddressLine1: "Sin Coordenadas",
			District:     "La Florida",
			Province:     "Santiago",
			State:        "Metropolitana",
		}

		err = upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := addressInfo.DocID(ctx)

		// Verify the record was inserted correctly
		var dbAddressInfo table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", docID).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Sin Coordenadas"))
		Expect(dbAddressInfo.TenantID.String()).To(Equal(tenant.ID.String()))
		// Las coordenadas deberían ser cero o un valor por defecto
	})

	It("should fail if database has no address_infos table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			AddressLine1: "Error Esperado",
			District:     "Providencia",
			Province:     "Santiago",
			State:        "Metropolitana",
		}

		upsert := NewUpsertAddressInfo(noTablesContainerConnection)
		err = upsert(ctx, addressInfo)

		Expect(err).To(HaveOccurred())
	})
})

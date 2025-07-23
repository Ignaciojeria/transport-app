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
		upsert = NewUpsertAddressInfo(conn, nil)
	})

	It("should insert addressInfo and its related entities if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		politicalArea := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				Code:            politicalArea.Code,
				AdminAreaLevel1: politicalArea.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea.AdminAreaLevel3,
				TimeZone:        politicalArea.TimeZone,
			},
			AddressLine1: "Av Providencia 1234",
			ZipCode:      "7500000",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6506, -33.4372}, // [lon, lat]
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}

		err = upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertPoliticalArea(conn)(ctx, politicalArea)
		Expect(err).ToNot(HaveOccurred())

		// Verify political area was created
		var dbPoliticalArea table.PoliticalArea
		err = conn.DB.WithContext(ctx).
			Table("political_areas").
			Where("document_id = ?", politicalArea.DocID(ctx)).
			First(&dbPoliticalArea).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPoliticalArea.AdminAreaLevel1).To(Equal("region metropolitana de santiago"))
		Expect(dbPoliticalArea.AdminAreaLevel2).To(Equal("santiago"))
		Expect(dbPoliticalArea.AdminAreaLevel3).To(Equal("la florida"))
		Expect(dbPoliticalArea.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verify addressInfo was created with correct references
		var dbAddressInfo table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", addressInfo.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Av Providencia 1234"))
		politicalAreaDoc := politicalArea.DocID(ctx).String()
		Expect(dbAddressInfo.PoliticalAreaDoc).To(Equal(politicalAreaDoc))
		Expect(dbAddressInfo.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should create new address and reuse existing political area when fields change", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		politicalArea := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		original := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: politicalArea.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea.AdminAreaLevel3,
				TimeZone:        politicalArea.TimeZone,
			},
			AddressLine1: "Dirección Original",
			ZipCode:      "7550000",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.5768, -33.4002},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original DocID
		originalDocID := original.DocID(ctx)

		// Get the ID of the political area
		politicalAreaID := politicalArea.DocID(ctx).String()

		modified := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				Code:            politicalArea.Code,
				AdminAreaLevel1: politicalArea.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea.AdminAreaLevel3,
				TimeZone:        politicalArea.TimeZone,
			},
			AddressLine1: "Dirección Modificada",
			ZipCode:      "7560000", // Cambiado
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.5800, -33.4100},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
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

		// Verify that the political area was reused (same ID)
		Expect(modifiedAddress.PoliticalAreaDoc).To(Equal(politicalAreaID))
	})

	It("should allow same address info for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		politicalArea1 := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		politicalArea2 := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		addressInfo1 := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: politicalArea1.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea1.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea1.AdminAreaLevel3,
				TimeZone:        politicalArea1.TimeZone,
			},
			AddressLine1: "Test Street",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6506, -33.4372},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}

		addressInfo2 := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: politicalArea2.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea2.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea2.AdminAreaLevel3,
				TimeZone:        politicalArea2.TimeZone,
			},
			AddressLine1: "Test Street",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6506, -33.4372},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
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

		// Verify that political areas were created for each tenant
		politicalArea1ID := politicalArea1.DocID(ctx1).String()
		politicalArea2ID := politicalArea2.DocID(ctx2).String()
		Expect(politicalArea1ID).ToNot(Equal(politicalArea2ID))
	})

	It("should update location coordinates correctly", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		politicalArea := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		original := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: politicalArea.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea.AdminAreaLevel3,
				TimeZone:        politicalArea.TimeZone,
			},
			AddressLine1: "Dirección Original",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.5975, -33.4566}, // [lon, lat]
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Actualizar con nuevas coordenadas
		modified := original
		modified.Coordinates.Point = orb.Point{-70.5800, -33.4100}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que la dirección se actualizó correctamente
		var dbAddressInfo table.AddressInfo
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())

		// Comparar coordenadas con una pequeña tolerancia
		Expect(dbAddressInfo.Longitude).To(BeNumerically("~", -70.5800, 0.0001))
		Expect(dbAddressInfo.Latitude).To(BeNumerically("~", -33.4100, 0.0001))
		Expect(dbAddressInfo.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should correctly handle address without coordinates", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		politicalArea := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: politicalArea.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea.AdminAreaLevel3,
				TimeZone:        politicalArea.TimeZone,
			},
			AddressLine1: "Sin Coordenadas",
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

		politicalArea := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: politicalArea.AdminAreaLevel1,
				AdminAreaLevel2: politicalArea.AdminAreaLevel2,
				AdminAreaLevel3: politicalArea.AdminAreaLevel3,
				TimeZone:        politicalArea.TimeZone,
			},
			AddressLine1: "Error Esperado",
		}

		upsert := NewUpsertAddressInfo(noTablesContainerConnection, nil)
		err = upsert(ctx, addressInfo)

		Expect(err).To(HaveOccurred())
	})
})

package tidbrepository

import (
	"context"

	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("GetStoredCoordinates", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})
	/*
		It("should return coordinates if address exists", func() {
			// Create a new tenant for this test
			tenant, ctx, err := CreateTestTenant(context.Background(), conn)
			Expect(err).ToNot(HaveOccurred())

			state := domain.State("Ñuñoa")
			province := domain.Province("Santiago")
			district := domain.District("Ñuñoa")

			addressInfo := domain.AddressInfo{
				PoliticalArea: domain.PoliticalArea{
					State:    string(state),
					Province: string(province),
					District: string(district),
					TimeZone: "America/Santiago",
				},
				AddressLine1: "Coords Existentes",
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

			// Guardar la dirección (esto también guardará state, province y district)
			upsert := NewUpsertAddressInfo(conn)
			err = upsert(ctx, addressInfo)
			Expect(err).ToNot(HaveOccurred())

			// Verificar que la dirección pertenece al tenant correcto
			var dbAddressInfo table.AddressInfo
			err = conn.DB.WithContext(ctx).
				Table("address_infos").
				Where("document_id = ?", addressInfo.DocID(ctx)).
				First(&dbAddressInfo).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(dbAddressInfo.TenantID.String()).To(Equal(tenant.ID.String()))

			// Verificar que las entidades relacionadas se crearon correctamente
			var dbState table.State
			err = conn.DB.WithContext(ctx).
				Table("states").
				Where("document_id = ?", state.DocID(ctx).String()).
				First(&dbState).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(dbState.TenantID.String()).To(Equal(tenant.ID.String()))

			var dbProvince table.Province
			err = conn.DB.WithContext(ctx).
				Table("provinces").
				Where("document_id = ?", province.DocID(ctx).String()).
				First(&dbProvince).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(dbProvince.TenantID.String()).To(Equal(tenant.ID.String()))

			var dbDistrict table.District
			err = conn.DB.WithContext(ctx).
				Table("districts").
				Where("document_id = ?", district.DocID(ctx).String()).
				First(&dbDistrict).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(dbDistrict.TenantID.String()).To(Equal(tenant.ID.String()))

			getCoords := NewGetStoredCoordinates(conn)
			point, err := getCoords(ctx, addressInfo)
			Expect(err).ToNot(HaveOccurred())
			Expect(point.Lon()).To(BeNumerically("~", -70.6001, 0.0001))
			Expect(point.Lat()).To(BeNumerically("~", -33.4500, 0.0001))
		})*/

	It("should return empty point if address does not exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				State:    "Metropolitana",
				Province: "Santiago",
				District: "Peñalolén",
				TimeZone: "America/Santiago",
			},
			AddressLine1: "No Existe",
		}

		getCoords := NewGetStoredCoordinates(conn)
		point, err := getCoords(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())
		Expect(point).To(Equal(orb.Point{}))
	})

	It("should return error if address_infos table does not exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				State:    "Metropolitana",
				Province: "Santiago",
				District: "Independencia",
				TimeZone: "America/Santiago",
			},
			AddressLine1: "Error Forzado",
		}

		getCoords := NewGetStoredCoordinates(noTablesContainerConnection)
		_, err = getCoords(ctx, addressInfo)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("address_infos"))
	})
})

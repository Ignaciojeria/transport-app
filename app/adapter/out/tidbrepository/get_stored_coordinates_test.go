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

var _ = Describe("GetStoredCoordinates", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should return coordinates if address exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear y guardar State, Province y District primero
		state := domain.State("Ñuñoa")
		province := domain.Province("Santiago")
		district := domain.District("Ñuñoa")

		upsertState := NewUpsertState(conn)
		err = upsertState(ctx, state)
		Expect(err).ToNot(HaveOccurred())

		upsertProvince := NewUpsertProvince(conn)
		err = upsertProvince(ctx, province)
		Expect(err).ToNot(HaveOccurred())

		upsertDistrict := NewUpsertDistrict(conn)
		err = upsertDistrict(ctx, district)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			AddressLine1: "Coords Existentes",
			State:        state,
			Province:     province,
			District:     district,
			Location:     orb.Point{-70.6001, -33.4500},
		}

		// Guardar la dirección
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

		getCoords := NewGetStoredCoordinates(conn)
		point, err := getCoords(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())
		Expect(point.Lon()).To(BeNumerically("~", -70.6001, 0.0001))
		Expect(point.Lat()).To(BeNumerically("~", -33.4500, 0.0001))
	})

	It("should return empty point if address does not exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			AddressLine1: "No Existe",
			District:     "Peñalolén",
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
			AddressLine1: "Error Forzado",
			District:     "Independencia",
		}

		getCoords := NewGetStoredCoordinates(noTablesContainerConnection)
		_, err = getCoords(ctx, addressInfo)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("address_infos"))
	})
})

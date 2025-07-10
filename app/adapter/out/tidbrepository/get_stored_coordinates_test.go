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

	It("should return empty point if address does not exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		addressInfo := domain.AddressInfo{
			PoliticalArea: domain.PoliticalArea{
				AdminAreaLevel1: "Metropolitana",
				AdminAreaLevel2: "Santiago",
				AdminAreaLevel3: "Peñalolén",
				TimeZone:        "America/Santiago",
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
				AdminAreaLevel1: "Metropolitana",
				AdminAreaLevel2: "Santiago",
				AdminAreaLevel3: "Independencia",
				TimeZone:        "America/Santiago",
			},
			AddressLine1: "Error Forzado",
		}

		getCoords := NewGetStoredCoordinates(noTablesContainerConnection)
		_, err = getCoords(ctx, addressInfo)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("address_infos"))
	})
})

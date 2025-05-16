package tidbrepository

import (
	"context"

	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("GetStoredCoordinates", func() {

	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	It("should return coordinates if address exists", func() {
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "Coords Existentes",
			District:     "Ñuñoa",
			Location:     orb.Point{-70.6001, -33.4500},
		}

		// Guardar la dirección primero
		upsert := NewUpsertAddressInfo(connection)
		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		getCoords := NewGetStoredCoordinates(connection)
		point, err := getCoords(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())
		Expect(point.Lon()).To(BeNumerically("~", -70.6001, 0.0001))
		Expect(point.Lat()).To(BeNumerically("~", -33.4500, 0.0001))

	})

	It("should return empty point if address does not exist", func() {
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "No Existe",
			District:     "Peñalolén",
		}

		getCoords := NewGetStoredCoordinates(connection)
		point, err := getCoords(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())
		Expect(point).To(Equal(orb.Point{}))
	})

	It("should return error if address_infos table does not exist", func() {
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "Error Forzado",
			District:     "Independencia",
		}

		getCoords := NewGetStoredCoordinates(noTablesContainerConnection)
		_, err := getCoords(ctx, addressInfo)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("address_infos"))
	})
})

package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertVehicleCategory", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		// Clean the table before each test
		err := connection.DB.Exec("DELETE FROM vehicle_categories").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert vehicle category if not exists", func() {
		ctx := createOrgContext(organization1)

		vehicleCategory := domain.VehicleCategory{
			Type:                "VAN",
			MaxPackagesQuantity: 100,
		}

		upsert := NewUpsertVehicleCategory(connection)
		err := upsert(ctx, vehicleCategory)
		Expect(err).ToNot(HaveOccurred())

		var dbVehicleCategory table.VehicleCategory
		err = connection.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ?", vehicleCategory.DocID(ctx)).
			First(&dbVehicleCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicleCategory.Type).To(Equal("VAN"))
		Expect(dbVehicleCategory.MaxPackagesQuantity).To(Equal(100))
	})

	It("should update vehicle category if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.VehicleCategory{
			Type:                "TRUCK",
			MaxPackagesQuantity: 50,
		}

		upsert := NewUpsertVehicleCategory(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.VehicleCategory{
			Type:                "TRUCK",
			MaxPackagesQuantity: 75,
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbVehicleCategory table.VehicleCategory
		err = connection.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbVehicleCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicleCategory.Type).To(Equal("TRUCK"))
		Expect(dbVehicleCategory.MaxPackagesQuantity).To(Equal(75))
	})

	It("should not update vehicle category if fields are the same", func() {
		ctx := createOrgContext(organization1)

		original := domain.VehicleCategory{
			Type:                "VAN",
			MaxPackagesQuantity: 100,
		}

		upsert := NewUpsertVehicleCategory(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Try to update with same values
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		var dbVehicleCategory table.VehicleCategory
		err = connection.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbVehicleCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicleCategory.Type).To(Equal("VAN"))
		Expect(dbVehicleCategory.MaxPackagesQuantity).To(Equal(100))
	})
})

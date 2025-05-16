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

var _ = Describe("UpsertVehicle", func() {
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
		err := connection.DB.Exec("DELETE FROM vehicles").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert vehicle if not exists", func() {
		ctx := createOrgContext(organization1)

		vehicle := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
			Weight: struct {
				Value         int
				UnitOfMeasure string
			}{Value: 1000, UnitOfMeasure: "kg"},
			Insurance: struct {
				PolicyStartDate      string
				PolicyExpirationDate string
				PolicyRenewalDate    string
				MaxInsuranceCoverage struct {
					Amount   float64
					Currency string
				}
			}{
				PolicyStartDate:      "2024-01-01",
				PolicyExpirationDate: "2025-01-01",
				PolicyRenewalDate:    "2025-02-01",
				MaxInsuranceCoverage: struct {
					Amount   float64
					Currency string
				}{Amount: 10000, Currency: "CLP"},
			},
			TechnicalReview: struct {
				LastReviewDate string
				NextReviewDate string
				ReviewedBy     string
			}{
				LastReviewDate: "2023-12-01",
				NextReviewDate: "2024-12-01",
				ReviewedBy:     "Juan Mecánico",
			},
			Dimensions: struct {
				Width         float64
				Length        float64
				Height        int
				UnitOfMeasure string
			}{
				Width:         2.0,
				Length:        5.0,
				Height:        3,
				UnitOfMeasure: "m",
			},
		}

		upsert := NewUpsertVehicle(connection)
		err := upsert(ctx, vehicle)
		Expect(err).ToNot(HaveOccurred())

		var dbVehicle table.Vehicle
		err = connection.DB.WithContext(ctx).
			Table("vehicles").
			Where("document_id = ?", vehicle.DocID(ctx)).
			First(&dbVehicle).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle.Plate).To(Equal("ABC123"))
		Expect(dbVehicle.CertificateDate).To(Equal("2024-01-01"))
	})

	It("should update vehicle if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
		}

		upsert := NewUpsertVehicle(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-06-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 75,
			},
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbVehicle table.Vehicle
		err = connection.DB.WithContext(ctx).
			Table("vehicles").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbVehicle).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle.CertificateDate).To(Equal("2024-06-01"))
	})

	It("should not update vehicle if fields are the same", func() {
		ctx := createOrgContext(organization1)

		original := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
		}

		upsert := NewUpsertVehicle(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Try to update with same values
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		var dbVehicle table.Vehicle
		err = connection.DB.WithContext(ctx).
			Table("vehicles").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbVehicle).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbVehicle.Plate).To(Equal("ABC123"))
		Expect(dbVehicle.CertificateDate).To(Equal("2024-01-01"))
	})

	It("should allow same vehicle for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		vehicle1 := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
		}

		vehicle2 := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
		}

		upsert := NewUpsertVehicle(connection)

		err := upsert(ctx1, vehicle1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, vehicle2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("vehicles").
			Where("plate = ?", vehicle1.Plate).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(vehicle1.DocID(ctx1)).ToNot(Equal(vehicle2.DocID(ctx2)))
	})

	It("should fail if database has no vehicles table", func() {
		ctx := createOrgContext(organization1)

		vehicle := domain.Vehicle{
			Plate:           "ABC123",
			CertificateDate: "2024-01-01",
			VehicleCategory: domain.VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 50,
			},
		}

		upsert := NewUpsertVehicle(noTablesContainerConnection)
		err := upsert(ctx, vehicle)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("vehicles"))
	})
})

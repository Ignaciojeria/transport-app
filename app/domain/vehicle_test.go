package domain

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/biter777/countries"
)

var _ = Describe("Vehicle", func() {
	var org = Organization{ID: 1, Country: countries.CL}

	Describe("DocID", func() {
		It("should generate DocumentID based on Organization and Plate", func() {
			vehicle := Vehicle{
				Headers: Headers{Organization: org},
				Plate:   "ABC123",
			}

			Expect(vehicle.DocID()).To(Equal(Hash(org, "ABC123")))
		})
	})

	Describe("UpdateIfChanged", func() {
		var base Vehicle

		BeforeEach(func() {
			base = Vehicle{
				Headers:         Headers{Organization: org},
				Plate:           "ABC123",
				IsActive:        true,
				CertificateDate: "2024-01-01",
				VehicleCategory: VehicleCategory{
					Organization:        org,
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
		})

		It("should update plate and certificate date", func() {
			updated := base.UpdateIfChanged(Vehicle{
				Plate:           "XYZ789",
				CertificateDate: "2024-06-01",
			})

			Expect(updated.Plate).To(Equal("XYZ789"))
			Expect(updated.CertificateDate).To(Equal("2024-06-01"))
		})

		It("should update weight if non-zero", func() {
			updated := base.UpdateIfChanged(Vehicle{
				Weight: struct {
					Value         int
					UnitOfMeasure string
				}{Value: 2000, UnitOfMeasure: "ton"},
			})

			Expect(updated.Weight.Value).To(Equal(2000))
			Expect(updated.Weight.UnitOfMeasure).To(Equal("ton"))
		})

		It("should update insurance values", func() {
			updated := base.UpdateIfChanged(Vehicle{
				Insurance: struct {
					PolicyStartDate      string
					PolicyExpirationDate string
					PolicyRenewalDate    string
					MaxInsuranceCoverage struct {
						Amount   float64
						Currency string
					}
				}{
					PolicyStartDate:      "2024-02-01",
					PolicyExpirationDate: "2025-02-01",
					PolicyRenewalDate:    "2025-03-01",
					MaxInsuranceCoverage: struct {
						Amount   float64
						Currency string
					}{Amount: 20000, Currency: "USD"},
				},
			})

			Expect(updated.Insurance.PolicyStartDate).To(Equal("2024-02-01"))
			Expect(updated.Insurance.MaxInsuranceCoverage.Amount).To(Equal(20000.0))
			Expect(updated.Insurance.MaxInsuranceCoverage.Currency).To(Equal("USD"))
		})

		It("should update technical review info", func() {
			updated := base.UpdateIfChanged(Vehicle{
				TechnicalReview: struct {
					LastReviewDate string
					NextReviewDate string
					ReviewedBy     string
				}{
					LastReviewDate: "2024-01-10",
					NextReviewDate: "2025-01-10",
					ReviewedBy:     "Pedro Técnico",
				},
			})

			Expect(updated.TechnicalReview.LastReviewDate).To(Equal("2024-01-10"))
			Expect(updated.TechnicalReview.ReviewedBy).To(Equal("Pedro Técnico"))
		})

		It("should update dimensions", func() {
			updated := base.UpdateIfChanged(Vehicle{
				Dimensions: struct {
					Width         float64
					Length        float64
					Height        int
					UnitOfMeasure string
				}{
					Width:         2.5,
					Length:        6.0,
					Height:        4,
					UnitOfMeasure: "cm",
				},
			})

			Expect(updated.Dimensions.Width).To(Equal(2.5))
			Expect(updated.Dimensions.UnitOfMeasure).To(Equal("cm"))
		})

		It("should update vehicle category if type is set", func() {
			newCat := VehicleCategory{Organization: org, Type: "TRUCK", MaxPackagesQuantity: 100}
			updated := base.UpdateIfChanged(Vehicle{
				VehicleCategory: newCat,
			})

			Expect(updated.VehicleCategory.Type).To(Equal("TRUCK"))
			Expect(updated.VehicleCategory.MaxPackagesQuantity).To(Equal(100))
		})

		It("should not change values if no update data is provided", func() {
			updated := base.UpdateIfChanged(Vehicle{})
			Expect(updated).To(Equal(base))
		})
	})
})

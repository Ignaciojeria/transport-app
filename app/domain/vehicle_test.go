package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vehicle", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	Describe("DocID", func() {
		It("should generate DocumentID based on context and Plate", func() {
			vehicle := Vehicle{
				Plate: "ABC123",
			}

			Expect(vehicle.DocID(ctx)).To(Equal(Hash(ctx, "ABC123")))
		})

		It("should generate different IDs for different contexts", func() {
			ctx1 := buildCtx("org1", "CL")
			ctx2 := buildCtx("org2", "AR")

			vehicle := Vehicle{
				Plate: "ABC123",
			}

			Expect(vehicle.DocID(ctx1)).ToNot(Equal(vehicle.DocID(ctx2)))
		})
	})

	Describe("UpdateIfChanged", func() {
		var base Vehicle

		BeforeEach(func() {
			base = Vehicle{
				Plate:           "ABC123",
				IsActive:        true,
				CertificateDate: "2024-01-01",
				VehicleCategory: VehicleCategory{
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
			updated, changed := base.UpdateIfChanged(Vehicle{
				Plate:           "XYZ789",
				CertificateDate: "2024-06-01",
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Plate).To(Equal("XYZ789"))
			Expect(updated.CertificateDate).To(Equal("2024-06-01"))
		})

		It("should update weight if non-zero", func() {
			updated, changed := base.UpdateIfChanged(Vehicle{
				Weight: struct {
					Value         int
					UnitOfMeasure string
				}{Value: 2000, UnitOfMeasure: "ton"},
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Weight.Value).To(Equal(2000))
			Expect(updated.Weight.UnitOfMeasure).To(Equal("ton"))
		})

		It("should update insurance values", func() {
			updated, changed := base.UpdateIfChanged(Vehicle{
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

			Expect(changed).To(BeTrue())
			Expect(updated.Insurance.PolicyStartDate).To(Equal("2024-02-01"))
			Expect(updated.Insurance.MaxInsuranceCoverage.Amount).To(Equal(20000.0))
			Expect(updated.Insurance.MaxInsuranceCoverage.Currency).To(Equal("USD"))
		})

		It("should update technical review info", func() {
			updated, changed := base.UpdateIfChanged(Vehicle{
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

			Expect(changed).To(BeTrue())
			Expect(updated.TechnicalReview.LastReviewDate).To(Equal("2024-01-10"))
			Expect(updated.TechnicalReview.ReviewedBy).To(Equal("Pedro Técnico"))
		})

		It("should update dimensions", func() {
			updated, changed := base.UpdateIfChanged(Vehicle{
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

			Expect(changed).To(BeTrue())
			Expect(updated.Dimensions.Width).To(Equal(2.5))
			Expect(updated.Dimensions.UnitOfMeasure).To(Equal("cm"))
		})

		It("should update vehicle category if type is set", func() {
			newCat := VehicleCategory{Type: "TRUCK", MaxPackagesQuantity: 100}
			updated, changed := base.UpdateIfChanged(Vehicle{
				VehicleCategory: newCat,
			})

			Expect(changed).To(BeTrue())
			Expect(updated.VehicleCategory.Type).To(Equal("TRUCK"))
			Expect(updated.VehicleCategory.MaxPackagesQuantity).To(Equal(100))
		})

		It("should not report changes if provided values are the same", func() {
			updated, changed := base.UpdateIfChanged(Vehicle{
				Plate:           "ABC123",     // Mismo valor que base
				CertificateDate: "2024-01-01", // Mismo valor que base
				IsActive:        true,         // Mismo valor que base
			})

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(base))
		})

		It("should only update IsActive when it changes", func() {
			updated, changed := base.UpdateIfChanged(Vehicle{
				IsActive: false, // Valor diferente al original (true)
			})
			Expect(changed).To(BeTrue())
			Expect(updated.IsActive).To(BeFalse())
		})

		It("should update vehicle category MaxPackagesQuantity even if Type is the same", func() {
			newCat := VehicleCategory{Type: "VAN", MaxPackagesQuantity: 75}
			updated, changed := base.UpdateIfChanged(Vehicle{
				VehicleCategory: newCat,
			})

			Expect(changed).To(BeTrue())
			Expect(updated.VehicleCategory.Type).To(Equal("VAN"))
			Expect(updated.VehicleCategory.MaxPackagesQuantity).To(Equal(75))
		})
	})
})

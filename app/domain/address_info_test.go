package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("AddressInfo", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	Describe("ReferenceID / DocID", func() {
		It("should generate different reference IDs for different address lines", func() {
			addr1 := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
			}
			addr2 := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 5678", // distinto
			}

			Expect(addr1.DocID(ctx)).ToNot(Equal(addr2.DocID(ctx)))
		})

		It("should generate same ReferenceID for same input", func() {
			addr1 := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
			}
			addr2 := addr1

			Expect(addr1.DocID(ctx)).To(Equal(addr2.DocID(ctx)))
		})

		It("should confirm hash only includes specified fields", func() {
			addr1 := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
				ZipCode:      "7500000",
			}
			addr2 := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
				ZipCode:      "7550000", // No debería afectar el hash
			}

			Expect(addr1.DocID(ctx)).To(Equal(addr2.DocID(ctx)))
		})
	})

	Describe("UpdateIfChanged", func() {
		var original AddressInfo

		BeforeEach(func() {
			original = AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
				ZipCode:      "7500000",
				Coordinates: Coordinates{
					Point:  orb.Point{-33.4372, -70.6506},
					Source: "geocoding",
					Confidence: CoordinatesConfidence{
						Level:   0.9,
						Message: "High confidence",
						Reason:  "Exact match",
					},
				},
			}
		})

		It("should mark as changed when AddressLine1 is updated", func() {
			input := AddressInfo{AddressLine1: "Av Las Condes 5678"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine1).To(Equal("Av Las Condes 5678"))
		})

		It("should mark as changed when Coordinates are updated", func() {
			newPoint := orb.Point{-33.4400, -70.6600}
			input := AddressInfo{
				Coordinates: Coordinates{
					Point:  newPoint,
					Source: "geocoding",
					Confidence: CoordinatesConfidence{
						Level:   0.8,
						Message: "Medium confidence",
						Reason:  "Approximate match",
					},
				},
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Coordinates.Point).To(Equal(newPoint))
		})

		It("should mark as changed when confidence level is updated", func() {
			input := AddressInfo{
				Coordinates: Coordinates{
					Confidence: CoordinatesConfidence{
						Level:   0.7,
						Message: "Updated confidence",
						Reason:  "New reason",
					},
				},
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Coordinates.Confidence.Level).To(Equal(0.7))
			Expect(updated.Coordinates.Confidence.Message).To(Equal("Updated confidence"))
			Expect(updated.Coordinates.Confidence.Reason).To(Equal("New reason"))
		})

		It("should mark as changed when confidence level is zero", func() {
			input := AddressInfo{
				Coordinates: Coordinates{
					Confidence: CoordinatesConfidence{
						Level:   0.0,
						Message: "Zero confidence",
						Reason:  "Zero reason",
					},
				},
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.Coordinates.Confidence.Level).To(Equal(0.0))
			Expect(updated.Coordinates.Confidence.Message).To(Equal("Zero confidence"))
			Expect(updated.Coordinates.Confidence.Reason).To(Equal("Zero reason"))
		})

		It("should not mark as changed when confidence message is empty", func() {
			input := AddressInfo{
				Coordinates: Coordinates{
					Confidence: CoordinatesConfidence{
						Level:   original.Coordinates.Confidence.Level,
						Message: "",
						Reason:  original.Coordinates.Confidence.Reason,
					},
				},
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated.Coordinates.Confidence.Message).To(Equal(original.Coordinates.Confidence.Message))
		})

		It("should not mark as changed when confidence reason is empty", func() {
			input := AddressInfo{
				Coordinates: Coordinates{
					Confidence: CoordinatesConfidence{
						Level:   original.Coordinates.Confidence.Level,
						Message: original.Coordinates.Confidence.Message,
						Reason:  "",
					},
				},
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated.Coordinates.Confidence.Reason).To(Equal(original.Coordinates.Confidence.Reason))
		})

		It("should mark as changed when State is updated", func() {
			input := AddressInfo{PoliticalArea: PoliticalArea{AdminAreaLevel1: "Valparaíso"}}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.PoliticalArea.AdminAreaLevel1).To(Equal("Valparaíso"))
		})

		It("should mark as changed when Province is updated", func() {
			input := AddressInfo{PoliticalArea: PoliticalArea{AdminAreaLevel2: "Valparaíso"}}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.PoliticalArea.AdminAreaLevel2).To(Equal("Valparaíso"))
		})

		It("should mark as changed when District is updated", func() {
			input := AddressInfo{PoliticalArea: PoliticalArea{AdminAreaLevel3: "Las Condes"}}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.PoliticalArea.AdminAreaLevel3).To(Equal("Las Condes"))
		})

		It("should mark as changed when ZipCode is updated", func() {
			input := AddressInfo{ZipCode: "7550000"}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.ZipCode).To(Equal("7550000"))
		})

		It("should mark as changed when TimeZone is updated", func() {
			input := AddressInfo{PoliticalArea: PoliticalArea{TimeZone: "America/Argentina/Buenos_Aires"}}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.PoliticalArea.TimeZone).To(Equal("America/Argentina/Buenos_Aires"))
		})

		It("should not change when all fields are the same", func() {
			input := original
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore empty fields in input", func() {
			input := AddressInfo{
				Coordinates: Coordinates{
					Confidence: CoordinatesConfidence{
						Level:   original.Coordinates.Confidence.Level,
						Message: "",
						Reason:  "",
					},
				},
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeFalse())
			Expect(updated.Coordinates.Confidence.Level).To(Equal(original.Coordinates.Confidence.Level))
			Expect(updated.Coordinates.Confidence.Message).To(Equal(original.Coordinates.Confidence.Message))
			Expect(updated.Coordinates.Confidence.Reason).To(Equal(original.Coordinates.Confidence.Reason))
		})

		It("should update multiple fields at once", func() {
			input := AddressInfo{
				AddressLine1: "Av Las Condes 5678",
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Las Condes",
					TimeZone:        "America/Santiago",
				},
				ZipCode: "7550000",
			}
			updated, changed := original.UpdateIfChanged(input)
			Expect(changed).To(BeTrue())
			Expect(updated.AddressLine1).To(Equal("Av Las Condes 5678"))
			Expect(updated.PoliticalArea.AdminAreaLevel1).To(Equal("Metropolitana"))
			Expect(updated.PoliticalArea.AdminAreaLevel2).To(Equal("Santiago"))
			Expect(updated.PoliticalArea.AdminAreaLevel3).To(Equal("Las Condes"))
			Expect(updated.PoliticalArea.TimeZone).To(Equal("America/Santiago"))
			Expect(updated.ZipCode).To(Equal("7550000"))
		})
	})

	Describe("Normalize", func() {
		It("should normalize casing and spacing", func() {
			addr := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "   Metropolitana ",
					AdminAreaLevel2: " SANTIAGO  ",
					AdminAreaLevel3: "PROVIDENCIA   ",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "   AVENIDA   PROVIDENCIA   1234  ",
			}

			addr.ToLowerAndRemovePunctuation()

			Expect(addr.AddressLine1).To(Equal("avenida providencia 1234"))
			Expect(addr.PoliticalArea.AdminAreaLevel1).To(Equal("metropolitana"))
			Expect(addr.PoliticalArea.AdminAreaLevel2).To(Equal("santiago"))
			Expect(addr.PoliticalArea.AdminAreaLevel3).To(Equal("providencia"))
			Expect(addr.PoliticalArea.TimeZone).To(Equal("america/santiago"))
		})
	})

	Describe("FullAddress", func() {
		It("should concatenate address fields correctly", func() {
			addr := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
				ZipCode:      "7500000",
			}

			fullAddress := addr.FullAddress()

			Expect(fullAddress).To(ContainSubstring("Av Providencia 1234"))
			Expect(fullAddress).To(ContainSubstring("Providencia"))
			Expect(fullAddress).To(ContainSubstring("Santiago"))
			Expect(fullAddress).To(ContainSubstring("Metropolitana"))
			Expect(fullAddress).To(ContainSubstring("7500000"))
		})

		It("should only include specified fields in the address format", func() {
			addr := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
				ZipCode:      "7500000",
			}

			fullAddress := addr.FullAddress()
			Expect(fullAddress).To(Equal("Av Providencia 1234, Metropolitana, Santiago, Providencia, 7500000"))
		})

		It("should correctly handle empty fields", func() {
			addr := AddressInfo{
				PoliticalArea: PoliticalArea{
					AdminAreaLevel1: "Metropolitana",
					AdminAreaLevel2: "Santiago",
					AdminAreaLevel3: "Providencia",
					TimeZone:        "America/Santiago",
				},
				AddressLine1: "Av Providencia 1234",
			}

			fullAddress := addr.FullAddress()
			Expect(fullAddress).To(Equal("Av Providencia 1234, Metropolitana, Santiago, Providencia"))
		})

		It("should return just AddressLine1 if other fields are empty", func() {
			addr := AddressInfo{AddressLine1: "Av Providencia 1234"}
			Expect(addr.FullAddress()).To(Equal("Av Providencia 1234"))
		})
	})

	Describe("concatenateWithCommas", func() {
		It("should concatenate multiple non-empty strings with commas", func() {
			Expect(concatenateWithCommas("uno", "dos", "tres")).To(Equal("uno, dos, tres"))
		})

		It("should skip empty strings", func() {
			Expect(concatenateWithCommas("uno", "", "tres")).To(Equal("uno, tres"))
		})

		It("should return a single value without comma if only one non-empty string", func() {
			Expect(concatenateWithCommas("", "", "único")).To(Equal("único"))
		})

		It("should return empty string if all inputs are empty", func() {
			Expect(concatenateWithCommas("", "", "")).To(BeEmpty())
		})

		It("should return empty string if no values are passed", func() {
			Expect(concatenateWithCommas()).To(BeEmpty())
		})
	})
})

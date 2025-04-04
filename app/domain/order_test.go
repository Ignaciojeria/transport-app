package domain

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order UpdateIfChanged", func() {
	It("should return updated order and true if DeliveryInstructions changed", func() {
		original := Order{
			DeliveryInstructions: "Dejar en portería",
		}

		updated := Order{
			DeliveryInstructions: "Dejar con conserje",
		}

		result, changed := original.UpdateIfChanged(updated)

		Expect(changed).To(BeTrue())
		Expect(result.DeliveryInstructions).To(Equal("Dejar con conserje"))
	})

	It("should return same order and false if no fields changed", func() {
		original := Order{
			DeliveryInstructions: "Sin cambios",
		}

		updated := Order{}

		result, changed := original.UpdateIfChanged(updated)

		Expect(changed).To(BeFalse())
		Expect(result).To(Equal(original))
	})

	It("should return true when PromisedDate fields change", func() {
		newStart := time.Now()
		newEnd := newStart.Add(48 * time.Hour)

		updated := Order{
			PromisedDate: PromisedDate{
				DateRange: DateRange{
					StartDate: newStart,
					EndDate:   newEnd,
				},
				TimeRange: TimeRange{
					StartTime: "08:00",
					EndTime:   "18:00",
				},
				ServiceCategory: "next-day",
			},
		}

		_, changed := Order{}.UpdateIfChanged(updated)
		Expect(changed).To(BeTrue())
	})

	It("should return true when CollectAvailabilityDate changes", func() {
		date := time.Now()
		updated := Order{
			CollectAvailabilityDate: CollectAvailabilityDate{
				Date: date,
				TimeRange: TimeRange{
					StartTime: "10:00",
					EndTime:   "12:00",
				},
			},
		}

		_, changed := Order{}.UpdateIfChanged(updated)
		Expect(changed).To(BeTrue())
	})

	It("should return true when Items, Packages, References or TransportRequirements are replaced", func() {
		updated := Order{
			Items:                 []Item{{ReferenceID: "ITEM001"}},
			Packages:              []Package{{Lpn: "PKG001"}},
			References:            []Reference{{Type: "X", Value: "1"}},
			TransportRequirements: []Reference{{Type: "TEMP", Value: "COND"}},
		}

		_, changed := Order{}.UpdateIfChanged(updated)
		Expect(changed).To(BeTrue())
	})

	It("should return false when attempting to update with empty slices or zero values", func() {
		original := Order{
			Items: []Item{{ReferenceID: "ITEM001"}},
		}

		// newOrder no tiene nada nuevo (vacío)
		updated := Order{}

		result, changed := original.UpdateIfChanged(updated)

		Expect(changed).To(BeFalse())
		Expect(result).To(Equal(original))
	})
})

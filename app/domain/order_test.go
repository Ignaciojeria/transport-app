package domain

import (
	"context"
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

	It("should return true when Packages, References or TransportRequirements are replaced", func() {
		updated := Order{
			Packages: []Package{
				{
					Lpn: "PKG001",
					Items: []Item{
						{Sku: "ITEM001"},
					},
				},
			},
			References:            []Reference{{Type: "X", Value: "1"}},
			TransportRequirements: []Reference{{Type: "TEMP", Value: "COND"}},
		}

		_, changed := Order{}.UpdateIfChanged(updated)
		Expect(changed).To(BeTrue())
	})

	It("should return false when attempting to update with empty slices or zero values", func() {
		original := Order{
			Packages: []Package{
				{
					Lpn: "PKG001",
					Items: []Item{
						{Sku: "ITEM001"},
					},
				},
			},
		}

		// newOrder no tiene nada nuevo (vacío)
		updated := Order{}

		result, changed := original.UpdateIfChanged(updated)

		Expect(changed).To(BeFalse())
		Expect(result).To(Equal(original))
	})
})

var _ = Describe("Order Validate", func() {
	It("should return nil when all validations pass", func() {
		order := Order{
			PromisedDate: PromisedDate{
				DateRange: DateRange{
					StartDate: time.Now(),
					EndDate:   time.Now().Add(24 * time.Hour),
				},
				TimeRange: TimeRange{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
			},
			CollectAvailabilityDate: CollectAvailabilityDate{
				Date: time.Now(),
				TimeRange: TimeRange{
					StartTime: "10:00",
					EndTime:   "12:00",
				},
			},
			Packages: []Package{
				{
					Lpn: "PKG001",
					Items: []Item{
						{Sku: "ITEM001"},
					},
				},
			},
		}

		Expect(order.Validate()).To(Succeed())
	})

	It("should fail if promised date is invalid", func() {
		order := Order{
			PromisedDate: PromisedDate{
				DateRange: DateRange{
					StartDate: time.Now(),
					EndDate:   time.Now().Add(-24 * time.Hour),
				},
			},
		}
		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("PromisedDate"))
	})

	It("should fail if collect availability time is invalid", func() {
		order := Order{
			CollectAvailabilityDate: CollectAvailabilityDate{
				TimeRange: TimeRange{
					StartTime: "not-time",
				},
			},
		}
		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("CollectAvailabilityDate"))
	})

	// Eliminamos el test de validación de paquetes que ya no es necesario
})

var _ = Describe("Order Origin and Destination Comparison", func() {
	It("should return true when contacts are equal", func() {
		contact := Contact{
			FullName:     "John Doe",
			PrimaryEmail: "john@example.com",
			PrimaryPhone: "123456789",
			NationalID:   "12345678-9",
		}

		order := Order{
			Origin:      NodeInfo{AddressInfo: AddressInfo{Contact: contact}},
			Destination: NodeInfo{AddressInfo: AddressInfo{Contact: contact}},
		}

		Expect(order.IsOriginAndDestinationContactEqual()).To(BeTrue())
	})

	It("should return false when contacts differ", func() {
		order := Order{
			Origin: NodeInfo{AddressInfo: AddressInfo{Contact: Contact{
				FullName: "John",
			}}},
			Destination: NodeInfo{AddressInfo: AddressInfo{Contact: Contact{
				FullName: "Jane",
			}}},
		}

		Expect(order.IsOriginAndDestinationContactEqual()).To(BeFalse())
	})

	It("should return true when addresses are equal", func() {
		address := AddressInfo{AddressLine1: "Av. Siempre Viva 742"}

		order := Order{
			Origin:      NodeInfo{AddressInfo: address},
			Destination: NodeInfo{AddressInfo: address},
		}

		Expect(order.IsOriginAndDestinationAddressEqual()).To(BeTrue())
	})

	It("should return false when addresses differ", func() {
		order := Order{
			Origin:      NodeInfo{AddressInfo: AddressInfo{AddressLine1: "Calle A"}},
			Destination: NodeInfo{AddressInfo: AddressInfo{AddressLine1: "Calle B"}},
		}

		Expect(order.IsOriginAndDestinationAddressEqual()).To(BeFalse())
	})

	It("should return true when node reference IDs are equal", func() {
		order := Order{
			Origin:      NodeInfo{ReferenceID: "NODE-123"},
			Destination: NodeInfo{ReferenceID: "NODE-123"},
		}

		Expect(order.IsOriginAndDestinationNodeEqual()).To(BeTrue())
	})

	It("should return false when node reference IDs differ", func() {
		order := Order{
			Origin:      NodeInfo{ReferenceID: "NODE-123"},
			Destination: NodeInfo{ReferenceID: "NODE-999"},
		}

		Expect(order.IsOriginAndDestinationNodeEqual()).To(BeFalse())
	})
})

var _ = Describe("Order Validate - Additional Cases", func() {
	It("should fail if PromisedDate startTime has invalid format", func() {
		order := Order{
			PromisedDate: PromisedDate{
				TimeRange: TimeRange{
					StartTime: "25:00", // inválido
				},
			},
		}
		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("promised delivery startTime"))
	})

	It("should fail if PromisedDate endTime has invalid format", func() {
		order := Order{
			PromisedDate: PromisedDate{
				TimeRange: TimeRange{
					EndTime: "99:99", // inválido
				},
			},
		}
		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("promised delivery endTime"))
	})

	It("should fail if CollectAvailabilityDate endTime has invalid format", func() {
		order := Order{
			CollectAvailabilityDate: CollectAvailabilityDate{
				TimeRange: TimeRange{
					EndTime: "ab:cd", // inválido
				},
			},
		}
		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("collect endTime"))
	})

	// Eliminamos todos los test relacionados con la validación de paquetes e ItemReferences
})

var _ = Describe("Order DocID", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	It("should return the hash of the context and ReferenceID", func() {
		order := Order{
			ReferenceID: "REF-0001",
		}

		expected := HashByTenant(ctx, "REF-0001")
		Expect(order.DocID(ctx)).To(Equal(expected))
	})

	It("should generate different IDs for different contexts", func() {
		order := Order{
			ReferenceID: "REF-0001",
		}

		ctx1 := buildCtx("org1", "CL")
		ctx2 := buildCtx("org2", "AR")

		Expect(order.DocID(ctx1)).ToNot(Equal(order.DocID(ctx2)))
	})

	It("should generate different IDs for different reference IDs", func() {
		order1 := Order{
			ReferenceID: "REF-0001",
		}
		order2 := Order{
			ReferenceID: "REF-0002",
		}

		Expect(order1.DocID(ctx)).ToNot(Equal(order2.DocID(ctx)))
	})

	It("should fail if a package without LPN has no items", func() {
		order := Order{
			Packages: []Package{
				{Lpn: "", Items: nil}, // ❌ sin ítems
			},
		}

		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("must have at least one item"))
	})

	It("should fail if a package without LPN has an item with empty SKU", func() {
		order := Order{
			Packages: []Package{
				{
					Lpn: "",
					Items: []Item{
						{Sku: ""},
					},
				},
			},
		}

		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("must have a non-empty SKU"))
	})

	It("should succeed if a package without LPN has item with valid SKU", func() {
		order := Order{
			Packages: []Package{
				{
					Lpn: "",
					Items: []Item{
						{Sku: "ITEM001"},
					},
				},
			},
		}

		Expect(order.Validate()).To(Succeed())
	})

})

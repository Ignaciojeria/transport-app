package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VehicleCategory", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	Describe("DocID", func() {
		It("should generate unique ID based on context and Type", func() {
			vc1 := VehicleCategory{
				Type: "VAN",
			}
			vc2 := VehicleCategory{
				Type: "TRUCK",
			}

			Expect(vc1.DocID(ctx)).To(Equal(HashByTenant(ctx, "VAN")))
			Expect(vc1.DocID(ctx)).ToNot(Equal(vc2.DocID(ctx)))
		})

		It("should generate different IDs for different contexts", func() {
			ctx1 := buildCtx("org1", "CL")
			ctx2 := buildCtx("org2", "AR")

			vc := VehicleCategory{
				Type: "VAN",
			}

			Expect(vc.DocID(ctx1)).ToNot(Equal(vc.DocID(ctx2)))
		})
	})

	Describe("UpdateIfChanged", func() {
		var base VehicleCategory

		BeforeEach(func() {
			base = VehicleCategory{
				Type:                "VAN",
				MaxPackagesQuantity: 100,
			}
		})

		It("should update MaxPackagesQuantity if new value is non-zero", func() {
			updated, changed := base.UpdateIfChanged(VehicleCategory{
				MaxPackagesQuantity: 200,
			})

			Expect(changed).To(BeTrue())
			Expect(updated.MaxPackagesQuantity).To(Equal(200))
		})

		It("should not update MaxPackagesQuantity if new value is zero", func() {
			updated, changed := base.UpdateIfChanged(VehicleCategory{
				MaxPackagesQuantity: 0,
			})

			Expect(changed).To(BeFalse())
			Expect(updated.MaxPackagesQuantity).To(Equal(100))
		})

		It("should keep Type unchanged", func() {
			updated, changed := base.UpdateIfChanged(VehicleCategory{
				Type:                "TRUCK", // Este valor debería ser ignorado
				MaxPackagesQuantity: 300,
			})

			Expect(changed).To(BeTrue())
			Expect(updated.Type).To(Equal("VAN"))
			Expect(updated.MaxPackagesQuantity).To(Equal(300))
		})

		It("should not mark as changed when values are the same", func() {
			updated, changed := base.UpdateIfChanged(VehicleCategory{
				MaxPackagesQuantity: 100, // Mismo valor que base
			})

			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(base))
		})
	})
})

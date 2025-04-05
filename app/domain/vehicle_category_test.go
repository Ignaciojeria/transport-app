package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VehicleCategory", func() {
	var org = Organization{ID: 1, Country: countries.CL}

	Describe("DocID", func() {
		It("should generate unique ID based on Organization and Type", func() {
			vc1 := VehicleCategory{
				Organization: org,
				Type:         "VAN",
			}
			vc2 := VehicleCategory{
				Organization: org,
				Type:         "TRUCK",
			}

			Expect(vc1.DocID()).To(Equal(Hash(org, "VAN")))
			Expect(vc1.DocID()).ToNot(Equal(vc2.DocID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		var base VehicleCategory

		BeforeEach(func() {
			base = VehicleCategory{
				Organization:        org,
				Type:                "VAN",
				MaxPackagesQuantity: 100,
			}
		})

		It("should update MaxPackagesQuantity if new value is non-zero", func() {
			updated := base.UpdateIfChanged(VehicleCategory{
				MaxPackagesQuantity: 200,
			})

			Expect(updated.MaxPackagesQuantity).To(Equal(200))
		})

		It("should not update MaxPackagesQuantity if new value is zero", func() {
			updated := base.UpdateIfChanged(VehicleCategory{
				MaxPackagesQuantity: 0,
			})

			Expect(updated.MaxPackagesQuantity).To(Equal(100))
		})

		It("should keep Organization and Type unchanged", func() {
			updated := base.UpdateIfChanged(VehicleCategory{
				MaxPackagesQuantity: 300,
			})

			Expect(updated.Organization).To(Equal(org))
			Expect(updated.Type).To(Equal("VAN"))
		})
	})
})

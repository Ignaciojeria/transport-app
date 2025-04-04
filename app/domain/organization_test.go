package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Organization", func() {
	Describe("GetOrgKey", func() {
		It("should return the correct key", func() {
			org := Organization{
				ID:      123,
				Country: countries.CL,
			}
			Expect(org.GetOrgKey()).To(Equal("123-CL"))
		})
	})

	Describe("SetKey", func() {
		It("should set ID and Country correctly from a valid key", func() {
			var org Organization
			err := org.SetKey("456-AR")
			Expect(err).ToNot(HaveOccurred())
			Expect(org.ID).To(Equal(int64(456)))
			Expect(org.Country).To(Equal(countries.AR))
		})

		It("should return error if key format is invalid (missing dash)", func() {
			var org Organization
			err := org.SetKey("456AR")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid key format"))
		})

		It("should return error if ID part is not a number", func() {
			var org Organization
			err := org.SetKey("abc-CL")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid ID in key"))
		})

		It("should set unknown country if country code is invalid", func() {
			var org Organization
			err := org.SetKey("789-XXX")
			Expect(err).ToNot(HaveOccurred())
			Expect(org.ID).To(Equal(int64(789)))
			Expect(org.Country).To(Equal(countries.Unknown))
		})
	})
})

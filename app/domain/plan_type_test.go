package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlanType", func() {

	var (
		org1 = Organization{ID: 1, Country: countries.CL}
		org2 = Organization{ID: 2, Country: countries.AR}
	)

	Describe("DocID", func() {
		It("should generate different document IDs for different organizations", func() {
			planType1 := PlanType{Organization: org1, Value: "dispatch"}
			planType2 := PlanType{Organization: org2, Value: "dispatch"}

			Expect(planType1.DocID()).ToNot(Equal(planType2.DocID()))
		})

		It("should generate different document IDs for different plan types", func() {
			planType1 := PlanType{Organization: org1, Value: "dispatch"}
			planType2 := PlanType{Organization: org1, Value: "web"}

			Expect(planType1.DocID()).ToNot(Equal(planType2.DocID()))
		})

		It("should generate the same document ID for same org and plan type", func() {
			planType1 := PlanType{Organization: org1, Value: "pickup"}
			planType2 := PlanType{Organization: org1, Value: "pickup"}

			Expect(planType1.DocID()).To(Equal(planType2.DocID()))
		})
	})
})
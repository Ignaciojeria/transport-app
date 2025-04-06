package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlanningStatus", func() {

	var (
		org1 = Organization{ID: 1, Country: countries.CL}
		org2 = Organization{ID: 2, Country: countries.AR}
	)

	Describe("DocID", func() {
		It("should generate different document IDs for different organizations", func() {
			planningStatus1 := PlanningStatus{Organization: org1, Value: "pending"}
			planningStatus2 := PlanningStatus{Organization: org2, Value: "pending"}

			Expect(planningStatus1.DocID()).ToNot(Equal(planningStatus2.DocID()))
		})

		It("should generate different document IDs for different status values", func() {
			planningStatus1 := PlanningStatus{Organization: org1, Value: "pending"}
			planningStatus2 := PlanningStatus{Organization: org1, Value: "in_progress"}

			Expect(planningStatus1.DocID()).ToNot(Equal(planningStatus2.DocID()))
		})

		It("should generate the same document ID for same org and status value", func() {
			planningStatus1 := PlanningStatus{Organization: org1, Value: "completed"}
			planningStatus2 := PlanningStatus{Organization: org1, Value: "completed"}

			Expect(planningStatus1.DocID()).To(Equal(planningStatus2.DocID()))
		})
	})
})

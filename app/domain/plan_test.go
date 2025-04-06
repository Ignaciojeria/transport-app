package domain

import (
	"time"

	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {
	var (
		org1 = Organization{ID: 1, Country: countries.CL}
		org2 = Organization{ID: 2, Country: countries.AR}
		now  = time.Now()
		tomorrow = now.AddDate(0, 0, 1)
	)

	Describe("DocID", func() {
		It("should generate different document IDs for different organizations", func() {
			plan1 := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				Headers:     Headers{Organization: org2},
				ReferenceID: "PLAN-001",
			}

			Expect(plan1.DocID()).ToNot(Equal(plan2.DocID()))
		})

		It("should generate different document IDs for different reference IDs", func() {
			plan1 := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-002",
			}

			Expect(plan1.DocID()).ToNot(Equal(plan2.DocID()))
		})

		It("should generate the same document ID for same org and reference ID", func() {
			plan1 := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
			}

			Expect(plan1.DocID()).To(Equal(plan2.DocID()))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should update planned date if provided", func() {
			original := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: tomorrow,
			}

			updated := original.UpdateIfChanged(newPlan)
			Expect(updated.ReferenceID).To(Equal("PLAN-001")) // No debe cambiar
			Expect(updated.PlannedDate).To(Equal(tomorrow))
		})

		It("should not update planned date when empty value is provided", func() {
			original := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: time.Time{}, // zero value
			}

			updated := original.UpdateIfChanged(newPlan)
			Expect(updated.ReferenceID).To(Equal("PLAN-001"))
			Expect(updated.PlannedDate).To(Equal(now))
		})

		It("should ignore ReferenceID even if provided in newPlan", func() {
			original := Plan{
				Headers:     Headers{Organization: org1},
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				ReferenceID: "PLAN-002", // Esto no debería afectar
				PlannedDate: tomorrow,
			}

			updated := original.UpdateIfChanged(newPlan)
			Expect(updated.ReferenceID).To(Equal("PLAN-001")) // Debe mantener el original
			Expect(updated.PlannedDate).To(Equal(tomorrow))
		})
	})
})
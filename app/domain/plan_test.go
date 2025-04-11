package domain

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {
	var (
		ctx1, ctx2 context.Context
		now        = time.Now()
		tomorrow   = now.AddDate(0, 0, 1)
	)

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate different document IDs for different contexts", func() {
			plan1 := Plan{
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				ReferenceID: "PLAN-001",
			}

			Expect(plan1.DocID(ctx1)).ToNot(Equal(plan2.DocID(ctx2)))
		})

		It("should generate different document IDs for different reference IDs", func() {
			plan1 := Plan{
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				ReferenceID: "PLAN-002",
			}

			Expect(plan1.DocID(ctx1)).ToNot(Equal(plan2.DocID(ctx1)))
		})

		It("should generate the same document ID for same context and reference ID", func() {
			plan1 := Plan{
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				ReferenceID: "PLAN-001",
			}

			Expect(plan1.DocID(ctx1)).To(Equal(plan2.DocID(ctx1)))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should update planned date if provided and different", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: tomorrow,
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeTrue())
			Expect(updated.ReferenceID).To(Equal("PLAN-001")) // No debe cambiar
			Expect(updated.PlannedDate).To(Equal(tomorrow))
		})

		It("should not update planned date when empty value is provided", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: time.Time{}, // zero value
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeFalse())
			Expect(updated.ReferenceID).To(Equal("PLAN-001"))
			Expect(updated.PlannedDate).To(Equal(now))
		})

		It("should not mark as changed if same value is provided", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: now, // mismo valor
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore ReferenceID even if provided in newPlan", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				ReferenceID: "PLAN-002", // Esto no debería afectar
				PlannedDate: tomorrow,
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeTrue())
			Expect(updated.ReferenceID).To(Equal("PLAN-001")) // Debe mantener el original
			Expect(updated.PlannedDate).To(Equal(tomorrow))
		})
	})
})

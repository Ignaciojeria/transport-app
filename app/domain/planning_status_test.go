package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlanningStatus", func() {
	var (
		ctx1 context.Context
		ctx2 context.Context
	)

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate different document IDs for different contexts", func() {
			planningStatus := PlanningStatus{Value: "pending"}

			Expect(planningStatus.DocID(ctx1)).ToNot(Equal(planningStatus.DocID(ctx2)))
		})

		It("should generate different document IDs for different status values", func() {
			planningStatus1 := PlanningStatus{Value: "pending"}
			planningStatus2 := PlanningStatus{Value: "in_progress"}

			Expect(planningStatus1.DocID(ctx1)).ToNot(Equal(planningStatus2.DocID(ctx1)))
		})

		It("should generate the same document ID for same context and status value", func() {
			planningStatus1 := PlanningStatus{Value: "completed"}
			planningStatus2 := PlanningStatus{Value: "completed"}

			Expect(planningStatus1.DocID(ctx1)).To(Equal(planningStatus2.DocID(ctx1)))
		})
	})
})

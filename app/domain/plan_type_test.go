package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlanType", func() {
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
			planType := PlanType{Value: "dispatch"}

			Expect(planType.DocID(ctx1)).ToNot(Equal(planType.DocID(ctx2)))
		})

		It("should generate different document IDs for different plan types", func() {
			planType1 := PlanType{Value: "dispatch"}
			planType2 := PlanType{Value: "web"}

			Expect(planType1.DocID(ctx1)).ToNot(Equal(planType2.DocID(ctx1)))
		})

		It("should generate the same document ID for same context and plan type", func() {
			planType1 := PlanType{Value: "pickup"}
			planType2 := PlanType{Value: "pickup"}

			Expect(planType1.DocID(ctx1)).To(Equal(planType2.DocID(ctx1)))
		})
	})
})

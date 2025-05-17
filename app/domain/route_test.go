package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route", func() {
	var (
		ctx1, ctx2 context.Context
		baseRoute  Route
	)

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
		baseRoute = Route{
			ReferenceID: "route-001",
			Origin: NodeInfo{
				ReferenceID: "origin-001",
				Name:        "Origin Location",
			},
			Destination: NodeInfo{
				ReferenceID: "dest-001",
				Name:        "Destination Location",
			},
			Vehicle: Vehicle{
				Plate: "ABC123",
			},
			Orders: []Order{
				{ReferenceID: "order-001"},
			},
		}
	})

	Describe("DocID", func() {
		It("should generate different IDs for different contexts", func() {
			Expect(baseRoute.DocID(ctx1)).ToNot(Equal(baseRoute.DocID(ctx2)))
		})

		It("should generate different IDs for different reference IDs", func() {
			route1 := Route{ReferenceID: "route-001"}
			route2 := Route{ReferenceID: "route-002"}

			Expect(route1.DocID(ctx1)).ToNot(Equal(route2.DocID(ctx1)))
		})

		It("should generate the same ID for same context and reference ID", func() {
			route1 := Route{ReferenceID: "route-001"}
			route2 := Route{ReferenceID: "route-001"}

			Expect(route1.DocID(ctx1)).To(Equal(route2.DocID(ctx1)))
		})
	})
})

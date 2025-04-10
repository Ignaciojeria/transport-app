package domain

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Headers", func() {
	Describe("DocID", func() {
		It("should generate ID based on context, commerce and consumer", func() {
			ctx := buildCtx("org1", "CL")

			headers := Headers{
				Commerce: "store-1",
				Consumer: "customer-1",
			}

			Expect(headers.DocID(ctx)).To(Equal(Hash(ctx, "store-1", "customer-1")))
		})

		It("should generate different IDs for different contexts", func() {
			ctx1 := buildCtx("org1", "CL")
			ctx2 := buildCtx("org2", "AR")

			headers := Headers{
				Commerce: "store-1",
				Consumer: "customer-1",
			}

			Expect(headers.DocID(ctx1)).ToNot(Equal(headers.DocID(ctx2)))
		})

		It("should generate different IDs for different Commerce values", func() {
			ctx := buildCtx("org1", "CL")

			headers1 := Headers{
				Commerce: "store-1",
				Consumer: "customer-1",
			}
			headers2 := Headers{
				Commerce: "store-2",
				Consumer: "customer-1",
			}

			Expect(headers1.DocID(ctx)).ToNot(Equal(headers2.DocID(ctx)))
		})

		It("should generate different IDs for different Consumer values", func() {
			ctx := buildCtx("org1", "CL")

			headers1 := Headers{
				Commerce: "store-1",
				Consumer: "customer-1",
			}
			headers2 := Headers{
				Commerce: "store-1",
				Consumer: "customer-2",
			}

			Expect(headers1.DocID(ctx)).ToNot(Equal(headers2.DocID(ctx)))
		})
	})
})

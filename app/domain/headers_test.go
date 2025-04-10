package domain

import (
	"context"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

func buildCtx(tenantID, country string) context.Context {
	ctx := context.Background()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID)
	cntry, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, country)
	bag, _ := baggage.New(tID, cntry)
	return baggage.ContextWithBaggage(ctx, bag)
}

var _ = Describe("Headers", func() {
	Describe("DocID", func() {
		It("should generate unique ID based on tenantID, country, commerce and consumer", func() {
			ctx1 := buildCtx("org1", "CL")
			ctx2 := buildCtx("org1", "CL")
			ctx3 := buildCtx("org2", "AR")

			headers1 := Headers{
				Commerce: "store-1",
				Consumer: "customer-1",
			}
			headers2 := Headers{
				Commerce: "store-2",
				Consumer: "customer-1",
			}
			headers3 := Headers{
				Commerce: "store-1",
				Consumer: "customer-1",
			}

			Expect(headers1.DocID(ctx1)).To(Equal(Hash(ctx1, "store-1", "customer-1")))
			Expect(headers1.DocID(ctx1)).ToNot(Equal(headers2.DocID(ctx2)))
			Expect(headers1.DocID(ctx1)).ToNot(Equal(headers3.DocID(ctx3)))
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

	Describe("UpdateIfChanged", func() {
		var baseHeaders Headers

		BeforeEach(func() {
			baseHeaders = Headers{
				Commerce: "original-store",
				Consumer: "original-customer",
			}
		})

		It("should update Consumer if new value is not empty", func() {
			newHeaders := Headers{Consumer: "new-customer"}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Consumer).To(Equal("new-customer"))
			Expect(updated.Commerce).To(Equal("original-store"))
		})

		It("should not update Consumer if new value is empty", func() {
			newHeaders := Headers{Consumer: ""}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Consumer).To(Equal("original-customer"))
		})

		It("should update Commerce if new value is not empty", func() {
			newHeaders := Headers{Commerce: "new-store"}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Commerce).To(Equal("new-store"))
			Expect(updated.Consumer).To(Equal("original-customer"))
		})

		It("should not update Commerce if new value is empty", func() {
			newHeaders := Headers{Commerce: ""}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Commerce).To(Equal("original-store"))
		})

		It("should update multiple fields at once", func() {
			newHeaders := Headers{
				Commerce: "new-store",
				Consumer: "new-customer",
			}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Commerce).To(Equal("new-store"))
			Expect(updated.Consumer).To(Equal("new-customer"))
		})

		It("should update only non-empty fields", func() {
			newHeaders := Headers{
				Commerce: "",
				Consumer: "new-customer",
			}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated.Commerce).To(Equal("original-store"))
			Expect(updated.Consumer).To(Equal("new-customer"))
		})

		It("should not change anything if all new values are empty", func() {
			newHeaders := Headers{
				Commerce: "",
				Consumer: "",
			}

			updated := baseHeaders.UpdateIfChanged(newHeaders)

			Expect(updated).To(Equal(baseHeaders))
		})
	})
})

package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("Province", func() {
	Describe("String", func() {
		It("should return the string value of the province", func() {
			p := Province("Santiago")
			Expect(p.String()).To(Equal("Santiago"))
		})
	})

	Describe("DocID", func() {
		It("should generate a deterministic DocID based on tenant and province", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			p := Province("Santiago")
			docID := p.DocID(ctx)

			orgKey := "org1-CL"
			joined := strings.Join([]string{orgKey, "Santiago"}, "|")
			sum := sha256.Sum256([]byte(joined))
			expected := hex.EncodeToString(sum[:])

			Expect(string(docID)).To(Equal(expected))
		})

		It("should return different DocIDs for different provinces", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			p1 := Province("Santiago")
			p2 := Province("Valparaíso")

			Expect(p1.DocID(ctx)).ToNot(Equal(p2.DocID(ctx)))
		})

		It("should return different DocIDs for same province with different tenants", func() {
			ctx1 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctx2 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org2",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			p := Province("Santiago")

			Expect(p.DocID(ctx1)).ToNot(Equal(p.DocID(ctx2)))
		})
	})
})

var _ = Describe("Province", func() {
	Describe("IsEmpty", func() {
		It("should return true when province is empty", func() {
			Expect(Province("").IsEmpty()).To(BeTrue())
		})

		It("should return true when province contains only spaces", func() {
			Expect(Province("   ").IsEmpty()).To(BeTrue())
		})

		It("should return false when province has content", func() {
			Expect(Province("Santiago").IsEmpty()).To(BeFalse())
		})
	})

	Describe("Equals", func() {
		It("should return true when provinces are exactly equal", func() {
			Expect(Province("Santiago").Equals(Province("Santiago"))).To(BeTrue())
		})

		It("should return false when provinces differ", func() {
			Expect(Province("Santiago").Equals(Province("Chacabuco"))).To(BeFalse())
		})

		It("should return false when one is empty", func() {
			Expect(Province("").Equals(Province("Santiago"))).To(BeFalse())
		})
	})
})

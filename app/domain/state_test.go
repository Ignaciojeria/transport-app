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

var _ = Describe("State", func() {
	Describe("String", func() {
		It("should return the string value of the state", func() {
			s := State("Región Metropolitana de Santiago")
			Expect(s.String()).To(Equal("Región Metropolitana de Santiago"))
		})
	})

	Describe("DocID", func() {
		It("should generate a deterministic DocID based on tenant and state", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			s := State("Región Metropolitana de Santiago")
			docID := s.DocID(ctx)

			orgKey := "org1-CL"
			joined := strings.Join([]string{orgKey, "Región Metropolitana de Santiago"}, "|")
			sum := sha256.Sum256([]byte(joined))
			expected := hex.EncodeToString(sum[:])

			Expect(string(docID)).To(Equal(expected))
		})

		It("should return different DocIDs for different states", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			s1 := State("Región Metropolitana de Santiago")
			s2 := State("Valparaíso")

			Expect(s1.DocID(ctx)).ToNot(Equal(s2.DocID(ctx)))
		})

		It("should return different DocIDs for same state with different tenants", func() {
			ctx1 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctx2 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org2",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			s := State("Región Metropolitana de Santiago")

			Expect(s.DocID(ctx1)).ToNot(Equal(s.DocID(ctx2)))
		})
	})

})

var _ = Describe("State", func() {
	Describe("IsEmpty", func() {
		It("should return true when state is empty", func() {
			Expect(State("").IsEmpty()).To(BeTrue())
		})

		It("should return true when state contains only spaces", func() {
			Expect(State("   ").IsEmpty()).To(BeTrue())
		})

		It("should return false when state has content", func() {
			Expect(State("Metropolitana").IsEmpty()).To(BeFalse())
		})
	})

	Describe("Equals", func() {
		It("should return true when states are exactly equal", func() {
			Expect(State("Valparaíso").Equals(State("Valparaíso"))).To(BeTrue())
		})

		It("should return false when states differ", func() {
			Expect(State("Valparaíso").Equals(State("Metropolitana"))).To(BeFalse())
		})

		It("should return false when one is empty", func() {
			Expect(State("").Equals(State("Valparaíso"))).To(BeFalse())
		})
	})
})

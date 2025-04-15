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
		It("should generate a deterministic DocID based on country and state", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "CL",
			))

			s := State("Región Metropolitana de Santiago")
			docID := s.DocID(ctx)

			joined := strings.Join([]string{"CL", "state", "Región Metropolitana de Santiago"}, "|")

			sum := sha256.Sum256([]byte(joined))
			expected := hex.EncodeToString(sum[:16])

			Expect(string(docID)).To(Equal(expected))
		})

		It("should return different DocIDs for different states", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "CL",
			))

			s1 := State("Región Metropolitana de Santiago")
			s2 := State("Valparaíso")

			Expect(s1.DocID(ctx)).ToNot(Equal(s2.DocID(ctx)))
		})

		It("should return different DocIDs for same state with different countries", func() {
			ctxCL := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctxAR := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "AR",
			))

			s := State("Región Metropolitana de Santiago")

			Expect(s.DocID(ctxCL)).ToNot(Equal(s.DocID(ctxAR)))
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

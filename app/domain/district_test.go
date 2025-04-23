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

var _ = Describe("District", func() {
	Describe("String", func() {
		It("should return the string value of the district", func() {
			d := District("Providencia")
			Expect(d.String()).To(Equal("Providencia"))
		})
	})

	Describe("DocID", func() {
		It("should generate a deterministic DocID based on country and district", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "CL",
			))

			d := District("Ñuñoa")
			docID := d.DocID(ctx)

			// Expected hash: HashByCountry(context.Background(), "CL", "Ñuñoa")
			joined := strings.Join([]string{"CL", "district", "Ñuñoa"}, "|")
			sum := sha256.Sum256([]byte(joined))
			expected := hex.EncodeToString(sum[:])

			Expect(string(docID)).To(Equal(expected))
		})

		It("should return different DocIDs for different districts", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "CL",
			))

			d1 := District("Providencia")
			d2 := District("Las Condes")

			Expect(d1.DocID(ctx)).ToNot(Equal(d2.DocID(ctx)))
		})

		It("should return different DocIDs for same district but different countries", func() {
			ctxCL := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctxAR := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantCountry, "AR",
			))

			d := District("Providencia")
			Expect(d.DocID(ctxCL)).ToNot(Equal(d.DocID(ctxAR)))
		})
	})
})

var _ = Describe("District", func() {
	Describe("IsEmpty", func() {
		It("should return true when district is empty", func() {
			Expect(District("").IsEmpty()).To(BeTrue())
		})

		It("should return true when district contains only spaces", func() {
			Expect(District("   ").IsEmpty()).To(BeTrue())
		})

		It("should return false when district has content", func() {
			Expect(District("Providencia").IsEmpty()).To(BeFalse())
		})
	})

	Describe("Equals", func() {
		It("should return true when districts are exactly equal", func() {
			Expect(District("Ñuñoa").Equals(District("Ñuñoa"))).To(BeTrue())
		})

		It("should return false when districts differ", func() {
			Expect(District("Providencia").Equals(District("Las Condes"))).To(BeFalse())
		})

		It("should return false when one is empty", func() {
			Expect(District("").Equals(District("Ñuñoa"))).To(BeFalse())
		})
	})
})

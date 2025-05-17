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
		It("should generate a deterministic DocID based on tenant and district", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			d := District("Providencia")
			docID := d.DocID(ctx)

			orgKey := "org1-CL"
			joined := strings.Join([]string{orgKey, "Providencia"}, "|")
			sum := sha256.Sum256([]byte(joined))
			expected := hex.EncodeToString(sum[:])

			Expect(string(docID)).To(Equal(expected))
		})

		It("should return different DocIDs for different districts", func() {
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			d1 := District("Providencia")
			d2 := District("Las Condes")

			Expect(d1.DocID(ctx)).ToNot(Equal(d2.DocID(ctx)))
		})

		It("should return different DocIDs for same district with different tenants", func() {
			ctx1 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org1",
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctx2 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "org2",
				sharedcontext.BaggageTenantCountry, "CL",
			))

			d := District("Providencia")

			Expect(d.DocID(ctx1)).ToNot(Equal(d.DocID(ctx2)))
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
			Expect(District("Ñuñoa").IsEmpty()).To(BeFalse())
		})
	})

	Describe("Equals", func() {
		It("should return true when districts are exactly equal", func() {
			Expect(District("Ñuñoa").Equals(District("Ñuñoa"))).To(BeTrue())
		})

		It("should return false when districts differ", func() {
			Expect(District("Ñuñoa").Equals(District("Las Condes"))).To(BeFalse())
		})

		It("should return false when one is empty", func() {
			Expect(District("").Equals(District("Ñuñoa"))).To(BeFalse())
		})
	})
})

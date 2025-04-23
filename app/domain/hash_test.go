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

var _ = Describe("HashByCountry", func() {
	It("should generate a deterministic hash based on country and input", func() {
		// Arrange
		country := "CL"
		ctx := context.Background()

		countryMember, err := baggage.NewMember(sharedcontext.BaggageTenantCountry, country)
		Expect(err).ToNot(HaveOccurred())

		bag, err := baggage.New(countryMember)
		Expect(err).ToNot(HaveOccurred())

		ctx = baggage.ContextWithBaggage(ctx, bag)

		// Act
		hash := HashByCountry(ctx, "some", "value")

		// Expected hash logic
		joined := strings.Join([]string{country, "some", "value"}, "|")
		sum := sha256.Sum256([]byte(joined))
		expected := hex.EncodeToString(sum[:]) // truncate to 128 bits

		// Assert
		Expect(string(hash)).To(Equal(expected))
	})

	It("should generate different hashes for different countries", func() {
		ctxCL := baggage.ContextWithBaggage(context.Background(), mustBaggage(
			sharedcontext.BaggageTenantCountry, "CL",
		))
		ctxAR := baggage.ContextWithBaggage(context.Background(), mustBaggage(
			sharedcontext.BaggageTenantCountry, "AR",
		))

		hashCL := HashByCountry(ctxCL, "x")
		hashAR := HashByCountry(ctxAR, "x")

		Expect(hashCL).ToNot(Equal(hashAR))
	})

	var _ = Describe("HashByTenant", func() {
		It("should generate a deterministic hash based on tenant and country and inputs", func() {
			// Arrange
			tenantID := "123"
			country := "CL"
			ctx := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, tenantID,
				sharedcontext.BaggageTenantCountry, country,
			))

			// Act
			hash := HashByTenant(ctx, "some", "value")

			// Expected hash logic
			orgKey := tenantID + "-" + country
			joined := strings.Join([]string{orgKey, "some", "value"}, "|")
			sum := sha256.Sum256([]byte(joined))
			expected := hex.EncodeToString(sum[:]) // truncate to 128 bits

			// Assert
			Expect(string(hash)).To(Equal(expected))
		})

		It("should generate different hashes for different tenant or country", func() {
			ctx1 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "1",
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctx2 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "2", // cambia tenant
				sharedcontext.BaggageTenantCountry, "CL",
			))
			ctx3 := baggage.ContextWithBaggage(context.Background(), mustBaggage(
				sharedcontext.BaggageTenantID, "1",
				sharedcontext.BaggageTenantCountry, "AR", // cambia country
			))

			Expect(HashByTenant(ctx1, "x")).ToNot(Equal(HashByTenant(ctx2, "x")))
			Expect(HashByTenant(ctx1, "x")).ToNot(Equal(HashByTenant(ctx3, "x")))
		})
	})

})

func mustBaggage(kv ...string) baggage.Baggage {
	Expect(len(kv)%2).To(Equal(0), "mustBaggage requires even number of arguments (key-value pairs)")

	members := make([]baggage.Member, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		m, err := baggage.NewMember(kv[i], kv[i+1])
		Expect(err).ToNot(HaveOccurred())
		members = append(members, m)
	}

	bag, err := baggage.New(members...)
	Expect(err).ToNot(HaveOccurred())
	return bag
}

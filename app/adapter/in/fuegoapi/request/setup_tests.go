package request

import (
	"context"
	"testing"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

func TestDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Suite")
}

func buildCtx(tenantID, country string) context.Context {
	ctx := context.Background()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID)
	cntry, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, country)
	bag, _ := baggage.New(tID, cntry)
	return baggage.ContextWithBaggage(ctx, bag)
}

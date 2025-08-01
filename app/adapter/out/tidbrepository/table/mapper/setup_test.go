package mapper

import (
	"context"
	"testing"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

func TestTableMapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Table Mapper Suite")
}

func buildCtx(tenantID, country string) context.Context {
	ctx := context.Background()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID)
	cntry, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, country)
	bag, _ := baggage.New(tID, cntry)
	return baggage.ContextWithBaggage(ctx, bag)
}

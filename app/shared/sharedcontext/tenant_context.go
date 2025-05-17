package sharedcontext

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

// WithOnlyTenantContext creates a new context with only tenant information from the source context,
// removing any other baggage information while maintaining the cancellation chain
func WithOnlyTenantContext(ctx context.Context) (context.Context, error) {
	bag := baggage.FromContext(ctx)
	tenantID := bag.Member(BaggageTenantID).Value()
	country := bag.Member(BaggageTenantCountry).Value()

	tenantMember, err := baggage.NewMember(BaggageTenantID, tenantID)
	if err != nil {
		return nil, err
	}

	countryMember, err := baggage.NewMember(BaggageTenantCountry, country)
	if err != nil {
		return nil, err
	}

	cleanBag, err := baggage.New(tenantMember, countryMember)
	if err != nil {
		return nil, err
	}

	return baggage.ContextWithBaggage(ctx, cleanBag), nil
}

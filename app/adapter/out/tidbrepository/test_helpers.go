package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/baggage"
)

// CreateTestTenant creates a new tenant for testing purposes
func CreateTestTenant(ctx context.Context, conn database.ConnectionFactory) (domain.Tenant, context.Context, error) {
	// Create a new tenant ID
	tenantID := uuid.New()

	// Create tenant entity
	tenant := domain.Tenant{
		ID:      tenantID,
		Country: countries.CL,
		Name:    "Test Tenant " + tenantID.String(),
	}

	// Create context with tenant information
	orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID.String())
	countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, tenant.Country.String())
	channelMember, _ := baggage.NewMember(sharedcontext.BaggageChannel, "test")
	bag, _ := baggage.New(orgIDMember, countryMember, channelMember)
	tenantCtx := baggage.ContextWithBaggage(ctx, bag)

	// Save tenant
	saveTenant := NewSaveTenant(conn, NewSaveFSMTransition(conn))
	savedTenant, err := saveTenant(tenantCtx, domain.Tenant{
		ID:      tenantID,
		Name:    "Test Tenant " + tenantID.String(),
		Country: tenant.Country,
	})
	if err != nil {
		return domain.Tenant{}, nil, err
	}

	return savedTenant, tenantCtx, nil
}

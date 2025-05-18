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
		Operator: domain.Operator{
			Contact: domain.Contact{
				PrimaryEmail: "test-" + tenantID.String() + "@example.com",
			},
		},
	}

	// Create context with tenant information
	orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID.String())
	countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, tenant.Country.String())
	bag, _ := baggage.New(orgIDMember, countryMember)
	tenantCtx := baggage.ContextWithBaggage(ctx, bag)

	// Create account first
	account := domain.Account{
		Email: tenant.Operator.Contact.PrimaryEmail,
	}
	upsertAccount := NewUpsertAccount(conn)
	err := upsertAccount(tenantCtx, account)
	if err != nil {
		return domain.Tenant{}, nil, err
	}

	// Save tenant
	saveTenant := NewSaveTenant(conn)
	savedTenant, err := saveTenant(tenantCtx, tenant)
	if err != nil {
		return domain.Tenant{}, nil, err
	}

	return savedTenant, tenantCtx, nil
}

package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapContactToTable(ctx context.Context, e domain.Contact) table.Contact {
	return table.Contact{
		FullName:       e.FullName,
		Email:          e.PrimaryEmail,
		Phone:          e.PrimaryPhone,
		Documents:      mapDocuments(e.Documents), // Serializar a JSON
		NationalID:     e.NationalID,
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		DocumentID:     string(e.DocID(ctx)),
	}
}

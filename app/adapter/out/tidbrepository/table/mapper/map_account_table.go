package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapAccountTable(e domain.Operator) table.Account {
	return table.Account{
		Email:    e.Contact.PrimaryEmail,
		IsActive: true,
	}
}

func MapAddressInfoTable(ctx context.Context, e domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		State:          e.State.String(),
		District:       e.District.String(),
		AddressLine1:   e.AddressLine1,
		DocumentID:     string(e.DocID(ctx)),
		Latitude:       e.Location[1],
		Longitude:      e.Location[0],
		ZipCode:        e.ZipCode,
		TimeZone:       e.TimeZone,
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		Province:       e.Province.String(),
	}
}

func mapDocuments(docs []domain.Document) table.JSONReference {
	// Crear un slice para mapear los documentos
	mapped := make(table.JSONReference, len(docs))

	// Iterar sobre los documentos y mapearlos
	for i, d := range docs {
		mapped[i] = table.Reference{
			Type:  d.Type,
			Value: d.Value,
		}
	}

	return mapped
}

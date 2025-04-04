package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Operator) table.Account {
	return table.Account{
		Email:    e.Contact.PrimaryEmail,
		IsActive: true,
	}
}

func MapAddressInfoTable(e domain.AddressInfo, organizationCountryID int64) table.AddressInfo {
	return table.AddressInfo{
		State:          e.State,
		District:       e.District,
		AddressLine1:   e.AddressLine1,
		DocumentID:     string(e.DocID()),
		Latitude:       e.Location[1],
		Longitude:      e.Location[0],
		ZipCode:        e.ZipCode,
		TimeZone:       e.TimeZone,
		OrganizationID: organizationCountryID,
		Province:       e.Province,
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

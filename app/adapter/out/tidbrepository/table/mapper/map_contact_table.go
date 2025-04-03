package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapContactToTable(e domain.Contact, organizationID int64) table.Contact {
	return table.Contact{
		ID:             e.ID,
		FullName:       e.FullName,
		Email:          e.PrimaryEmail,
		Phone:          e.PrimaryPhone,
		Documents:      mapDocuments(e.Documents), // Serializar a JSON
		NationalID:     e.NationalID,
		OrganizationID: organizationID,
		DocumentID:     string(e.DocID()),
	}
}

package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapContactToTable(e domain.Contact, organizationCountryID int64) table.Contact {
	return table.Contact{
		ID:                    e.ID,
		FullName:              e.FullName,
		Email:                 e.Email,
		Phone:                 e.Phone,
		Documents:             mapDocuments(e.Documents), // Serializar a JSON
		NationalID:            e.NationalID,
		OrganizationCountryID: organizationCountryID,
	}
}

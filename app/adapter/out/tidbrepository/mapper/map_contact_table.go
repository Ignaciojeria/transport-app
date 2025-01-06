package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapContactToTable(e domain.Contact) table.Contact {
	return table.Contact{
		FullName:  e.FullName,
		Email:     e.Email,
		Phone:     e.Phone,
		Documents: serializeToJSON(e.Documents), // Serializar a JSON
	}
}

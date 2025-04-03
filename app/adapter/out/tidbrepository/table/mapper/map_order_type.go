package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderType(ot domain.OrderType) table.OrderType {
	return table.OrderType{
		Type:           ot.Type,
		Description:    ot.Description,
		OrganizationID: ot.Organization.ID,
		DocumentID:     ot.DocID(),
	}
}

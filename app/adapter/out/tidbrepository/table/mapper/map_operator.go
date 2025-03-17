package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOperator(e domain.Operator) table.Account {
	var contactIDPtr *int64
	if e.Contact.ID != 0 {
		contactIDPtr = &e.Contact.ID
	}
	var originNodeInfoIDPtr *int64
	if e.OriginNode.ID != 0 {
		originNodeInfoIDPtr = &e.OriginNode.ID
	}
	var orgIDPtr *int64
	if e.Organization.ID != 0 {
		orgIDPtr = &e.Organization.ID
	}
	return table.Account{
		ID:               e.ID,
		ReferenceID:      e.ReferenceID,
		ContactID:        contactIDPtr,
		IsActive:         true,
		OriginNodeInfoID: originNodeInfoIDPtr,
		OrganizationID:   orgIDPtr,
		Type:             "operator",
	}
}

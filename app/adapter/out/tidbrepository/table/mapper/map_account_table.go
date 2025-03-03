package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Account) table.Account {
	var contactIDPtr *int64
	if e.Contact.ID != 0 {
		contactIDPtr = &e.Contact.ID
	}

	var originNodeInfoIDPtr *int64
	if e.Origin.ID != 0 {
		originNodeInfoIDPtr = &e.Origin.ID
	}
	return table.Account{
		ID:                    e.ID,
		ContactID:             contactIDPtr,
		IsActive:              true,
		OriginNodeInfoID:      originNodeInfoIDPtr,
		OrganizationCountryID: e.Organization.OrganizationCountryID,
	}
}

func MapNodeInfoTable(e domain.NodeInfo) table.NodeInfo {
	var contactID, addressInfoID, nodeTypeID *int64
	contactID = &e.Contact.ID
	addressInfoID = &e.AddressInfo.ID
	nodeTypeID = &e.NodeType.ID
	if e.Contact.ID == 0 {
		contactID = nil
	}
	if e.AddressInfo.ID == 0 {
		addressInfoID = nil
	}
	if e.NodeType.ID == 0 {
		nodeTypeID = nil
	}
	var nodeName *string = &e.Name
	if e.Name == "" {
		nodeName = nil
	}
	return table.NodeInfo{
		ID:                    e.ID,
		ReferenceID:           string(e.ReferenceID),
		Name:                  nodeName,
		NodeTypeID:            nodeTypeID,
		ContactID:             contactID,
		OrganizationCountryID: e.Organization.OrganizationCountryID,
		AddressID:             addressInfoID,
	}
}

func MapAddressInfoTable(e domain.AddressInfo, organizationCountryID int64) table.AddressInfo {
	return table.AddressInfo{
		ID:                    e.ID,
		State:                 e.State,
		Locality:              e.Locality,
		District:              e.District,
		AddressLine1:          e.AddressLine1,
		AddressLine2:          e.AddressLine2,
		AddressLine3:          e.AddressLine3,
		RawAddress:            e.RawAddress(),
		Latitude:              e.Location[1],
		Longitude:             e.Location[0],
		ZipCode:               e.ZipCode,
		TimeZone:              e.TimeZone,
		OrganizationCountryID: organizationCountryID,
		Province:              e.Province,
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

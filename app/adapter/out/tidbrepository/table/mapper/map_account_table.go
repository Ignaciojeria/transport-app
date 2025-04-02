package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Account) table.Account {
	return table.Account{
		ID:       e.ID,
		Email:    e.Email,
		IsActive: true,
	}
}

func MapNodeInfoTable(e domain.NodeInfo) table.NodeInfo {
	var contactID, addressInfoID, nodeTypeID *int64
	contactID = &e.Contact.ID

	nodeTypeID = &e.NodeType.ID
	if e.Contact.ID == 0 {
		contactID = nil
	}

	if e.NodeType.ID == 0 {
		nodeTypeID = nil
	}
	var nodeName *string = &e.Name
	if e.Name == "" {
		nodeName = nil
	}
	return table.NodeInfo{
		ID:             e.ID,
		ReferenceID:    string(e.ReferenceID),
		Name:           nodeName,
		NodeTypeID:     nodeTypeID,
		ContactID:      contactID,
		OrganizationID: e.Organization.ID,
		AddressID:      addressInfoID,
	}
}

func MapAddressInfoTable(e domain.AddressInfo, organizationCountryID int64) table.AddressInfo {
	return table.AddressInfo{
		State: e.State,
		//	Locality:       e.Locality,
		District:     e.District,
		AddressLine1: e.AddressLine1,
		//	AddressLine2:   e.AddressLine2,
		//	AddressLine3:   e.AddressLine3,
		//RawAddress:     e.FullAddress(),
		ReferenceID:    e.DocID(),
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

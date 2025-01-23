package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Account, originNodeInfoID int64, contactId int64, organizationCountryID int64) table.Account {
	return table.Account{
		ID:                    e.ID,
		ContactID:             contactId,
		IsActive:              true,
		OriginNodeInfoID:      originNodeInfoID,
		OrganizationCountryID: organizationCountryID,
	}
}

func MapNodeInfoTable(e domain.NodeInfo) table.NodeInfo {
	var operatorID, addressInfoID *int64
	if e.Operator.ID == 0 {
		operatorID = nil
	}
	if e.AddressInfo.ID == 0 {
		addressInfoID = nil
	}
	return table.NodeInfo{
		ID:                    e.ID,
		ReferenceID:           string(e.ReferenceID),
		Name:                  e.Name,
		Type:                  e.Type,
		OperatorID:            operatorID,
		OrganizationCountryID: e.Organization.OrganizationCountryID,
		AddressID:             addressInfoID,
	}
}

func MapAddressInfoTable(e domain.AddressInfo, organizationCountryID int64) table.AddressInfo {
	return table.AddressInfo{
		ID:                    e.ID,
		State:                 e.State,
		County:                e.County,
		District:              e.District,
		AddressLine1:          e.AddressLine1,
		AddressLine2:          e.AddressLine2,
		AddressLine3:          e.AddressLine3,
		RawAddress:            e.RawAddress(),
		Latitude:              e.Latitude,
		Longitude:             e.Longitude,
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

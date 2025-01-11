package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Account, originNodeInfoID int64, contactId int64, organizationCountryID int64) table.Account {
	return table.Account{
		ID:                    0,
		ContactID:             contactId,
		IsActive:              true,
		OriginNodeInfoID:      originNodeInfoID,
		OrganizationCountryID: organizationCountryID,
	}
}

func MapNodeInfoTable(e domain.NodeInfo, organizationCountryID int64, addressID int64) table.NodeInfo {
	return table.NodeInfo{
		ID:                    0,
		ReferenceID:           string(e.ReferenceID),
		Name:                  e.Name,
		Type:                  e.Type,
		OperatorID:            0,
		OrganizationCountryID: organizationCountryID,
		AddressID:             addressID,
		//NodeReferences: MapReferencesTable(e.References),
	}
}

func MapAddressInfoTable(e domain.AddressInfo, organizationCountryID int64) table.AddressInfo {
	return table.AddressInfo{
		ID:                    0, // ID inicializado en 0
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
	}
}

func MapReferencesTable(refs []domain.Reference) []table.NodeReference {
	var tableRefs []table.NodeReference
	for _, ref := range refs {
		tableRefs = append(tableRefs, table.NodeReference{
			ID:    0, // ID inicializado en 0
			Type:  ref.Type,
			Value: ref.Value,
		})
	}
	return tableRefs
}

func mapDocuments(docs []domain.Document) table.JSONDocuments {
	// Crear un slice para mapear los documentos
	mapped := make(table.JSONDocuments, len(docs))

	// Iterar sobre los documentos y mapearlos
	for i, d := range docs {
		mapped[i] = table.Document{
			Type:  d.Type,
			Value: d.Value,
		}
	}

	// Retornar los documentos mapeados como table.JSONDocuments
	return mapped
}

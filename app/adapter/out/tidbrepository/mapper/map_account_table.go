package mapper

import (
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Account) table.Account {
	return table.Account{
		ID: 0, // Establecer ID en 0
		Contact: table.Contact{
			ID:         0,
			FullName:   e.Contact.FullName,
			Email:      e.Contact.Email,
			Phone:      e.Contact.Phone,
			NationalID: e.Contact.NationalID,
		},
		IsActive: true, // Ejemplo de valor predeterminado para IsActive
		// Mapear el origen
		OriginNodeInfoID: 0, // ID inicializado en 0
		OriginNodeInfo:   MapNodeInfoTable(e.Origin.NodeInfo),
	}
}

func MapNodeInfoTable(e domain.NodeInfo) table.NodeInfo {
	return table.NodeInfo{
		ID:          0,
		ReferenceID: string(e.ReferenceID),
		Name:        e.Name,
		Type:        e.Type,
		OperatorID:  0,
		Operator: table.Operator{
			ID:   0,
			Type: e.Operator.Type,
			Contact: table.Contact{
				ID:         0,
				FullName:   e.Operator.Contact.FullName,
				Email:      e.Operator.Contact.Email,
				Phone:      e.Operator.Contact.Phone,
				NationalID: e.Operator.Contact.NationalID,
			},
		},
		NodeReferences: MapReferencesTable(e.References),
	}
}

func MapAddressInfoTable(e domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		ID:           0, // ID inicializado en 0
		State:        e.State,
		County:       e.County,
		District:     e.District,
		AddressLine1: e.AddressLine1,
		AddressLine2: e.AddressLine2,
		AddressLine3: e.AddressLine3,
		RawAddress:   e.RawAddress(),
		Latitude:     e.Latitude,
		Longitude:    e.Longitude,
		ZipCode:      e.ZipCode,
		TimeZone:     e.TimeZone,
	}
}

func MapReferencesTable(refs []domain.References) []table.NodeReferences {
	var tableRefs []table.NodeReferences
	for _, ref := range refs {
		tableRefs = append(tableRefs, table.NodeReferences{
			ID:    0, // ID inicializado en 0
			Type:  ref.Type,
			Value: ref.Value,
		})
	}
	return tableRefs
}

func serializeToJSON(data interface{}) []byte {
	// Serializar los métodos de contacto o documentos a JSON
	serialized, _ := json.Marshal(data)
	return serialized
}

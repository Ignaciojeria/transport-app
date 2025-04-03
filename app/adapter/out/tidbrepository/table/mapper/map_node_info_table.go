package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapNodeInfoTable(e domain.NodeInfo) table.NodeInfo {
	references := table.JSONReference{}
	for _, v := range e.References {
		references = append(references, table.Reference{
			Type:  v.Type,
			Value: v.Value,
		})
	}
	return table.NodeInfo{
		ID:             e.ID,
		ReferenceID:    string(e.ReferenceID),
		DocumentID:     string(e.DocID()),
		NodeTypeDoc:    string(e.NodeType.DocID()),
		Name:           e.Name,
		ContactDoc:     string(e.Contact.DocID()),
		OrganizationID: e.Organization.ID,
		AddressInfoDoc: string(e.AddressInfo.DocID()),
		AddressLine2:   e.AddressLine2,
		AddressLine3:   e.AddressLine3,
		NodeReferences: references,
	}
}

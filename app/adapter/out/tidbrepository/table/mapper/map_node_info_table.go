package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapNodeInfoTable(ctx context.Context, e domain.NodeInfo) table.NodeInfo {
	references := table.JSONReference{}
	for _, v := range e.References {
		references = append(references, table.Reference{
			Type:  v.Type,
			Value: v.Value,
		})
	}
	return table.NodeInfo{
		ReferenceID:    string(e.ReferenceID),
		DocumentID:     string(e.DocID(ctx)),
		NodeTypeDoc:    string(e.NodeType.DocID(ctx)),
		Name:           e.Name,
		ContactDoc:     string(e.AddressInfo.Contact.DocID(ctx)),
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		AddressInfoDoc: string(e.AddressInfo.DocID(ctx)),
		AddressLine2:   e.AddressLine2,
		AddressLine3:   e.AddressLine3,
		NodeReferences: references,
	}
}

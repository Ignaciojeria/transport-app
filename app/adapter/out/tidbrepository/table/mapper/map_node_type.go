package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapNodeType(nt domain.NodeType) table.NodeType {
	return table.NodeType{
		OrganizationCountryID: nt.Organization.OrganizationCountryID,
		ID:                    nt.ID,
		Name:                  nt.Value,
	}
}

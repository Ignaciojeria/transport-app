package views

import "transport-app/app/domain"

type SearchNodesView []struct {
	NodeName    string `gorm:"column:node_name"`
	ReferenceID string `gorm:"column:reference_id"`
}

func (sn SearchNodesView) Map() []domain.NodeInfo {
	var nodes []domain.NodeInfo
	for _, s := range sn {
		name := s.NodeName // Handle potential nil case if NodeName is a pointer in domain.NodeInfo
		nodes = append(nodes, domain.NodeInfo{
			Name:        name,
			ReferenceID: domain.ReferenceID(s.ReferenceID),
		})
	}
	return nodes
}

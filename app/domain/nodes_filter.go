package domain

type NodesFilter struct {
	Pagination      Pagination
	RequestedFields map[string]any
	ReferenceIds    []string
	Name            *string
	NodeType        *NodeTypeFilter
	References      []ReferenceFilter
}

type NodeTypeFilter struct {
	Value string
}

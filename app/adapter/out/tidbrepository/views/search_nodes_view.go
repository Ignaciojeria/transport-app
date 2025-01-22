package views

type SearchNodesView struct {
	NodeName    string `gorm:"column:node_name"`
	ReferenceID string `gorm:"column:reference_id"`
}

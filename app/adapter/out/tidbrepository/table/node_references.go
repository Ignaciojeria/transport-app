package table

type NodeReferences struct {
	ID         int64    `gorm:"primaryKey"`
	DocumentID string   `gorm:"type:char(64);uniqueIndex:idx_doc_node"`
	Type       string   `gorm:"not null"`
	Value      string   `gorm:"not null"`
	NodeDoc    string   `gorm:"type:char(64);uniqueIndex:idx_doc_node"`
	Node       NodeInfo `gorm:"-"`
}

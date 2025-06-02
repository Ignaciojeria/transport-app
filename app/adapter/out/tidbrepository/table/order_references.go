package table

type OrderReferences struct {
	ID         int64  `gorm:"primaryKey"`
	DocumentID string `gorm:"type:char(64);uniqueIndex:idx_doc_order"`
	Type       string `gorm:"not null"`
	Value      string `gorm:"not null"`
	OrderDoc   string `gorm:"type:char(64);uniqueIndex:idx_doc_order"`
	Order      Order  `gorm:"-"`
}

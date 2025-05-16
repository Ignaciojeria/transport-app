package table

type OrderReferences struct {
	ID         int64  `gorm:"primaryKey"`
	DocumentID string `gorm:"type:char(64);uniqueIndex"`
	Type       string `gorm:"not null"`
	Value      string `gorm:"not null"`
	OrderDoc   string `gorm:"type:char(64)"`
	Order      Order  `gorm:"-"`
}

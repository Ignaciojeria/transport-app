package table

type DeliveryUnitsLabels struct {
	ID              int64  `gorm:"primaryKey"`
	DocumentID      string `gorm:"type:char(64);uniqueIndex"`
	Type            string `gorm:"not null"`
	Value           string `gorm:"not null"`
	DeliveryUnitDoc string `gorm:"type:char(64)"`
	OrderDoc        string `gorm:"type:char(64)"`
}

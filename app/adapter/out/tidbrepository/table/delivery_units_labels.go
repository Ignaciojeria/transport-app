package table

type DeliveryUnitsLabels struct {
	ID              int64  `gorm:"primaryKey"`
	DocumentID      string `gorm:"type:char(64);uniqueIndex:idx_label_unique"`
	Type            string `gorm:"not null"`
	Value           string `gorm:"not null"`
	DeliveryUnitDoc string `gorm:"type:char(64);uniqueIndex:idx_label_unique"`
}

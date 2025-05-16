package table

type OrderDeliveryUnit struct {
	ID                   int64  `gorm:"primaryKey"`
	DeliveryUnitStatusID *int64 `gorm:"default null;index"`
	DeliveryUnitStatus   Status `gorm:"foreignKey:DeliveryUnitStatusID"`
	OrderDoc             string `gorm:"type:char(64);"`
	DeliveryUnitDoc      string `gorm:"type:char(64);"`
}

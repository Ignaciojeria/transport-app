package table

type OrderDeliveryUnit struct {
	ID              int64  `gorm:"primaryKey"`
	OrderDoc        string `gorm:"type:char(64);"`
	DeliveryUnitDoc string `gorm:"type:char(64);"`
}

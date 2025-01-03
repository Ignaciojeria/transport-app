package table

type Organization struct {
	ID      int64  `gorm:"primaryKey"`
	Country string `gorm:"not null;index:idx_country_email,unique"`
	Email   string `gorm:"type:varchar(255);not null;index:idx_country_email,unique"`
	Name    string `gorm:"type:varchar(255);not null;"`
}

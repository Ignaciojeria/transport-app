package table

type DeliveryUnitsSkills struct {
	ID              int64  `gorm:"primaryKey"`
	Skill           string `gorm:"type:varchar(100);uniqueIndex:idx_skill_unique"`
	DeliveryUnitDoc string `gorm:"type:char(64);uniqueIndex:idx_skill_unique"`
}

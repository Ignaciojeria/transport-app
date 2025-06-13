package table

type DeliveryUnitsSkills struct {
	ID              int64  `gorm:"primaryKey"`
	SkillDoc        string `gorm:"type:char(64);uniqueIndex:idx_skill_unique"`
	DeliveryUnitDoc string `gorm:"type:char(64);uniqueIndex:idx_skill_unique"`
}

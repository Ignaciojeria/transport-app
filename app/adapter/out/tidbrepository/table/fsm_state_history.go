package table

import (
	"gorm.io/gorm"
)

type FSMStateHistory struct {
	gorm.Model
	TraceID        string `gorm:"type:text;not null;index;uniqueIndex:fsm_unique_transition"`
	IdempotencyKey string `gorm:"type:text"`
	Workflow       string `gorm:"type:text;not null;uniqueIndex:fsm_unique_transition"`
	State          string `gorm:"type:text;not null;uniqueIndex:fsm_unique_transition"`
}

package table

import "time"

type OrderPlanningHistory struct {
	ID                  int64 `gorm:"primaryKey;autoIncrement"`
	PlanReferenceID     string
	RouteReferenceID    string
	OrderReferenceID    string
	OperatorReferenceID string
	JSONPlanLocation    JSONPlanLocation `gorm:"type:json"`
	VersionDate         time.Time
}

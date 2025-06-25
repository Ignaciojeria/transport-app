package optimization

import (
	"time"
)

type OptimizedFleet struct {
	PlanReferenceID string
	PlannedDate     time.Time
	Routes          []OptimizedRoute
	Unassigned      OptimizedUnassigned
}

type OptimizedRoute struct {
	VehiclePlate string
	Steps        []OptimizedStep
	Cost         int64
	Duration     int64
}

type OptimizedStep struct {
	Type       string // start, pickup, delivery, end
	VisitIndex int    // índice de visita original
	Location   Location
	Arrival    int64
	Orders     []Order
}

type OptimizedUnassigned struct {
	Orders   []Order
	Vehicles []Vehicle
	Origins  []NodeInfo
}

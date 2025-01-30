package domain

import "time"

type Plan struct {
	ID             int64
	ReferenceID    string
	Date           time.Time
	Routes         []Route
	PlanningStatus PlanningStatus
	PlanType       PlanType
}

type PlanType struct {
	ID    int64
	Value string
}

type PlanningStatus struct {
	ID    int64
	Value string
}

type Route struct {
	ID      int64
	Vehicle Vehicle
	Orders  []Order
}

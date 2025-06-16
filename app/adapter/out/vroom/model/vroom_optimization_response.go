package model

// VroomOptimizationResponse represents the complete response from VROOM API
type VroomOptimizationResponse struct {
	Code       int             `json:"code"`
	Error      string          `json:"error,omitempty"`
	Unassigned []UnassignedJob `json:"unassigned,omitempty"`
	Routes     []Route         `json:"routes,omitempty"`
}

// UnassignedJob represents jobs that couldn't be assigned to any vehicle
type UnassignedJob struct {
	ID       int        `json:"id"`
	Location [2]float64 `json:"location,omitempty"`
	Reason   string     `json:"reason"`
}

// Route represents a vehicle's route with its assigned jobs
type Route struct {
	Vehicle     int     `json:"vehicle"`
	Cost        int     `json:"cost"`
	Service     int     `json:"service"`
	Duration    int     `json:"duration"`
	WaitingTime int     `json:"waiting_time"`
	Priority    float64 `json:"priority"`
	Steps       []Step  `json:"steps"`
	Geometry    string  `json:"geometry,omitempty"`
}

// Step represents a single step in a route (pickup, delivery, or break)
type Step struct {
	Type           string         `json:"type"` // "start", "job", "pickup", "delivery", "break", "end"
	Arrival        int            `json:"arrival"`
	Duration       int            `json:"duration"`
	Service        int            `json:"service"`
	WaitingTime    int            `json:"waiting_time"`
	Job            int            `json:"job,omitempty"`
	Location       [2]float64     `json:"location,omitempty"`
	Load           []int          `json:"load,omitempty"`
	Distance       int            `json:"distance,omitempty"`
	Setup          int            `json:"setup,omitempty"`
	Shipment       int            `json:"shipment,omitempty"`
	Pickup         int            `json:"pickup,omitempty"`
	Delivery       int            `json:"delivery,omitempty"`
	Description    string         `json:"description,omitempty"`
	CustomUserData map[string]any `json:"custom_user_data,omitempty"`
}

package model

type VroomOptimizationRequest struct {
	Vehicles  []VroomVehicle  `json:"vehicles"`
	Jobs      []VroomJob      `json:"jobs,omitempty"`      // Para entregas simples
	Shipments []VroomShipment `json:"shipments,omitempty"` // Para entregas con pickup + delivery
	Options   *VroomOptions   `json:"options,omitempty"`
}

// --- Vehicles ---
type VroomVehicle struct {
	ID         int         `json:"id"`
	Start      *[2]float64 `json:"start,omitempty"`       // [lon, lat]
	End        *[2]float64 `json:"end,omitempty"`         // [lon, lat]
	Capacity   []int64     `json:"capacity,omitempty"`    // Ej: [peso, volumen]
	Skills     []int       `json:"skills,omitempty"`      // Habilidades codificadas como enteros
	TimeWindow []int       `json:"time_window,omitempty"` // [start, end] en segundos desde medianoche
}

// --- Jobs (entrega directa sin pickup) ---
type VroomJob struct {
	ID             int            `json:"id"`
	Location       [2]float64     `json:"location"`                   // [lon, lat]
	Service        int64          `json:"service,omitempty"`          // En segundos
	Amount         []int64        `json:"amount,omitempty"`           // Ej: [peso, volumen]
	Skills         []int          `json:"skills,omitempty"`           // Habilidades requeridas
	TimeWindows    [][]int        `json:"time_windows,omitempty"`     // [[start, end]]
	Priority       int            `json:"priority,omitempty"`         // 0-100
	CustomUserData map[string]any `json:"custom_user_data,omitempty"` // Metadata opcional
}

// --- Shipments (con pickup + delivery) ---
type VroomShipment struct {
	ID             int            `json:"id"`
	Pickup         VroomStep      `json:"pickup"`
	Delivery       VroomStep      `json:"delivery"`
	Amount         []int64        `json:"amount,omitempty"`
	Skills         []int          `json:"skills,omitempty"`
	TimeWindows    [][]int        `json:"time_windows,omitempty"`
	Service        int64          `json:"service,omitempty"`
	CustomUserData map[string]any `json:"custom_user_data,omitempty"`
}

type VroomStep struct {
	ID          int        `json:"id"`                     // Unique identifier for the step
	Location    [2]float64 `json:"location"`               // [lon, lat]
	TimeWindows [][]int    `json:"time_windows,omitempty"` // [[start, end]]
}

// --- Options ---
type VroomOptions struct {
	G                bool `json:"g,omitempty"`                      // Retornar geometría
	Steps            bool `json:"steps,omitempty"`                  // Incluir pasos
	Overview         bool `json:"overview,omitempty"`               // Ruta simplificada
	Polygons         bool `json:"polygons,omitempty"`               // Polígonos de zonas
	MinimizeVehicles bool `json:"minimize_vehicle_usage,omitempty"` // Usar menos vehículos
}

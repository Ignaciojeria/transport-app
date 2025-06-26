package optimization

// Coordinates representa coordenadas geográficas
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// NodeInfo representa información de un nodo
type NodeInfo struct {
	ReferenceID string
}

// Location representa una ubicación con coordenadas e información del nodo
type Location struct {
	Latitude  float64
	Longitude float64
	NodeInfo  NodeInfo
}

// TimeWindow representa una ventana de tiempo
type TimeWindow struct {
	Start string
	End   string
}

// Capacity representa la capacidad de un vehículo
type Capacity struct {
	Insurance             int64
	Volume                int64
	Weight                int64
	DeliveryUnitsQuantity int64
}

// Contact representa información de contacto
type Contact struct {
	Email      string
	Phone      string
	NationalID string
	FullName   string
}

// PoliticalArea representa información política/geográfica
type PoliticalArea struct {
	Code     string
	District string
	Province string
	State    string
}

// VisitLocation representa una ubicación de visita con toda su información
type VisitLocation struct {
	Coordinates   Coordinates
	ServiceTime   int64
	Contact       Contact
	Skills        []string
	TimeWindow    TimeWindow
	NodeInfo      NodeInfo
	PoliticalArea PoliticalArea
}

// Item representa un artículo
type Item struct {
	Sku string
}

// DeliveryUnit representa una unidad de entrega
type DeliveryUnit struct {
	Items     []Item
	Insurance int64
	Volume    int64
	Weight    int64
	Lpn       string
}

// Order representa una orden
type Order struct {
	DeliveryUnits  []DeliveryUnit
	ReferenceID    string
	SequenceNumber *int
}

// Vehicle representa un vehículo
type Vehicle struct {
	Plate         string
	StartLocation Location
	EndLocation   Location
	Skills        []string
	TimeWindow    TimeWindow
	Capacity      Capacity
	PoliticalArea PoliticalArea
}

// Visit representa una visita con pickup, delivery y órdenes
type Visit struct {
	Pickup   VisitLocation
	Delivery VisitLocation
	Orders   []Order
}

// Optimization representa la estructura principal de optimización
type FleetOptimization struct {
	PlanReferenceID string
	Vehicles        []Vehicle
	Visits          []Visit
}

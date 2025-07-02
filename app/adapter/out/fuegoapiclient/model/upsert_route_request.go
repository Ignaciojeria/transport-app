package model

type UpsertRouteRequest struct {
	ReferenceID string   `json:"referenceID" example:"ROUTE-001"`
	Plan        Plan     `json:"plan"`
	Vehicle     Vehicle  `json:"vehicle"`
	Geometry    Geometry `json:"geometry"`
	Visits      []Visit  `json:"visits"`
	CreatedAt   string   `json:"createdAt" example:"2025-01-15T10:30:00Z"`
}

type Plan struct {
	ReferenceID string `json:"referenceID" example:"PLAN-001"`
}

type Vehicle struct {
	Plate string `json:"plate" example:"ABCD12"`
}

type Geometry struct {
	Encoding string `json:"encoding" example:"polyline"`
	Type     string `json:"type" example:"linestring"`
	Value    string `json:"value" example:"_p~iF~ps|U_ulLnnqC_mqNvxq@"`
}

type Visit struct {
	Type                 string      `json:"type" example:"delivery"`
	AddressInfo          AddressInfo `json:"addressInfo"`
	NodeInfo             NodeInfo    `json:"nodeInfo"`
	DeliveryInstructions string      `json:"deliveryInstructions" example:"Entregar en recepción"`
	SequenceNumber       int         `json:"sequenceNumber" example:"1"`
	Orders               []Order     `json:"orders"`
}

type AddressInfo struct {
	AddressLine1  string        `json:"addressLine1" example:"Av. Providencia 1234"`
	AddressLine2  string        `json:"addressLine2" example:"Oficina 567"`
	Contact       Contact       `json:"contact"`
	Coordinates   Coordinates   `json:"coordinates"`
	PoliticalArea PoliticalArea `json:"politicalArea"`
	ZipCode       string        `json:"zipCode" example:"7500000"`
}

type Contact struct {
	Email      string `json:"email" example:"cliente@ejemplo.com"`
	FullName   string `json:"fullName" example:"Juan Pérez"`
	NationalID string `json:"nationalID" example:"12345678-9"`
	Phone      string `json:"phone" example:"+56912345678"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude" example:"-33.5147889"`
	Longitude float64 `json:"longitude" example:"-70.6130425"`
}

type PoliticalArea struct {
	Code     string `json:"code" example:"cl-rm-providencia"`
	District string `json:"district" example:"providencia"`
	Province string `json:"province" example:"santiago"`
	State    string `json:"state" example:"region metropolitana de santiago"`
}

type NodeInfo struct {
	ReferenceID string `json:"referenceID" example:"NODE-001"`
}

type Order struct {
	DeliveryUnits []DeliveryUnit `json:"deliveryUnits"`
	ReferenceID   string         `json:"referenceID" example:"ORDER-001"`
}

type DeliveryUnit struct {
	Items []Item `json:"items"`
	Lpn   string `json:"lpn" example:"LPN-789012"`
}

type Item struct {
	Sku string `json:"sku" example:"SKU-123456"`
}

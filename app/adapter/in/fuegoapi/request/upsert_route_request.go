package request

import (
	"errors"
	"transport-app/app/domain"
)

type UpsertRouteRequest struct {
	DocumentID  string `json:"documentID"`
	ReferenceID string `json:"referenceID,omitempty" example:"ROUTE-001"`
	CreatedAt   string `json:"createdAt,omitempty" example:"2025-01-15T10:30:00Z"`

	// Información del plan de optimización original
	PlanReferenceID string `json:"planReferenceID,omitempty" example:"PLAN-001"`

	// Información del vehículo
	Vehicle UpsertRouteVehicle `json:"vehicle,omitempty"`

	// Geometría de la ruta
	Geometry UpsertRouteGeometry `json:"geometry,omitempty"`

	// Visitas de la ruta (pickup y delivery)
	Visits []UpsertRouteVisit `json:"visits,omitempty"`
}

type UpsertRouteVehicle struct {
	Plate         string                     `json:"plate,omitempty" example:"ABCD12"`
	StartLocation UpsertRouteVehicleLocation `json:"startLocation,omitempty"`
	EndLocation   UpsertRouteVehicleLocation `json:"endLocation,omitempty"`
	Skills        []string                   `json:"skills,omitempty"`
	TimeWindow    UpsertRouteTimeWindow      `json:"timeWindow,omitempty"`
	Capacity      UpsertRouteVehicleCapacity `json:"capacity,omitempty"`
}

type UpsertRouteVehicleLocation struct {
	AddressInfo UpsertRouteAddressInfo `json:"addressInfo,omitempty"`
	NodeInfo    UpsertRouteNodeInfo    `json:"nodeInfo,omitempty"`
}

type UpsertRouteAddressInfo struct {
	AddressLine1  string                   `json:"addressLine1,omitempty"`
	AddressLine2  string                   `json:"addressLine2,omitempty"`
	Contact       UpsertRouteContact       `json:"contact,omitempty"`
	Coordinates   UpsertRouteCoordinates   `json:"coordinates,omitempty"`
	PoliticalArea UpsertRoutePoliticalArea `json:"politicalArea,omitempty"`
	ZipCode       string                   `json:"zipCode,omitempty"`
}

type UpsertRouteContact struct {
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	NationalID string `json:"nationalID,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}

type UpsertRouteCoordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type UpsertRoutePoliticalArea struct {
	Code            string `json:"code,omitempty"`
	AdminAreaLevel1 string `json:"adminAreaLevel1,omitempty"`
	AdminAreaLevel2 string `json:"adminAreaLevel2,omitempty"`
	AdminAreaLevel3 string `json:"adminAreaLevel3,omitempty"`
	AdminAreaLevel4 string `json:"adminAreaLevel4,omitempty"`
}

type UpsertRouteVehicleCapacity struct {
	Volume                int64 `json:"volume,omitempty" example:"1000"`
	Weight                int64 `json:"weight,omitempty" example:"1000"`
	Insurance             int64 `json:"insurance,omitempty" example:"10000"`
	DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity,omitempty"`
}

type UpsertRouteTimeWindow struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type UpsertRouteGeometry struct {
	Encoding string `json:"encoding,omitempty" example:"polyline"`
	Type     string `json:"type,omitempty" example:"linestring"`
	Value    string `json:"value,omitempty" example:"_p~iF~ps|U_ulLnnqC_mqNvxq@"`
}

type UpsertRouteVisit struct {
	Type                 string                 `json:"type,omitempty" example:"delivery"`
	Instructions         string                 `json:"instructions,omitempty"`
	AddressInfo          UpsertRouteAddressInfo `json:"addressInfo,omitempty"`
	NodeInfo             UpsertRouteNodeInfo    `json:"nodeInfo,omitempty"`
	DeliveryInstructions string                 `json:"deliveryInstructions,omitempty" example:"Entregar en recepción"`
	SequenceNumber       int                    `json:"sequenceNumber,omitempty" example:"1"`
	ServiceTime          int64                  `json:"serviceTime,omitempty"`
	TimeWindow           UpsertRouteTimeWindow  `json:"timeWindow,omitempty"`
	Orders               []UpsertRouteOrder     `json:"orders,omitempty"`
	// Motivo de no asignación (solo para rutas UNASSIGNED)
	UnassignedReason string `json:"unassignedReason,omitempty" example:"Vehicle capacity exceeded"`
}

type UpsertRouteNodeInfo struct {
	ReferenceID string `json:"referenceID,omitempty" example:"NODE-001"`
}

type UpsertRouteOrder struct {
	ReferenceID   string                    `json:"referenceID,omitempty" example:"ORDER-001"`
	DocumentID    string                    `json:"documentID"`
	DeliveryUnits []UpsertRouteDeliveryUnit `json:"deliveryUnits,omitempty"`
}

type UpsertRouteDeliveryUnit struct {
	DocumentID string            `json:"documentID"`
	Items      []UpsertRouteItem `json:"items,omitempty"`
	Volume     int64             `json:"volume,omitempty" example:"1000"`
	Weight     int64             `json:"weight,omitempty" example:"1000"`
	Price      int64             `json:"price,omitempty" example:"10000"`
	Lpn        string            `json:"lpn,omitempty" example:"LPN-789012"`
	Skills     []string          `json:"skills,omitempty"`
}

type UpsertRouteItem struct {
	Sku         string `json:"sku,omitempty" example:"SKU-123456"`
	Description string `json:"description,omitempty" example:"Descripción del artículo"`
	Quantity    int    `json:"quantity,omitempty" example:"1"`
}

// FlattenForExcel convierte la respuesta anidada en una estructura plana para Excel
// incluyendo solo fleet_id, visit_id y sequence
func (r *UpsertRouteRequest) FlattenForExcel() []map[string]interface{} {
	var rows []map[string]interface{}

	// Procesar visitas para generar filas con fleet_id, visit_id y sequence
	for _, visit := range r.Visits {
		for _, order := range visit.Orders {
			visitRow := make(map[string]interface{})

			// fleet_id (usando la patente del vehículo)
			visitRow["fleet_id"] = r.Vehicle.Plate

			// visit_id (usando el ReferenceID de la orden)
			visitRow["visit_id"] = order.ReferenceID

			// sequence
			visitRow["sequence"] = visit.SequenceNumber

			rows = append(rows, visitRow)
		}
	}

	return rows
}

// Map convierte la request a entidades del dominio
func (r *UpsertRouteRequest) Map() (domain.Route, error) {
	if r == nil {
		return domain.Route{}, errors.New("request cannot be nil")
	}

	// Mapear vehículo (solo campos básicos que coinciden)
	vehicle := domain.Vehicle{
		Plate: r.Vehicle.Plate,
	}

	// Mapear geometría
	geometry := domain.RouteGeometry{
		Encoding: r.Geometry.Encoding,
		Type:     r.Geometry.Type,
		Value:    r.Geometry.Value,
	}

	// Mapear órdenes desde las visitas
	var orders []domain.Order
	for _, v := range r.Visits {
		for _, o := range v.Orders {
			order := domain.Order{
				ReferenceID: domain.ReferenceID(o.ReferenceID),
			}
			orders = append(orders, order)
		}
	}

	// Crear ruta con la estructura del dominio
	route := domain.Route{
		ReferenceID: r.ReferenceID,
		Vehicle:     vehicle,
		Orders:      orders,
		Geometry:    geometry,
		TimeWindow: domain.TimeWindow{
			Start: r.Vehicle.TimeWindow.Start,
			End:   r.Vehicle.TimeWindow.End,
		},
	}

	return route, nil
}

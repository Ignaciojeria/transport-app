package request

import "transport-app/app/shared/projection/deliveryunits"

type UpsertRouteRequest struct {
	// Información básica de la ruta
	VisitKeyMapping   map[string]string
	VehicleKeyMapping map[string]string

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
	DeliveryUnits []UpsertRouteDeliveryUnit `json:"deliveryUnits,omitempty"`
}

type UpsertRouteDeliveryUnit struct {
	Items  []UpsertRouteItem `json:"items,omitempty"`
	Volume int64             `json:"volume,omitempty" example:"1000"`
	Weight int64             `json:"weight,omitempty" example:"1000"`
	Price  int64             `json:"price,omitempty" example:"10000"`
	Lpn    string            `json:"lpn,omitempty" example:"LPN-789012"`
	Skills []string          `json:"skills,omitempty"`
}

type UpsertRouteItem struct {
	Sku         string `json:"sku,omitempty" example:"SKU-123456"`
	Description string `json:"description,omitempty" example:"Descripción del artículo"`
	Quantity    int64  `json:"quantity,omitempty" example:"1"`
}

// FlattenForExcel convierte la respuesta anidada en una estructura plana para Excel
// usando los mappings para convertir las constantes del dominio a nombres de columnas
func (r *UpsertRouteRequest) FlattenForExcel() []map[string]interface{} {
	proj := deliveryunits.NewProjection()
	var rows []map[string]interface{}

	// Procesar visitas usando round robin (cada delivery unit genera una fila)
	for _, visit := range r.Visits {
		for _, order := range visit.Orders {
			for _, deliveryUnit := range order.DeliveryUnits {
				visitRow := make(map[string]interface{})

				// Mapear campos de la visita usando VisitKeyMapping
				for domainKey, excelColumn := range r.VisitKeyMapping {
					if excelColumn != "" {
						switch domainKey {
						case proj.ReferenceID().String():
							visitRow[excelColumn] = order.ReferenceID
						case proj.DestinationContactFullName().String():
							visitRow[excelColumn] = visit.AddressInfo.Contact.FullName
						case proj.DestinationContactPhone().String():
							visitRow[excelColumn] = visit.AddressInfo.Contact.Phone
						case proj.DestinationAddressLine1().String():
							visitRow[excelColumn] = visit.AddressInfo.AddressLine1
						case proj.DestinationCoordinatesLatitude().String():
							visitRow[excelColumn] = visit.AddressInfo.Coordinates.Latitude
						case proj.DestinationCoordinatesLongitude().String():
							visitRow[excelColumn] = visit.AddressInfo.Coordinates.Longitude
						case proj.DeliveryUnitVolume().String():
							visitRow[excelColumn] = deliveryUnit.Volume
						case proj.DeliveryUnitWeight().String():
							visitRow[excelColumn] = deliveryUnit.Weight
						case "price":
							visitRow[excelColumn] = 0.0 // Valor por defecto
						}
					}
				}

				// Agregar la patente/ID del vehículo a cada fila de visita
				if r.Vehicle.Plate != "" {
					for domainKey, excelColumn := range r.VehicleKeyMapping {
						if excelColumn != "" && (domainKey == "id" || domainKey == "plate") {
							visitRow[excelColumn] = r.Vehicle.Plate
							break // Solo necesitamos agregar la patente una vez
						}
					}
				}

				if len(visitRow) > 0 {
					rows = append(rows, visitRow)
				}
			}
		}
	}

	return rows
}

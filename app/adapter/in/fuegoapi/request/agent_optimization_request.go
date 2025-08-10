package request

import (
	"strconv"
	"strings"
	"transport-app/app/adapter/out/agents/model"
	"transport-app/app/shared/projection/deliveryunits"

	"github.com/google/uuid"
)

type AgentOptimizationRequest struct {
	Fleet  []map[string]interface{} `json:"fleet"`
	Visits []map[string]interface{} `json:"visits"`
}

func (r *AgentOptimizationRequest) ToOptimizeFleetRequest() OptimizeFleetRequest {
	// Mapear vehículos
	vehicles := make([]OptimizeFleetVehicle, len(r.Fleet))
	for i, vehicle := range r.Fleet {
		vehicles[i] = r.mapVehicle(vehicle)
	}

	// Mapear visitas
	visits := make([]OptimizeFleetVisit, len(r.Visits))
	for i, visit := range r.Visits {
		visits[i] = r.mapVisit(visit)
	}

	return OptimizeFleetRequest{
		PlanReferenceID: uuid.New().String(),
		Vehicles:        vehicles,
		Visits:          visits,
	}
}

func (r *AgentOptimizationRequest) mapVehicle(vehicle map[string]interface{}) OptimizeFleetVehicle {
	// Usar las constantes de VehicleFieldMappingSchema
	return OptimizeFleetVehicle{
		Plate: r.getStringValue(vehicle, model.VehicleKeyID),
		StartLocation: OptimizeFleetVehicleLocation{
			AddressInfo: OptimizeFleetAddressInfo{
				Coordinates: OptimizeFleetCoordinates{
					Latitude:  r.getFloatValue(vehicle, model.VehicleKeyStartLocationLatitude),
					Longitude: r.getFloatValue(vehicle, model.VehicleKeyStartLocationLongitude),
				},
			},
		},
		EndLocation: OptimizeFleetVehicleLocation{
			AddressInfo: OptimizeFleetAddressInfo{
				Coordinates: OptimizeFleetCoordinates{
					Latitude:  r.getFloatValue(vehicle, model.VehicleKeyEndLocationLatitude),
					Longitude: r.getFloatValue(vehicle, model.VehicleKeyEndLocationLongitude),
				},
			},
		},
		Capacity: OptimizeFleetVehicleCapacity{
			Weight:    int64(r.getFloatValue(vehicle, model.VehicleKeyWeight)),
			Volume:    int64(r.getFloatValue(vehicle, model.VehicleKeyVolume)),
			Insurance: int64(r.getFloatValue(vehicle, model.VehicleKeyInsurance)),
		},
	}
}

func (r *AgentOptimizationRequest) mapVisit(visit map[string]interface{}) OptimizeFleetVisit {
	projection := deliveryunits.NewProjection()

	return OptimizeFleetVisit{
		Delivery: OptimizeFleetVisitLocation{
			AddressInfo: OptimizeFleetAddressInfo{
				AddressLine1: r.getStringValue(visit, projection.DestinationAddressLine1().String()),
				Contact: OptimizeFleetContact{
					FullName: r.getStringValue(visit, projection.DestinationContactFullName().String()),
					Phone:    r.getStringValue(visit, projection.DestinationContactPhone().String()),
				},
				Coordinates: OptimizeFleetCoordinates{
					Latitude:  r.getFloatValue(visit, projection.DestinationCoordinatesLatitude().String()),
					Longitude: r.getFloatValue(visit, projection.DestinationCoordinatesLongitude().String()),
				},
			},
		},
		Orders: []OptimizeFleetOrder{
			{
				ReferenceID: r.getStringValue(visit, projection.ReferenceID().String()),
				DeliveryUnits: []OptimizeFleetDeliveryUnit{
					{
						Volume: int64(r.getFloatValue(visit, projection.DeliveryUnitVolume().String())),
						Weight: int64(r.getFloatValue(visit, projection.DeliveryUnitWeight().String())),
					},
				},
			},
		},
	}
}

// Helper methods para extraer valores de manera segura
func (r *AgentOptimizationRequest) getStringValue(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists && value != nil {
		if str, ok := value.(string); ok {
			return strings.TrimSpace(str)
		}
	}
	return ""
}

func (r *AgentOptimizationRequest) getFloatValue(data map[string]interface{}, key string) float64 {
	if value, exists := data[key]; exists && value != nil {
		switch v := value.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case string:
			if str := strings.TrimSpace(v); str != "" {
				if f, err := strconv.ParseFloat(str, 64); err == nil {
					return f
				}
			}
		}
	}
	return 0.0
}

// IterateVisitsInBatches itera sobre las visitas en lotes del tamaño especificado
// y ejecuta la función callback para cada lote
func (r *AgentOptimizationRequest) IterateVisitsInBatches(batchSize int, callback func([]map[string]interface{}) error) error {
	if batchSize <= 0 {
		batchSize = 100 // tamaño por defecto
	}

	totalVisits := len(r.Visits)
	if totalVisits == 0 {
		return nil
	}

	for i := 0; i < totalVisits; i += batchSize {
		end := i + batchSize
		if end > totalVisits {
			end = totalVisits
		}

		batch := r.Visits[i:end]
		if err := callback(batch); err != nil {
			return err
		}
	}

	return nil
}

// GetVisitsBatch devuelve un lote específico de visitas
func (r *AgentOptimizationRequest) GetVisitsBatch(startIndex, batchSize int) []map[string]interface{} {
	if startIndex < 0 || batchSize <= 0 {
		return nil
	}

	totalVisits := len(r.Visits)
	if startIndex >= totalVisits {
		return nil
	}

	endIndex := startIndex + batchSize
	if endIndex > totalVisits {
		endIndex = totalVisits
	}

	return r.Visits[startIndex:endIndex]
}

// GetTotalVisitsCount devuelve el número total de visitas
func (r *AgentOptimizationRequest) GetTotalVisitsCount() int {
	return len(r.Visits)
}

// GetTotalBatches calcula el número total de lotes necesarios para procesar todas las visitas
func (r *AgentOptimizationRequest) GetTotalBatches(batchSize int) int {
	if batchSize <= 0 {
		batchSize = 100
	}

	totalVisits := len(r.Visits)
	if totalVisits == 0 {
		return 0
	}

	return (totalVisits + batchSize - 1) / batchSize
}

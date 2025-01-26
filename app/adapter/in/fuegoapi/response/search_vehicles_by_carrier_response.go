package response

import "transport-app/app/domain"

type SearchVehiclesByCarrierResponse struct {
	ReferenceID         string  `json:"referenceID"`
	Plate               string  `json:"plate"`
	Category            string  `json:"category"`
	MaxPackageQuantity  int     `json:"maxPackageQuantity"`
	LoadInsuranceAmount float64 `json:"loadInsuranceAmount"`
	LoadCapacity        float64 `json:"loadCapacity"`
}

func MapSearchVehiclesByCarrierResponse(vehicles []domain.Vehicle) []SearchVehiclesByCarrierResponse {
	response := make([]SearchVehiclesByCarrierResponse, len(vehicles))
	for i, vehicle := range vehicles {
		response[i] = SearchVehiclesByCarrierResponse{
			ReferenceID:         vehicle.ReferenceID,
			Plate:               vehicle.Plate,
			Category:            vehicle.VehicleCategory.Type,
			MaxPackageQuantity:  vehicle.VehicleCategory.MaxPackagesQuantity,
			LoadInsuranceAmount: vehicle.Insurance.MaxInsuranceCoverage.Amount,
			LoadCapacity:        0,
		}
	}
	return response
}

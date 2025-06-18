package mapper

import (
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/vroom/model"
)

func MapOptimizationRequest(req request.OptimizationRequest) (model.VroomOptimizationRequest, error) {
	var vehicles []model.VroomVehicle
	for i, v := range req.Vehicles {
		vehicles = append(vehicles, model.VroomVehicle{
			ID: i + 1,
			Start: &[2]float64{
				v.StartLocation.Longitude,
				v.StartLocation.Latitude,
			},
			End: &[2]float64{
				v.EndLocation.Longitude,
				v.EndLocation.Latitude,
			},
			Capacity: []int64{
				v.Capacity.Weight,
				v.Capacity.DeliveryUnitsQuantity,
				v.Capacity.Insurance,
			},
			Skills:     mapSkills(v.Skills),
			TimeWindow: parseTimeRange(v.TimeWindow.Start, v.TimeWindow.End),
		})
	}

	var shipments []model.VroomShipment
	for i, visit := range req.Visits {
		pickup := model.VroomStep{
			ID: (i+1)*2 - 1, // Odd numbers for pickup
			Location: [2]float64{
				visit.PickupLocation.Longitude,
				visit.PickupLocation.Latitude,
			},
			TimeWindows: [][]int{parseTimeRange(visit.TimeWindow.Start, visit.TimeWindow.End)},
		}
		delivery := model.VroomStep{
			ID: (i + 1) * 2, // Even numbers for delivery
			Location: [2]float64{
				visit.DispatchLocation.Longitude,
				visit.DispatchLocation.Latitude,
			},
			TimeWindows: [][]int{parseTimeRange(visit.TimeWindow.Start, visit.TimeWindow.End)},
		}

		shipments = append(shipments, model.VroomShipment{
			ID:       i + 1,
			Pickup:   pickup,
			Delivery: delivery,
			Amount: []int64{
				visit.CapacityUsage.Weight,
				visit.CapacityUsage.DeliveryUnitsQuantity,
				visit.CapacityUsage.Insurance,
			},
			Skills:  mapSkills(visit.Skills),
			Service: visit.ServiceTime,
			CustomUserData: map[string]any{
				"orders": visit.Orders,
			},
		})
	}

	return model.VroomOptimizationRequest{
		Vehicles:  vehicles,
		Shipments: shipments,
		Options: &model.VroomOptions{
			G:                true,
			Steps:            true,
			Overview:         true,
			MinimizeVehicles: true,
		},
	}, nil
}

func parseTimeRange(start, end string) []int {
	// Convierte "08:00" a segundos desde medianoche
	// Retorna [inicio, fin] en segundos
	return []int{
		toSeconds(start),
		toSeconds(end),
	}
}

func toSeconds(timeStr string) int {
	// espera formato "HH:MM"
	var h, m int
	fmt.Sscanf(timeStr, "%02d:%02d", &h, &m)
	return h*3600 + m*60
}

func mapSkills(skills []string) []int {
	var mapped []int
	for _, skill := range skills {
		// Esto es un stub, deberías tener un map[string]int si necesitas traducirlos
		mapped = append(mapped, hashStringToInt(skill))
	}
	return mapped
}

func hashStringToInt(s string) int {
	// Ejemplo básico, podrías usar un hash real
	h := 0
	for i := 0; i < len(s); i++ {
		h += int(s[i])
	}
	return h
}

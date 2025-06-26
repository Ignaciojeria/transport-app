package request

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateMassiveTestData genera datos de prueba masivos para testing de rendimiento
func GenerateMassiveTestData() OptimizeFleetRequest {
	// Semilla para generar datos consistentes
	rand.Seed(time.Now().UnixNano())

	// Coordenadas base de La Florida, Chile
	baseLat := -33.5225
	baseLon := -70.575

	// Generar 10 vehículos con ubicaciones en Florida, Chile
	vehicles := generateVehicles(baseLat, baseLon)

	// Generar 10 visitas (5 por vehículo) con órdenes
	visits := generateVisits(baseLat, baseLon)

	return OptimizeFleetRequest{
		PlanReferenceID: "MASSIVE_TEST_PLAN",
		Vehicles:        vehicles,
		Visits:          visits,
	}
}

// generateVehicles genera 10 vehículos con capacidades variadas y ubicaciones en Florida, Chile
func generateVehicles(baseLat, baseLon float64) []struct {
	Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
	StartLocation struct {
		Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
		Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
		NodeInfo  struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"startLocation"`
	EndLocation struct {
		Latitude  float64 `json:"latitude" example:"-33.45" description:"Ending point latitude"`
		Longitude float64 `json:"longitude" example:"-70.66" description:"Ending point longitude"`
		NodeInfo  struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"endLocation"`
	Skills     []string `json:"skills" description:"Vehicle capabilities such as size or equipment requirements. eg: XL, heavy, etc"`
	TimeWindow struct {
		Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
		End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
	} `json:"timeWindow"`
	Capacity struct {
		Insurance             int64 `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
		Volume                int64 `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
		Weight                int64 `json:"weight" example:"1000" description:"Maximum weight in grams"`
		DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
	} `json:"capacity"`
} {
	var vehicles []struct {
		Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
		StartLocation struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
			NodeInfo  struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"startLocation"`
		EndLocation struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Ending point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Ending point longitude"`
			NodeInfo  struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"endLocation"`
		Skills     []string `json:"skills" description:"Vehicle capabilities such as size or equipment requirements. eg: XL, heavy, etc"`
		TimeWindow struct {
			Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
			End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
		} `json:"timeWindow"`
		Capacity struct {
			Insurance             int64 `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
			Volume                int64 `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
			Weight                int64 `json:"weight" example:"1000" description:"Maximum weight in grams"`
			DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
		} `json:"capacity"`
	}

	// Capacidades para los 10 vehículos (iguales para balancear)
	capacities := []struct {
		deliveryUnits int64
		insurance     int64
		volume        int64
		weight        int64
	}{}
	for i := 0; i < 10; i++ {
		capacities = append(capacities, struct {
			deliveryUnits int64
			insurance     int64
			volume        int64
			weight        int64
		}{30, 999999999, 999999999, 999999999})
	}

	// Puntos de partida diferentes para cada vehículo en Florida, Chile
	startLocations := []struct {
		lat float64
		lon float64
		ref string
	}{}
	for i := 0; i < 10; i++ {
		latOffset := (rand.Float64() - 0.5) * 0.05 // ±0.025 grados ≈ ±2.5km
		lonOffset := (rand.Float64() - 0.5) * 0.05 // ±0.025 grados ≈ ±2.5km
		startLocations = append(startLocations, struct {
			lat float64
			lon float64
			ref string
		}{baseLat + latOffset, baseLon + lonOffset, fmt.Sprintf("depot-la-florida-%d", i+1)})
	}

	for i := 0; i < 10; i++ {
		plate := fmt.Sprintf("SERV-%d", 80+i)
		vehicle := struct {
			Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
			StartLocation struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
				NodeInfo  struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			} `json:"startLocation"`
			EndLocation struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Ending point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Ending point longitude"`
				NodeInfo  struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			} `json:"endLocation"`
			Skills     []string `json:"skills" description:"Vehicle capabilities such as size or equipment requirements. eg: XL, heavy, etc"`
			TimeWindow struct {
				Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
				End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
			} `json:"timeWindow"`
			Capacity struct {
				Insurance             int64 `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
				Volume                int64 `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
				Weight                int64 `json:"weight" example:"1000" description:"Maximum weight in grams"`
				DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
			} `json:"capacity"`
		}{
			Plate: plate,
			StartLocation: struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
				NodeInfo  struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			}{
				Latitude:  startLocations[i].lat,
				Longitude: startLocations[i].lon,
				NodeInfo: struct {
					ReferenceID string `json:"referenceID"`
				}{
					ReferenceID: startLocations[i].ref,
				},
			},
			Skills: []string{"delivery", "express"},
			TimeWindow: struct {
				Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
				End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
			}{
				Start: "08:00",
				End:   "18:00",
			},
			Capacity: struct {
				Insurance             int64 `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
				Volume                int64 `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
				Weight                int64 `json:"weight" example:"1000" description:"Maximum weight in grams"`
				DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
			}{
				Insurance:             capacities[i].insurance,
				Volume:                capacities[i].volume,
				Weight:                capacities[i].weight,
				DeliveryUnitsQuantity: capacities[i].deliveryUnits,
			},
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles
}

// generateVisits genera 10 visitas (5 por vehículo) con coordenadas aleatorias en La Florida
func generateVisits(baseLat, baseLon float64) []struct {
	Pickup struct {
		Coordinates struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
		} `json:"coordinates"`
		ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Contact     struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			FullName   string `json:"fullName"`
		} `json:"contact"`
		Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"pickup"`
	Delivery struct {
		Coordinates struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
		} `json:"coordinates"`
		ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Contact     struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			FullName   string `json:"fullName"`
		} `json:"contact"`
		Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"delivery"`
	Orders []struct {
		DeliveryUnits []struct {
			Items []struct {
				Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
			} `json:"items"`
			Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
			Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
			Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
			Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
		} `json:"deliveryUnits"`
		ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
	} `json:"orders"`
} {
	var visits []struct {
		Pickup struct {
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
			} `json:"coordinates"`
			ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Contact     struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
			Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"pickup"`
		Delivery struct {
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
			} `json:"coordinates"`
			ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Contact     struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
			Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"delivery"`
		Orders []struct {
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
				} `json:"items"`
				Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
				Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
				Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
				Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
		} `json:"orders"`
	}

	// Nombres y datos de contacto para generar datos realistas
	names := []string{
		"Juan Gonzalez", "Maria Perez", "Ana Rodriguez", "Carlos Morales", "Lucia Herrera",
		"Roberto Silva", "Patricia Vargas", "Fernando Castro", "Gabriela Torres", "Miguel Ruiz",
	}

	// Generar 10 visitas (5 por vehículo)
	for i := 0; i < 10; i++ {
		// Coordenadas aleatorias dentro de La Florida (aproximadamente 5km de radio)
		latOffset := (rand.Float64() - 0.5) * 0.05 // ±0.025 grados ≈ ±2.5km
		lonOffset := (rand.Float64() - 0.5) * 0.05 // ±0.025 grados ≈ ±2.5km

		deliveryLat := baseLat + latOffset
		deliveryLon := baseLon + lonOffset

		// Para algunas visitas, agregar pickup (shipment) - aproximadamente 10% de las visitas
		hasPickup := i%10 == 0 // Solo la primera visita tendrá pickup

		visit := struct {
			Pickup struct {
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
					Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
				} `json:"coordinates"`
				ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
				Contact     struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
				TimeWindow struct {
					Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
					End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
				} `json:"timeWindow"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			} `json:"pickup"`
			Delivery struct {
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
					Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
				} `json:"coordinates"`
				ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
				Contact     struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
				TimeWindow struct {
					Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
					End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
				} `json:"timeWindow"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			} `json:"delivery"`
			Orders []struct {
				DeliveryUnits []struct {
					Items []struct {
						Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
					} `json:"items"`
					Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
					Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
					Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
					Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
				} `json:"deliveryUnits"`
				ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
			} `json:"orders"`
		}{}

		// Configurar delivery
		visit.Delivery.Coordinates.Latitude = deliveryLat
		visit.Delivery.Coordinates.Longitude = deliveryLon
		visit.Delivery.ServiceTime = 30
		visit.Delivery.Skills = []string{"delivery"}
		visit.Delivery.TimeWindow.Start = "08:00"
		visit.Delivery.TimeWindow.End = "17:00"
		visit.Delivery.NodeInfo.ReferenceID = fmt.Sprintf("delivery-%03d", i+1)

		// Configurar pickup si es necesario
		if hasPickup {
			pickupLat := deliveryLat + (rand.Float64()-0.5)*0.01
			pickupLon := deliveryLon + (rand.Float64()-0.5)*0.01
			visit.Pickup.Coordinates.Latitude = pickupLat
			visit.Pickup.Coordinates.Longitude = pickupLon
			visit.Pickup.ServiceTime = 3000
			visit.Pickup.Skills = []string{"delivery"}
			visit.Pickup.TimeWindow.Start = "08:00"
			visit.Pickup.TimeWindow.End = "17:00"
			visit.Pickup.NodeInfo.ReferenceID = fmt.Sprintf("pickup-%03d", i+1)
		}

		// Generar datos de contacto
		nameIndex := i % len(names)
		name := names[nameIndex]
		firstName := name[:len(name)/2]
		lastName := name[len(name)/2:]
		email := fmt.Sprintf("%s.%s@email.com", firstName, lastName)
		phone := fmt.Sprintf("+569%08d", 80000000+i)
		nationalID := fmt.Sprintf("%d-%d", 10000000+i, (i%9)+1)

		visit.Delivery.Contact.Email = email
		visit.Delivery.Contact.FullName = name
		visit.Delivery.Contact.NationalID = nationalID
		visit.Delivery.Contact.Phone = phone

		if hasPickup {
			visit.Pickup.Contact.Email = email
			visit.Pickup.Contact.FullName = name
			visit.Pickup.Contact.NationalID = nationalID
			visit.Pickup.Contact.Phone = phone
		}

		// Generar órdenes - cada visita tiene 1 orden
		order := struct {
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
				} `json:"items"`
				Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
				Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
				Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
				Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
		}{}

		// Generar delivery unit
		insurance := 9000 + rand.Int63n(12000) // Entre 9000 y 21000
		volume := insurance / 10               // Proporcional al insurance
		weight := volume                       // Mismo valor que volume

		order.DeliveryUnits = []struct {
			Items []struct {
				Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
			} `json:"items"`
			Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
			Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
			Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
			Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
		}{
			{
				Items: []struct {
					Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
				}{
					{Sku: fmt.Sprintf("SKU%03d", i+1)},
				},
				Insurance: insurance,
				Volume:    volume,
				Weight:    weight,
				Lpn:       fmt.Sprintf("LPN%03d", i+1),
			},
		}

		order.ReferenceID = fmt.Sprintf("ORD%03d", i+1)
		visit.Orders = append(visit.Orders, order)

		visits = append(visits, visit)
	}

	return visits
}

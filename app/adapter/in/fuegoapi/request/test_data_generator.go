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

	// Coordenadas del punto de inicio único para todos los vehículos
	baseLat := -33.4505803
	baseLon := -70.7857318

	// Generar 10 vehículos con capacidad para 100 visitas cada uno
	vehicles := generateVehicles(baseLat, baseLon)

	// Generar 1000 visitas distribuidas en las tres zonas
	visits := generateVisits()

	return OptimizeFleetRequest{
		PlanReferenceID: "MASSIVE_TEST_PLAN",
		Vehicles:        vehicles,
		Visits:          visits,
	}
}

// generateVehicles genera 10 vehículos con capacidad para 100 visitas cada uno
func generateVehicles(baseLat, baseLon float64) []struct {
	Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
	PoliticalArea struct {
		Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
		District string `json:"district" example:"la florida" description:"District name"`
		Province string `json:"province" example:"santiago" description:"Province name"`
		State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
	} `json:"politicalArea"`
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
		PoliticalArea struct {
			Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
			District string `json:"district" example:"la florida" description:"District name"`
			Province string `json:"province" example:"santiago" description:"Province name"`
			State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
		} `json:"politicalArea"`
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

	// Capacidades para los 10 vehículos (cada uno puede llevar 100 delivery units)
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
		}{100, 999999999, 999999999, 999999999})
	}

	// Todos los vehículos parten desde la misma ubicación
	startLocation := struct {
		lat float64
		lon float64
		ref string
	}{baseLat, baseLon, "depot-central"}

	for i := 0; i < 10; i++ {
		plate := fmt.Sprintf("SERV-%d", 80+i)
		vehicle := struct {
			Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			} `json:"politicalArea"`
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
			PoliticalArea: struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			}{
				Code:     "cl-rm-la-florida",
				District: "la florida",
				Province: "santiago",
				State:    "region metropolitana de santiago",
			},
			StartLocation: struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
				NodeInfo  struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			}{
				Latitude:  startLocation.lat,
				Longitude: startLocation.lon,
				NodeInfo: struct {
					ReferenceID string `json:"referenceID"`
				}{
					ReferenceID: startLocation.ref,
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

// generateVisits genera 1000 visitas distribuidas en tres zonas de Santiago
func generateVisits() []struct {
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
		PoliticalArea struct {
			Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
			District string `json:"district" example:"la florida" description:"District name"`
			Province string `json:"province" example:"santiago" description:"Province name"`
			State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
		} `json:"politicalArea"`
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
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			} `json:"politicalArea"`
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

	// Coordenadas base de las tres zonas
	laFloridaBase := struct{ lat, lon float64 }{-33.5225, -70.575}       // La Florida
	santiagoCentroBase := struct{ lat, lon float64 }{-33.4489, -70.6693} // Santiago Centro
	lasCondesBase := struct{ lat, lon float64 }{-33.4167, -70.5833}      // Las Condes

	// Nombres para generar datos realistas
	names := []string{
		"Juan Gonzalez", "Maria Perez", "Ana Rodriguez", "Carlos Morales", "Lucia Herrera",
		"Roberto Silva", "Patricia Vargas", "Fernando Castro", "Gabriela Torres", "Miguel Ruiz",
		"Carmen Vega", "Diego Mendoza", "Sofia Rojas", "Alejandro Fuentes", "Valentina Soto",
		"Francisco Herrera", "Daniela Morales", "Sebastian Silva", "Camila Torres", "Matias Ruiz",
	}

	// Generar 1000 visitas
	for i := 0; i < 1000; i++ {
		var deliveryLat, deliveryLon float64
		var zoneName string
		var politicalArea struct {
			code     string
			district string
			province string
			state    string
		}

		// Distribuir visitas en las tres zonas
		if i < 300 {
			// 300 visitas en La Florida
			latOffset := (rand.Float64() - 0.5) * 0.05 // ±0.025 grados ≈ ±2.5km
			lonOffset := (rand.Float64() - 0.5) * 0.05
			deliveryLat = laFloridaBase.lat + latOffset
			deliveryLon = laFloridaBase.lon + lonOffset
			zoneName = "la-florida"
			politicalArea = struct {
				code     string
				district string
				province string
				state    string
			}{
				code:     "cl-rm-la-florida",
				district: "la florida",
				province: "santiago",
				state:    "region metropolitana de santiago",
			}
		} else if i < 600 {
			// 300 visitas en Santiago Centro
			latOffset := (rand.Float64() - 0.5) * 0.03 // ±0.015 grados ≈ ±1.5km
			lonOffset := (rand.Float64() - 0.5) * 0.03
			deliveryLat = santiagoCentroBase.lat + latOffset
			deliveryLon = santiagoCentroBase.lon + lonOffset
			zoneName = "santiago-centro"
			politicalArea = struct {
				code     string
				district string
				province string
				state    string
			}{
				code:     "cl-rm-santiago-centro",
				district: "santiago centro",
				province: "santiago",
				state:    "region metropolitana de santiago",
			}
		} else {
			// 400 visitas en Las Condes
			latOffset := (rand.Float64() - 0.5) * 0.04 // ±0.02 grados ≈ ±2km
			lonOffset := (rand.Float64() - 0.5) * 0.04
			deliveryLat = lasCondesBase.lat + latOffset
			deliveryLon = lasCondesBase.lon + lonOffset
			zoneName = "las-condes"
			politicalArea = struct {
				code     string
				district string
				province string
				state    string
			}{
				code:     "cl-rm-las-condes",
				district: "las condes",
				province: "santiago",
				state:    "region metropolitana de santiago",
			}
		}

		// Para algunas visitas, agregar pickup (shipment) - aproximadamente 5% de las visitas
		hasPickup := i%20 == 0 // Solo 5% de las visitas tendrán pickup

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
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
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
		visit.Delivery.NodeInfo.ReferenceID = fmt.Sprintf("delivery-%s-%04d", zoneName, i+1)
		visit.Delivery.PoliticalArea.Code = politicalArea.code
		visit.Delivery.PoliticalArea.District = politicalArea.district
		visit.Delivery.PoliticalArea.Province = politicalArea.province
		visit.Delivery.PoliticalArea.State = politicalArea.state

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
			visit.Pickup.NodeInfo.ReferenceID = fmt.Sprintf("pickup-%s-%04d", zoneName, i+1)
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
					{Sku: fmt.Sprintf("SKU%04d", i+1)},
				},
				Insurance: insurance,
				Volume:    volume,
				Weight:    weight,
				Lpn:       fmt.Sprintf("LPN%04d", i+1),
			},
		}

		order.ReferenceID = fmt.Sprintf("ORD%04d", i+1)
		visit.Orders = append(visit.Orders, order)

		visits = append(visits, visit)
	}

	return visits
}

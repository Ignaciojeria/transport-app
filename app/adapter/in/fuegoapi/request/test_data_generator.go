package request

/*
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
	StartLocation struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
			AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
			Contact      struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
				Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
			} `json:"coordinates"`
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			} `json:"politicalArea"`
			ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"startLocation"`
	EndLocation struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
			AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
			Contact      struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Ending point latitude"`
				Longitude float64 `json:"longitude" example:"-70.6130425" description:"Ending point longitude"`
			} `json:"coordinates"`
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			} `json:"politicalArea"`
			ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
		} `json:"addressInfo"`
		NodeInfo struct {
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
			AddressInfo struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code            string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					AdminAreaLevel1 string `json:"adminAreaLevel1" example:"region metropolitana de santiago" description:"Administrative area level 1"`
					AdminAreaLevel2 string `json:"adminAreaLevel2" example:"santiago" description:"Administrative area level 2"`
					AdminAreaLevel3 string `json:"adminAreaLevel3" example:"la florida" description:"Administrative area level 3"`
					AdminAreaLevel4 string `json:"adminAreaLevel4" example:"" description:"Administrative area level 4"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"startLocation"`
		EndLocation struct {
			AddressInfo struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Ending point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Ending point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
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
			StartLocation struct {
				AddressInfo struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				} `json:"addressInfo"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			} `json:"startLocation"`
			EndLocation struct {
				AddressInfo struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Ending point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Ending point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				} `json:"addressInfo"`
				NodeInfo struct {
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
				AddressInfo struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				} `json:"addressInfo"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			}{
				AddressInfo: struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				}{
					AddressLine1: "Depósito Central",
					AddressLine2: "Bodega Principal",
					Contact: struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					}{
						Email:      "depot@transport.com",
						Phone:      "+56912345678",
						NationalID: "12345678-9",
						FullName:   "Depósito Central",
					},
					Coordinates: struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
					}{
						Latitude:  startLocation.lat,
						Longitude: startLocation.lon,
					},
					PoliticalArea: struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					}{
						Code:     "cl-rm-santiago-centro",
						District: "santiago centro",
						Province: "santiago",
						State:    "region metropolitana de santiago",
					},
					ZipCode: "7500000",
				},
				NodeInfo: struct {
					ReferenceID string `json:"referenceID"`
				}{
					ReferenceID: startLocation.ref,
				},
			},
			// EndLocation se mantiene vacío (solo origen)
			EndLocation: struct {
				AddressInfo struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Ending point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Ending point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				} `json:"addressInfo"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
			}{},
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
		Instructions string `json:"instructions" example:"Recoger en recepción" description:"Instructions for pickup"`
		AddressInfo  struct {
			AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
			AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
			Contact      struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Pickup point latitude"`
				Longitude float64 `json:"longitude" example:"-70.6130425" description:"Pickup point longitude"`
			} `json:"coordinates"`
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			} `json:"politicalArea"`
			ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
		ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow  struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
	} `json:"pickup"`
	Delivery struct {
		Instructions string `json:"instructions" example:"Entregar en recepción" description:"Instructions for delivery"`
		AddressInfo  struct {
			AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
			AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
			Contact      struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Delivery point latitude"`
				Longitude float64 `json:"longitude" example:"-70.6130425" description:"Delivery point longitude"`
			} `json:"coordinates"`
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
				District string `json:"district" example:"la florida" description:"District name"`
				Province string `json:"province" example:"santiago" description:"Province name"`
				State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
			} `json:"politicalArea"`
			ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
		ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow  struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
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
			Instructions string `json:"instructions" example:"Recoger en recepción" description:"Instructions for pickup"`
			AddressInfo  struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Pickup point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Pickup point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
			ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow  struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
		} `json:"pickup"`
		Delivery struct {
			Instructions string `json:"instructions" example:"Entregar en recepción" description:"Instructions for delivery"`
			AddressInfo  struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Delivery point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Delivery point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
			ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow  struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
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

		visit := struct {
			Pickup struct {
				Instructions string `json:"instructions" example:"Recoger en recepción" description:"Instructions for pickup"`
				AddressInfo  struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Pickup point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Pickup point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				} `json:"addressInfo"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
				ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
				Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
				TimeWindow  struct {
					Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
					End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
				} `json:"timeWindow"`
			} `json:"pickup"`
			Delivery struct {
				Instructions string `json:"instructions" example:"Entregar en recepción" description:"Instructions for delivery"`
				AddressInfo  struct {
					AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
					AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
					Contact      struct {
						Email      string `json:"email"`
						Phone      string `json:"phone"`
						NationalID string `json:"nationalID"`
						FullName   string `json:"fullName"`
					} `json:"contact"`
					Coordinates struct {
						Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Delivery point latitude"`
						Longitude float64 `json:"longitude" example:"-70.6130425" description:"Delivery point longitude"`
					} `json:"coordinates"`
					PoliticalArea struct {
						Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
						District string `json:"district" example:"la florida" description:"District name"`
						Province string `json:"province" example:"santiago" description:"Province name"`
						State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
					} `json:"politicalArea"`
					ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
				} `json:"addressInfo"`
				NodeInfo struct {
					ReferenceID string `json:"referenceID"`
				} `json:"nodeInfo"`
				ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
				Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
				TimeWindow  struct {
					Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
					End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
				} `json:"timeWindow"`
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
		visit.Delivery.Instructions = "Entregar en recepción"
		visit.Delivery.AddressInfo.Coordinates.Latitude = deliveryLat
		visit.Delivery.AddressInfo.Coordinates.Longitude = deliveryLon
		visit.Delivery.ServiceTime = 30
		visit.Delivery.Skills = []string{"delivery"}
		visit.Delivery.TimeWindow.Start = "08:00"
		visit.Delivery.TimeWindow.End = "23:00"
		visit.Delivery.NodeInfo.ReferenceID = fmt.Sprintf("delivery-%s-%04d", zoneName, i+1)
		visit.Delivery.AddressInfo.PoliticalArea.Code = politicalArea.code
		visit.Delivery.AddressInfo.PoliticalArea.District = politicalArea.district
		visit.Delivery.AddressInfo.PoliticalArea.Province = politicalArea.province
		visit.Delivery.AddressInfo.PoliticalArea.State = politicalArea.state

		// Generar datos de dirección para delivery
		visit.Delivery.AddressInfo.AddressLine1 = fmt.Sprintf("Calle %d", 1000+i)
		visit.Delivery.AddressInfo.AddressLine2 = fmt.Sprintf("Piso %d", (i%20)+1)
		visit.Delivery.AddressInfo.ZipCode = fmt.Sprintf("750%04d", 1000+i)

		// Pickup se mantiene vacío (no se generan pickups)
		visit.Pickup.Instructions = "Recoger en recepción"

		// Generar datos de contacto
		nameIndex := i % len(names)
		name := names[nameIndex]
		firstName := name[:len(name)/2]
		lastName := name[len(name)/2:]
		email := fmt.Sprintf("%s.%s@email.com", firstName, lastName)
		phone := fmt.Sprintf("+569%08d", 80000000+i)
		nationalID := fmt.Sprintf("%d-%d", 10000000+i, (i%9)+1)

		visit.Delivery.AddressInfo.Contact.Email = email
		visit.Delivery.AddressInfo.Contact.FullName = name
		visit.Delivery.AddressInfo.Contact.NationalID = nationalID
		visit.Delivery.AddressInfo.Contact.Phone = phone

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
*/

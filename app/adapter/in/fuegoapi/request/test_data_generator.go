package request

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// formatDecimalForGoogleSheets convierte un float64 a string con formato compatible con Google Sheets
func formatDecimalForGoogleSheets(value float64) string {
	// Usar coma como separador decimal para compatibilidad con Google Sheets
	return strings.Replace(fmt.Sprintf("%.7f", value), ".", ",", -1)
}

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
func generateVehicles(baseLat, baseLon float64) []OptimizeFleetVehicle {
	var vehicles []OptimizeFleetVehicle

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

		startAddressInfo := OptimizeFleetAddressInfo{
			AddressLine1: "Depósito Central",
			AddressLine2: "Bodega Principal",
			Contact: OptimizeFleetContact{
				Email:      "depot@transport.com",
				Phone:      "+56912345678",
				NationalID: "12345678-9",
				FullName:   "Depósito Central",
			},
			Coordinates: OptimizeFleetCoordinates{
				Latitude:  startLocation.lat,
				Longitude: startLocation.lon,
			},
			PoliticalArea: OptimizeFleetPoliticalArea{
				Code:            "cl-rm-santiago-centro",
				AdminAreaLevel1: "region metropolitana de santiago",
				AdminAreaLevel2: "santiago",
				AdminAreaLevel3: "santiago centro",
				AdminAreaLevel4: "",
			},
			ZipCode: "7500000",
		}

		vehicle := OptimizeFleetVehicle{
			Plate: plate,
			StartLocation: OptimizeFleetVehicleLocation{
				AddressInfo: startAddressInfo,
				NodeInfo:    OptimizeFleetNodeInfo{ReferenceID: startLocation.ref},
			},
			EndLocation: OptimizeFleetVehicleLocation{}, // EndLocation vacío
			Skills:      []string{"delivery", "express"},
			TimeWindow: OptimizeFleetTimeWindow{
				Start: "08:00",
				End:   "18:00",
			},
			Capacity: OptimizeFleetVehicleCapacity{
				Insurance: capacities[i].insurance,
				Volume:    capacities[i].volume,
				Weight:    capacities[i].weight,
			},
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles
}

// generateVisits genera 1000 visitas distribuidas en tres zonas de Santiago
func generateVisits() []OptimizeFleetVisit {
	var visits []OptimizeFleetVisit

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

		addressInfo := OptimizeFleetAddressInfo{
			AddressLine1: fmt.Sprintf("Calle %d", 1000+i),
			AddressLine2: fmt.Sprintf("Piso %d", (i%20)+1),
			Contact: OptimizeFleetContact{
				Email:      fmt.Sprintf("%s@email.com", strings.ToLower(strings.ReplaceAll(names[i%len(names)], " ", "."))),
				Phone:      fmt.Sprintf("+569%08d", 80000000+i),
				NationalID: fmt.Sprintf("%d-%d", 10000000+i, (i%9)+1),
				FullName:   names[i%len(names)],
			},
			Coordinates: OptimizeFleetCoordinates{
				Latitude:  deliveryLat,
				Longitude: deliveryLon,
			},
			PoliticalArea: OptimizeFleetPoliticalArea{
				Code:            politicalArea.code,
				AdminAreaLevel1: politicalArea.state,
				AdminAreaLevel2: politicalArea.province,
				AdminAreaLevel3: politicalArea.district,
				AdminAreaLevel4: "",
			},
			ZipCode: fmt.Sprintf("750%04d", 1000+i),
		}

		visit := OptimizeFleetVisit{
			Pickup: OptimizeFleetVisitLocation{}, // Pickup vacío
			Delivery: OptimizeFleetVisitLocation{
				Instructions: "Entregar en recepción",
				AddressInfo:  addressInfo,
				NodeInfo:     OptimizeFleetNodeInfo{ReferenceID: fmt.Sprintf("delivery-%s-%04d", zoneName, i+1)},
				ServiceTime:  30,
				TimeWindow: OptimizeFleetTimeWindow{
					Start: "08:00",
					End:   "23:00",
				},
			},
			Orders: []OptimizeFleetOrder{
				{
					DeliveryUnits: []OptimizeFleetDeliveryUnit{
						{
							Items: []OptimizeFleetItem{
								{Sku: fmt.Sprintf("SKU%04d", i+1)},
							},
							Price:  9000 + rand.Int63n(12000),        // Entre 9000 y 21000
							Volume: (9000 + rand.Int63n(12000)) / 10, // Proporcional al insurance
							Weight: (9000 + rand.Int63n(12000)) / 10, // Mismo valor que volume
							Lpn:    fmt.Sprintf("LPN%04d", i+1),
							Skills: []string{"delivery"},
						},
					},
					ReferenceID: fmt.Sprintf("ORD%04d", i+1),
				},
			},
		}

		visits = append(visits, visit)
	}

	return visits
}

// GenerateFleetAndVisitsCSV genera dos archivos CSV: uno para fleet y otro para visits
func GenerateFleetAndVisitsCSV() error {
	// Generar datos de prueba
	fleetData := generateFleetData()
	visitsData := generateVisitsData()

	// Generar CSV de fleet
	if err := generateFleetCSV(fleetData); err != nil {
		return fmt.Errorf("error generando CSV de fleet: %w", err)
	}

	// Generar CSV de visits
	if err := generateVisitsCSV(visitsData); err != nil {
		return fmt.Errorf("error generando CSV de visits: %w", err)
	}

	return nil
}

// FleetData representa la estructura de datos para el CSV de fleet
type FleetData struct {
	ID                     string `json:"id"`
	Insurance              string `json:"insurance"`
	StartLocationLatitude  string `json:"startLocationLatitude"`
	StartLocationLongitude string `json:"startLocationLongitude"`
	Volume                 string `json:"volume"`
	Weight                 string `json:"weight"`
	MaxPackageQuantity     string `json:"maxPackageQuantity"`
}

// VisitData representa la estructura de datos para el CSV de visits
type VisitData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	Volume      string `json:"volume"`
	Weight      string `json:"weight"`
}

// generateFleetData genera datos de fleet para el CSV
func generateFleetData() []FleetData {
	var fleet []FleetData

	// Coordenadas base del depósito
	baseLat := -33.4505803
	baseLon := -70.7857318

	// Generar 10 vehículos
	for i := 0; i < 10; i++ {
		fleet = append(fleet, FleetData{
			ID:                     fmt.Sprintf("vehicle_%d", i+1),
			Insurance:              "100000",
			StartLocationLatitude:  formatDecimalForGoogleSheets(baseLat),
			StartLocationLongitude: formatDecimalForGoogleSheets(baseLon),
			Volume:                 "10000",
			Weight:                 "10000",
			MaxPackageQuantity:     "1000",
		})
	}

	return fleet
}

// generateVisitsData genera datos de visits para el CSV
func generateVisitsData() []VisitData {
	var visits []VisitData

	// Nombres para generar datos realistas
	names := []string{
		"Ignacio Jeriaa", "Juan Gonzalez", "Maria Perez", "Ana Rodriguez", "Carlos Morales",
		"Lucia Herrera", "Roberto Silva", "Patricia Vargas", "Fernando Castro", "Gabriela Torres",
	}

	// Coordenadas base de las tres zonas
	laFloridaBase := struct{ lat, lon float64 }{-33.5225, -70.575}
	santiagoCentroBase := struct{ lat, lon float64 }{-33.4489, -70.6693}
	lasCondesBase := struct{ lat, lon float64 }{-33.4167, -70.5833}

	// Generar 1000 visitas
	for i := 0; i < 1000; i++ {
		var deliveryLat, deliveryLon float64
		var zoneName string

		// Distribuir visitas en las tres zonas
		if i < 300 {
			// 300 visitas en La Florida
			latOffset := (rand.Float64() - 0.5) * 0.05
			lonOffset := (rand.Float64() - 0.5) * 0.05
			deliveryLat = laFloridaBase.lat + latOffset
			deliveryLon = laFloridaBase.lon + lonOffset
			zoneName = "la-florida"
		} else if i < 600 {
			// 300 visitas en Santiago Centro
			latOffset := (rand.Float64() - 0.5) * 0.03
			lonOffset := (rand.Float64() - 0.5) * 0.03
			deliveryLat = santiagoCentroBase.lat + latOffset
			deliveryLon = santiagoCentroBase.lon + lonOffset
			zoneName = "santiago-centro"
		} else {
			// 400 visitas en Las Condes
			latOffset := (rand.Float64() - 0.5) * 0.04
			lonOffset := (rand.Float64() - 0.5) * 0.04
			deliveryLat = lasCondesBase.lat + latOffset
			deliveryLon = lasCondesBase.lon + lonOffset
			zoneName = "las-condes"
		}

		// Generar ID único para la visita
		visitID := fmt.Sprintf("%d%s", 100+i, strings.ToUpper(zoneName[:2]))

		// Generar dirección realista
		address := fmt.Sprintf("Calle %d, Piso %d, %s", 1000+i, (i%20)+1, zoneName)

		visits = append(visits, VisitData{
			ID:          visitID,
			Name:        names[i%len(names)],
			Phone:       fmt.Sprintf("569%08d", 80000000+i),
			Address:     address,
			Latitude:    formatDecimalForGoogleSheets(deliveryLat),
			Longitude:   formatDecimalForGoogleSheets(deliveryLon),
			Description: "bebida 350ml",
			Price:       "100",
			Quantity:    "13",
			Volume:      "2",
			Weight:      "12",
		})
	}

	return visits
}

// generateFleetCSV genera el archivo CSV de fleet
func generateFleetCSV(fleet []FleetData) error {
	file, err := os.Create("fleet.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezados
	headers := []string{"id", "insurance", "startLocationLatitude", "startLocationLongitude", "volume", "weight", "maxPackageQuantity"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Escribir datos
	for _, f := range fleet {
		row := []string{
			f.ID,
			f.Insurance,
			f.StartLocationLatitude,
			f.StartLocationLongitude,
			f.Volume,
			f.Weight,
			f.MaxPackageQuantity,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// generateVisitsCSV genera el archivo CSV de visits
func generateVisitsCSV(visits []VisitData) error {
	file, err := os.Create("visits.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezados
	headers := []string{"id", "name", "phone", "address", "latitude", "longitude", "description", "price", "quantity", "volume", "weight"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Escribir datos
	for _, v := range visits {
		row := []string{
			v.ID,
			v.Name,
			v.Phone,
			v.Address,
			v.Latitude,
			v.Longitude,
			v.Description,
			v.Price,
			v.Quantity,
			v.Volume,
			v.Weight,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

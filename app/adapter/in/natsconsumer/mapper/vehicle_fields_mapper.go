package mapper

type VehicleFieldMapper struct {
	synonyms    map[string][]string // canonical -> []alts
	index       map[string]string   // normalized alt -> canonical
	orderedKeys []string
}

const (
	VehicleKeyEndLocationLatitude    = "endLocationLatitude"
	VehicleKeyEndLocationLongitude   = "endLocationLongitude"
	VehicleKeyID                     = "id"
	VehicleKeyInsurance              = "insurance"
	VehicleKeyStartLocationLatitude  = "startLocationLatitude"
	VehicleKeyStartLocationLongitude = "startLocationLongitude"
	VehicleKeyVolume                 = "volume"
	VehicleKeyWeight                 = "weight"
	VehicleMaxPackageQuantity        = "maxPackageQuantity"
)

func NewVehicleFieldMapper() *VehicleFieldMapper {
	syn := map[string][]string{
		VehicleKeyEndLocationLatitude:    {"end location (latitude,longitude)", "end_location_lat", "end_latitude", "destino_lat", "destination_lat", "end_lat", "end_coordinates_lat", "latitud_destino", "latitud_final"},
		VehicleKeyEndLocationLongitude:   {"end location (latitude,longitude)", "end_location_lon", "end_longitude", "destino_lon", "destination_lon", "end_lon", "end_coordinates_lon", "longitud_destino", "longitud_final"},
		VehicleKeyID:                     {"id", "vehicle_id", "vehiculo_id", "identificador", "identifier", "id_vehiculo", "vehiculo", "auto", "carro", "coche"},
		VehicleKeyInsurance:              {"insurance (currency)", "insurance", "seguro", "insurance_amount", "insurance_value", "seguro_monto", "seguro_valor", "valor_seguro", "monto_seguro", "cobertura"},
		VehicleKeyStartLocationLatitude:  {"start location (lat,lon)", "start_location_lat", "start_latitude", "origen_lat", "origin_lat", "start_lat", "start_coordinates_lat", "latitud_origen", "latitud_inicial"},
		VehicleKeyStartLocationLongitude: {"start location (lat,lon)", "start_location_lon", "start_longitude", "origen_lon", "origin_lon", "start_lon", "start_coordinates_lon", "longitud_origen", "longitud_inicial"},
		VehicleKeyVolume:                 {"volume (cm3)", "volume", "volumen", "volume_cm3", "volumen_cm3", "cm3", "cubic_cm", "cubic_centimeters", "capacidad", "espacio"},
		VehicleKeyWeight:                 {"weight (grams)", "weight", "peso", "weight_grams", "peso_gramos", "grams", "gramos", "carga_maxima", "peso_maximo"},
		VehicleMaxPackageQuantity:        {"max package quantity", "max_package_quantity", "max_paquetes", "max_paquete", "max_paquetes_maximo", "max_paquete_maximo", "max_paquetes_maximo", "max_paquete_maximo"},
	}

	ordered := []string{
		VehicleKeyEndLocationLatitude,
		VehicleKeyEndLocationLongitude,
		VehicleKeyID,
		VehicleKeyInsurance,
		VehicleKeyStartLocationLatitude,
		VehicleKeyStartLocationLongitude,
		VehicleKeyVolume,
		VehicleKeyWeight,
		VehicleMaxPackageQuantity,
	}

	idx := buildIndex(syn)
	return &VehicleFieldMapper{synonyms: syn, orderedKeys: ordered, index: idx}
}

// Map rewrites input keys -> canonical keys. Order of columns is irrelevant.
// MapMetadata returns a map where canonical keys point to the original input keys.
func (m *VehicleFieldMapper) Map(input map[string]any) map[string]string {
	out := make(map[string]string, len(m.orderedKeys))

	// Inicializa el mapa de salida con valores vacíos para todos los campos canónicos.
	// Esto es crucial para simular el comportamiento del agente.
	for _, k := range m.orderedKeys {
		out[k] = ""
	}

	used := make(map[string]bool, len(m.synonyms))
	for k, _ := range input {
		if c, ok := m.index[norm(k)]; ok && !used[c] {
			out[c] = k // Aquí cambiamos: asignamos la clave original 'k' como valor.
			used[c] = true
		}
	}
	return out
}

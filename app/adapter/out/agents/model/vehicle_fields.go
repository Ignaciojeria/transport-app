package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/genai"
)

type VehicleFieldMappingSchema struct {
	schema      *genai.Schema
	synonyms    map[string][]string
	orderedKeys []string
}

func NewVehicleFieldMappingSchema() *VehicleFieldMappingSchema {
	props := map[string]*genai.Schema{
		VehicleKeyEndLocationLatitude:    {Type: genai.TypeString},
		VehicleKeyEndLocationLongitude:   {Type: genai.TypeString},
		VehicleKeyID:                     {Type: genai.TypeString},
		VehicleKeyInsurance:              {Type: genai.TypeString},
		VehicleKeyStartLocationLatitude:  {Type: genai.TypeString},
		VehicleKeyStartLocationLongitude: {Type: genai.TypeString},
		VehicleKeyVolume:                 {Type: genai.TypeString},
		VehicleKeyWeight:                 {Type: genai.TypeString},
	}

	schema := &genai.Schema{
		Type:       genai.TypeObject,
		Properties: props,
		Required: []string{
			VehicleKeyEndLocationLatitude, VehicleKeyEndLocationLongitude, VehicleKeyID, VehicleKeyInsurance,
			VehicleKeyStartLocationLatitude, VehicleKeyStartLocationLongitude, VehicleKeyVolume, VehicleKeyWeight,
		},
	}

	synonyms := map[string][]string{
		VehicleKeyEndLocationLatitude:    {"end location (latitude,longitude)", "end_location_lat", "end_latitude", "destino_lat", "destination_lat", "end_lat", "end_coordinates_lat"},
		VehicleKeyEndLocationLongitude:   {"end location (latitude,longitude)", "end_location_lon", "end_longitude", "destino_lon", "destination_lon", "end_lon", "end_coordinates_lon"},
		VehicleKeyID:                     {"id", "vehicle_id", "vehiculo_id", "identificador", "identifier"},
		VehicleKeyInsurance:              {"insurance (currency)", "insurance", "seguro", "insurance_amount", "insurance_value", "seguro_monto", "seguro_valor"},
		VehicleKeyStartLocationLatitude:  {"start location (lat,lon)", "start_location_lat", "start_latitude", "origen_lat", "origin_lat", "start_lat", "start_coordinates_lat"},
		VehicleKeyStartLocationLongitude: {"start location (lat,lon)", "start_location_lon", "start_longitude", "origen_lon", "origin_lon", "start_lon", "start_coordinates_lon"},
		VehicleKeyVolume:                 {"volume (cm3)", "volume", "volumen", "volume_cm3", "volumen_cm3", "cm3", "cubic_cm", "cubic_centimeters"},
		VehicleKeyWeight:                 {"weight (grams)", "weight", "peso", "weight_grams", "peso_gramos", "grams", "gramos"},
	}

	ordered := []string{
		VehicleKeyEndLocationLatitude, VehicleKeyEndLocationLongitude, VehicleKeyID, VehicleKeyInsurance,
		VehicleKeyStartLocationLatitude, VehicleKeyStartLocationLongitude, VehicleKeyVolume, VehicleKeyWeight,
	}

	return &VehicleFieldMappingSchema{
		schema:      schema,
		synonyms:    synonyms,
		orderedKeys: ordered,
	}
}

func (v *VehicleFieldMappingSchema) Schema() *genai.Schema { return v.schema }

func (v *VehicleFieldMappingSchema) Prompt(input interface{}) string {
	// Compacto -> menos tokens. Fallback simple si falla.
	in, err := json.Marshal(input)
	if err != nil {
		in = []byte("{}")
	}

	dict := v.renderSynonymsSection()

	return fmt.Sprintf(`
Eres un asistente de normalización. A partir del JSON de entrada, devuelve EXCLUSIVAMENTE UN SOLO OBJETO JSON que cumpla EXACTAMENTE este contrato por cada clave canónica:

{
  "endLocationLatitude": "string",
  "endLocationLongitude": "string",
  "id": "string", 
  "insurance": "string",
  "startLocationLatitude": "string",
  "startLocationLongitude": "string",
  "volume": "string",
  "weight": "string"
}

Entrada:
%s

Reglas:
- La salida DEBE ser SOLO UN OBJETO JSON válido (NO un array).
- Para cada campo, devuelve el NOMBRE EXACTO DE LA CLAVE de entrada que corresponde a ese campo.
- Coincidencia case-insensitive; trata espacios, guiones y guiones_bajos como equivalentes.
- Si hay múltiples claves candidatas, usa la prioridad según el diccionario (las primeras son más específicas).
- Si no encuentras una clave para un campo, deja el valor vacío "".
- Para coordenadas combinadas como "start location (lat,lon)", asigna la misma clave a ambos campos de lat y lon.

Diccionario de sinónimos:
%s

Ejemplo:
Si el input tiene {"id": "vehicle_1", "start location (lat,lon)": "-33.5148399,-70.6105876"}, devuelve:
{
  "endLocationLatitude": "",
  "endLocationLongitude": "",
  "id": "id",
  "insurance": "",
  "startLocationLatitude": "start location (lat,lon)",
  "startLocationLongitude": "start location (lat,lon)",
  "volume": "",
  "weight": ""
}
`,
		string(in),
		dict,
	)
}

// ----- helpers privados -----

func (v *VehicleFieldMappingSchema) renderSynonymsSection() string {
	// garantizar orden por v.orderedKeys y sinónimos alfabéticos para determinismo
	var b strings.Builder
	for _, k := range v.orderedKeys {
		syns := append([]string(nil), v.synonyms[k]...)
		sort.Strings(syns)
		if len(syns) > 0 {
			b.WriteString(fmt.Sprintf("- %s ← %s\n", k, strings.Join(syns, ", ")))
		}
	}
	return b.String()
}

// ===== 1) Constantes (fuente de verdad) =====
const (
	VehicleKeyEndLocationLatitude    = "endLocationLatitude"
	VehicleKeyEndLocationLongitude   = "endLocationLongitude"
	VehicleKeyID                     = "id"
	VehicleKeyInsurance              = "insurance"
	VehicleKeyStartLocationLatitude  = "startLocationLatitude"
	VehicleKeyStartLocationLongitude = "startLocationLongitude"
	VehicleKeyVolume                 = "volume"
	VehicleKeyWeight                 = "weight"
)

// SchemaForVehicleFieldMapping devuelve un schema para mapeo de claves de vehículos
func SchemaForVehicleFieldMapping() *genai.Schema {
	props := map[string]*genai.Schema{
		VehicleKeyEndLocationLatitude:    {Type: genai.TypeString},
		VehicleKeyEndLocationLongitude:   {Type: genai.TypeString},
		VehicleKeyID:                     {Type: genai.TypeString},
		VehicleKeyInsurance:              {Type: genai.TypeString},
		VehicleKeyStartLocationLatitude:  {Type: genai.TypeString},
		VehicleKeyStartLocationLongitude: {Type: genai.TypeString},
		VehicleKeyVolume:                 {Type: genai.TypeString},
		VehicleKeyWeight:                 {Type: genai.TypeString},
	}

	return &genai.Schema{
		Type:       genai.TypeObject,
		Properties: props,
		Required: []string{
			VehicleKeyEndLocationLatitude, VehicleKeyEndLocationLongitude, VehicleKeyID, VehicleKeyInsurance,
			VehicleKeyStartLocationLatitude, VehicleKeyStartLocationLongitude, VehicleKeyVolume, VehicleKeyWeight,
		},
	}
}

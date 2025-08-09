package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"transport-app/app/shared/projection/deliveryunits"

	"google.golang.org/genai"
)

type VisitFieldMappingSchema struct {
	schema      *genai.Schema
	synonyms    map[string][]string
	orderedKeys []string
}

func NewVisitFieldMappingSchema() *VisitFieldMappingSchema {
	projection := deliveryunits.NewProjection()

	props := map[string]*genai.Schema{
		VisitKeyTitle:                                         {Type: genai.TypeString},
		projection.ReferenceID().String():                     {Type: genai.TypeString},
		projection.DestinationContactFullName().String():      {Type: genai.TypeString},
		projection.DestinationContactPhone().String():         {Type: genai.TypeString},
		projection.DestinationAddressLine1().String():         {Type: genai.TypeString},
		projection.DestinationCoordinatesLatitude().String():  {Type: genai.TypeString},
		projection.DestinationCoordinatesLongitude().String(): {Type: genai.TypeString},
		projection.DeliveryUnitVolume().String():              {Type: genai.TypeString},
		projection.DeliveryUnitWeight().String():              {Type: genai.TypeString},
		VisitKeyPrice:                                         {Type: genai.TypeString},
	}

	schema := &genai.Schema{
		Type:       genai.TypeObject,
		Properties: props,
		Required: []string{
			VisitKeyTitle, projection.ReferenceID().String(), projection.DestinationContactFullName().String(),
			projection.DestinationContactPhone().String(), projection.DestinationAddressLine1().String(),
			projection.DestinationCoordinatesLatitude().String(), projection.DestinationCoordinatesLongitude().String(),
			projection.DeliveryUnitVolume().String(), projection.DeliveryUnitWeight().String(), VisitKeyPrice,
		},
	}

	synonyms := map[string][]string{
		VisitKeyTitle:                                         {"título", "titulo", "title", "heading", "nombre pedido", "nombre_orden", "order_title", "desc", "descripción", "descripcion"},
		projection.ReferenceID().String():                     {"id", "order_id", "reference_id", "reference", "ref", "folio", "numero", "nro", "num", "id_pedido", "id_orden"},
		projection.DestinationContactFullName().String():      {"apodo", "alias", "nickname", "nombre cliente", "nombre_cliente", "cliente", "nombre", "name", "full_name", "fullname", "contact_name", "customer_name", "razon social", "razón social", "business_name"},
		projection.DestinationContactPhone().String():         {"telf", "teléfono", "telefono", "fono", "celular", "móvil", "movil", "whatsapp", "phone", "phone_number", "contact_phone", "customer_phone"},
		projection.DestinationAddressLine1().String():         {"direccion", "dirección", "address", "street", "calle", "domicilio", "address_line", "address1", "address_1"},
		projection.DestinationCoordinatesLatitude().String():  {"lat", "latitude", "latitud"},
		projection.DestinationCoordinatesLongitude().String(): {"lon", "lng", "long", "longitude", "longitud"},
		projection.DeliveryUnitVolume().String():              {"total volume (cm3)", "total_volume_cm3", "volume_cm3", "volumen_cm3", "volumen", "volume", "cm3", "cubic_cm", "cubic_centimeters"},
		projection.DeliveryUnitWeight().String():              {"total weight (grams)", "weight_g", "weight_grams", "peso_gramos", "grams", "gramos", "weight", "peso"},
		VisitKeyPrice:                                         {"price", "precio", "amount", "monto", "cost", "costo"},
	}

	ordered := []string{
		VisitKeyTitle, projection.ReferenceID().String(), projection.DestinationContactFullName().String(),
		projection.DestinationContactPhone().String(), projection.DestinationAddressLine1().String(),
		projection.DestinationCoordinatesLatitude().String(), projection.DestinationCoordinatesLongitude().String(),
		projection.DeliveryUnitVolume().String(), projection.DeliveryUnitWeight().String(), VisitKeyPrice,
	}

	return &VisitFieldMappingSchema{
		schema:      schema,
		synonyms:    synonyms,
		orderedKeys: ordered,
	}
}

func (v *VisitFieldMappingSchema) Schema() *genai.Schema { return v.schema }

func (v *VisitFieldMappingSchema) Prompt(input interface{}) string {
	// Compacto -> menos tokens. Fallback simple si falla.
	in, err := json.Marshal(input)
	if err != nil {
		in = []byte("{}")
	}

	dict := v.renderSynonymsSection()

	return fmt.Sprintf(`
Eres un asistente de normalización. A partir del JSON de entrada, devuelve EXCLUSIVAMENTE UN SOLO OBJETO JSON que cumpla EXACTAMENTE este contrato por cada clave canónica:

{
  "title": "string",
  "referenceID": "string", 
  "fullName": "string",
  "phone": "string",
  "address": "string",
  "latitude": "string",
  "longitude": "string",
  "volume": "string",
  "weight": "string",
  "price": "string"
}

Entrada:
%s

Reglas:
- La salida DEBE ser SOLO UN OBJETO JSON válido (NO un array).
- Para cada campo, devuelve el NOMBRE EXACTO DE LA CLAVE de entrada que corresponde a ese campo.
- Coincidencia case-insensitive; trata espacios, guiones y guiones_bajos como equivalentes.
- Si hay múltiples claves candidatas, usa la prioridad según el diccionario (las primeras son más específicas).
- Si no encuentras una clave para un campo, deja el valor vacío "".

Diccionario de sinónimos:
%s

Ejemplo:
Si el input tiene {"nombre_cliente": "Juan", "telf": "123"}, devuelve:
{
  "title": "",
  "referenceID": "",
  "fullName": "nombre_cliente",
  "phone": "telf",
  "address": "",
  "latitude": "",
  "longitude": "",
  "volume": "",
  "weight": "",
  "price": ""
}
`,
		string(in),
		dict,
	)
}

// ----- helpers privados -----

func (v *VisitFieldMappingSchema) renderSynonymsSection() string {
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
	VisitKeyTitle = "title"
	VisitKeyPrice = "price"
)

// SchemaForFieldMapping devuelve un schema para mapeo de claves
func SchemaForFieldMapping() *genai.Schema {
	projection := deliveryunits.NewProjection()

	props := map[string]*genai.Schema{
		VisitKeyTitle:                                         {Type: genai.TypeString},
		projection.ReferenceID().String():                     {Type: genai.TypeString},
		projection.DestinationContactFullName().String():      {Type: genai.TypeString},
		projection.DestinationContactPhone().String():         {Type: genai.TypeString},
		projection.DestinationAddressLine1().String():         {Type: genai.TypeString},
		projection.DestinationCoordinatesLatitude().String():  {Type: genai.TypeString},
		projection.DestinationCoordinatesLongitude().String(): {Type: genai.TypeString},
		projection.DeliveryUnitVolume().String():              {Type: genai.TypeString},
		projection.DeliveryUnitWeight().String():              {Type: genai.TypeString},
		VisitKeyPrice:                                         {Type: genai.TypeString},
	}

	return &genai.Schema{
		Type:       genai.TypeObject,
		Properties: props,
		Required: []string{
			VisitKeyTitle, projection.ReferenceID().String(), projection.DestinationContactFullName().String(),
			projection.DestinationContactPhone().String(), projection.DestinationAddressLine1().String(),
			projection.DestinationCoordinatesLatitude().String(), projection.DestinationCoordinatesLongitude().String(),
			projection.DeliveryUnitVolume().String(), projection.DeliveryUnitWeight().String(), VisitKeyPrice,
		},
	}
}

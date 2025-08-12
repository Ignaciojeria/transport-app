package mapper

import (
	"transport-app/app/shared/projection/deliveryunits"
)

type VisitFieldMapper struct {
	synonyms    map[string][]string // canonical -> []alts
	index       map[string]string   // normalized alt -> canonical
	orderedKeys []string
}

func NewVisitFieldMapper() *VisitFieldMapper {
	proj := deliveryunits.NewProjection()

	syn := map[string][]string{
		"title":                     {"título", "titulo", "title", "heading", "nombre pedido", "nombre_orden", "order_title", "desc", "descripción", "descripcion"},
		proj.ReferenceID().String(): {"id", "order_id", "reference_id", "reference", "ref", "folio", "numero", "nro", "num", "id_pedido", "id_orden"},
		proj.DestinationContactFullName().String():      {"apodo", "alias", "nickname", "nombre cliente", "nombre_cliente", "cliente", "nombre", "name", "full_name", "fullname", "contact_name", "customer_name", "razon social", "razón social", "business_name"},
		proj.DestinationContactPhone().String():         {"telf", "teléfono", "telefono", "fono", "celular", "móvil", "movil", "whatsapp", "phone", "phone_number", "contact_phone", "customer_phone"},
		proj.DestinationAddressLine1().String():         {"direccion", "dirección", "address", "street", "calle", "domicilio", "address_line", "address1", "address_1"},
		proj.DestinationCoordinatesLatitude().String():  {"lat", "latitude", "latitud"},
		proj.DestinationCoordinatesLongitude().String(): {"lon", "lng", "long", "longitude", "longitud"},
		proj.DeliveryUnitVolume().String():              {"total volume (cm3)", "total_volume_cm3", "volume_cm3", "volumen_cm3", "volumen", "volume", "cm3", "cubic_cm", "cubic_centimeters"},
		proj.DeliveryUnitWeight().String():              {"total weight (grams)", "weight_g", "weight_grams", "peso_gramos", "grams", "gramos", "weight", "peso"},
		"price":                                         {"price", "precio", "amount", "monto", "cost", "costo"},
	}

	ordered := []string{
		"title",
		proj.ReferenceID().String(),
		proj.DestinationContactFullName().String(),
		proj.DestinationContactPhone().String(),
		proj.DestinationAddressLine1().String(),
		proj.DestinationCoordinatesLatitude().String(),
		proj.DestinationCoordinatesLongitude().String(),
		proj.DeliveryUnitVolume().String(),
		proj.DeliveryUnitWeight().String(),
		"price",
	}

	idx := buildIndex(syn)
	return &VisitFieldMapper{synonyms: syn, orderedKeys: ordered, index: idx}
}

// Map rewrites input keys -> canonical keys. Order of columns is irrelevant.
// MapMetadata returns a map where canonical keys point to the original input keys.
func (m *VisitFieldMapper) Map(input map[string]any) map[string]string {
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

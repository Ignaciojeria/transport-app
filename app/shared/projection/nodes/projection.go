package nodes

import (
	"reflect"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// Field representa un campo específico con métodos para verificación
type Field struct {
	path string
}

// Has verifica si este campo está en el mapa de campos solicitados
func (f Field) Has(requestedFields map[string]any) bool {
	_, exists := requestedFields[f.path]
	return exists
}

func (f Field) HasAnyPrefix(requestedFields map[string]any) bool {
	for field := range requestedFields {
		if field == f.path || hasPrefix(field, f.path) {
			return true
		}
	}
	return false
}

func hasPrefix(field, prefix string) bool {
	// ejemplo: prefix = "origin", field = "origin.addressInfo.contact.email"
	return len(field) > len(prefix) && field[:len(prefix)] == prefix && field[len(prefix)] == '.'
}

func init() {
	ioc.Registry(NewProjection)
}

func NewProjection() Projection {
	return Projection{}
}

// Projection representa un conjunto de campos solicitados por el cliente en la query GraphQL.
type Projection struct{}

// Métodos para campos de Location
func (p Projection) AddressInfo() Field {
	return Field{path: "addressInfo"}
}

func (p Projection) AddressLine1() Field {
	return Field{path: "addressInfo.addressLine1"}
}

func (p Projection) AddressLine2() Field {
	return Field{path: "addressInfo.addressLine2"}
}

func (p Projection) Contact() Field {
	return Field{path: "addressInfo.contact"}
}

func (p Projection) ContactEmail() Field {
	return Field{path: "addressInfo.contact.email"}
}

func (p Projection) ContactFullName() Field {
	return Field{path: "addressInfo.contact.fullName"}
}

func (p Projection) ContactNationalID() Field {
	return Field{path: "addressInfo.contact.nationalID"}
}

func (p Projection) ContactPhone() Field {
	return Field{path: "addressInfo.contact.phone"}
}

func (p Projection) ContactDocuments() Field {
	return Field{path: "addressInfo.contact.documents"}
}

func (p Projection) ContactAdditionalContactMethods() Field {
	return Field{path: "addressInfo.contact.additionalContactMethods"}
}

func (p Projection) Coordinates() Field {
	return Field{path: "addressInfo.coordinates"}
}

func (p Projection) CoordinatesLatitude() Field {
	return Field{path: "addressInfo.coordinates.latitude"}
}

func (p Projection) CoordinatesLongitude() Field {
	return Field{path: "addressInfo.coordinates.longitude"}
}

func (p Projection) CoordinatesSource() Field {
	return Field{path: "addressInfo.coordinates.source"}
}

func (p Projection) CoordinatesConfidence() Field {
	return Field{path: "addressInfo.coordinates.confidence"}
}

func (p Projection) CoordinatesConfidenceLevel() Field {
	return Field{path: "addressInfo.coordinates.confidence.level"}
}

func (p Projection) CoordinatesConfidenceMessage() Field {
	return Field{path: "addressInfo.coordinates.confidence.message"}
}

func (p Projection) CoordinatesConfidenceReason() Field {
	return Field{path: "addressInfo.coordinates.confidence.reason"}
}

func (p Projection) TimeZone() Field {
	return Field{path: "addressInfo.politicalArea.timeZone"}
}

func (p Projection) ZipCode() Field {
	return Field{path: "addressInfo.zipCode"}
}

func (p Projection) NodeInfo() Field {
	return Field{path: "nodeInfo"}
}

func (p Projection) NodeReferenceID() Field {
	return Field{path: "nodeInfo.referenceId"}
}

func (p Projection) NodeName() Field {
	return Field{path: "nodeInfo.name"}
}

func (p Projection) NodeType() Field {
	return Field{path: "nodeInfo.type"}
}

func (p Projection) NodeReferences() Field {
	return Field{path: "nodeInfo.references"}
}

// GetAllProjections devuelve un mapa con todas las proyecciones disponibles
func getAllProjections() map[string]Field {
	var p Projection
	projections := make(map[string]Field)

	// Obtener el tipo de Projection
	pType := reflect.TypeOf(p)
	pValue := reflect.ValueOf(p)

	// Iterar sobre todos los métodos de Projection
	for i := 0; i < pType.NumMethod(); i++ {
		method := pType.Method(i)

		// Solo considerar métodos que no reciben parámetros y devuelven Field
		if method.Type.NumIn() == 1 && method.Type.NumOut() == 1 {
			if method.Type.Out(0).Name() == "Field" {
				// Invocar el método para obtener el Field
				result := pValue.Method(i).Call(nil)[0].Interface().(Field)
				projections[method.Name] = result
			}
		}
	}

	return projections
}

// GetAllProjectionPaths devuelve un slice con todas las rutas de campos disponibles
func GetAllProjectionPaths() []string {
	projections := getAllProjections()
	paths := make([]string, 0, len(projections))

	for _, field := range projections {
		paths = append(paths, field.String())
	}

	return paths
}

// String devuelve la ruta del campo
func (f Field) String() string {
	return f.path
}

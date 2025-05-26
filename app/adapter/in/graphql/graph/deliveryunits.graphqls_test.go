package graph

import (
	_ "embed"
	"flag"
	"strings"
	"transport-app/app/shared/projection/deliveryunits"

	_ "embed"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

//go:embed deliveryunits.graphqls
var deliveryunitsSchema string

var verbose = flag.Bool("verbose", false, "enable verbose output")

var _ = Describe("GraphQL Schema", func() {
	var schema *ast.Schema
	var projectionPaths []string
	var schemaFields map[string]struct{}

	BeforeEach(func() {
		var err error

		// Cargar el esquema GraphQL
		schema, err = gqlparser.LoadSchema(&ast.Source{
			Name:  "deliveryunits.graphqls",
			Input: deliveryunitsSchema,
		})
		Expect(err).NotTo(HaveOccurred(), "failed to parse schema")

		// Obtener todas las rutas de proyección
		projectionPaths = deliveryunits.GetAllProjectionPaths()

		// Construir un mapa para verificar la existencia de campos en el esquema
		schemaFields = buildSchemaFieldsMap(schema)
	})

	Describe("Projections Validation", func() {
		It("should ensure all projections exist in the GraphQL schema", func() {
			// Verificar que cada proyección exista en el esquema
			for _, path := range projectionPaths {
				// Normalizar la ruta para comparación
				normalizedPath := normalizeProjectionPath(path)

				// Verificar si el campo existe en el esquema
				_, exists := schemaFields[normalizedPath]

				// Si verbose está habilitado, mostrar más información
				if *verbose {
					if exists {
						GinkgoWriter.Printf("✓ Projection '%s' exists in schema\n", path)
					} else {
						GinkgoWriter.Printf("✗ Projection '%s' (normalized as '%s') does not exist in schema\n",
							path, normalizedPath)
					}
				}

				Expect(exists).To(BeTrue(),
					"Projection path '%s' (normalized as '%s') does not exist in GraphQL schema",
					path, normalizedPath)
			}
		})

		// Esta prueba es opcional y solo imprime información, no falla si faltan proyecciones
		It("should report schema fields without projections (informational only)", func() {
			if !*verbose {
				Skip("Skipping verbose schema field check. Run with -verbose flag to see details.")
			}

			// Crear un mapa de rutas de proyección normalizadas para búsqueda rápida
			normalizedProjections := make(map[string]string)
			for _, path := range projectionPaths {
				normalizedProjections[normalizeProjectionPath(path)] = path
			}

			// Verificar cada campo del esquema
			missingCount := 0
			for schemaField := range schemaFields {
				if _, exists := normalizedProjections[schemaField]; !exists {
					GinkgoWriter.Printf("Schema field '%s' does not have a corresponding projection\n", schemaField)
					missingCount++
				}
			}

			GinkgoWriter.Printf("Found %d schema fields without corresponding projections\n", missingCount)
			// No hay expectativa aquí, ya que es solo informativo
		})
	})
})

// normalizeProjectionPath normaliza una ruta de proyección para compararla con el esquema
func normalizeProjectionPath(path string) string {
	// En este ejemplo, simplemente convertimos a minúsculas para comparación insensible a mayúsculas
	// Puedes agregar más lógica específica según tus convenciones de nomenclatura
	return strings.ToLower(path)
}

// buildSchemaFieldsMap construye un mapa de campos disponibles en el esquema GraphQL
func buildSchemaFieldsMap(schema *ast.Schema) map[string]struct{} {
	fields := make(map[string]struct{})

	// Verificar si el esquema contiene una consulta OrderConnection
	connectionType := schema.Types["DeliveryUnitsReport"]
	if connectionType != nil {
		// Si existe, buscar el tipo Order como parte de OrderConnection
		// Normalmente, OrderConnection tendría un campo "edges" con OrderEdge que a su vez tiene "node" de tipo Order
		for _, field := range connectionType.Fields {
			if field.Name == "edges" {
				edgeType := unwrapType(field.Type)
				if edgeTypeDef := schema.Types[edgeType]; edgeTypeDef != nil {
					for _, edgeField := range edgeTypeDef.Fields {
						if edgeField.Name == "node" {
							nodeType := unwrapType(edgeField.Type)
							if nodeTypeDef := schema.Types[nodeType]; nodeTypeDef != nil {
								collectFields(fields, "", nodeTypeDef, schema)
							}
						}
					}
				}
			}
			// También recopilar campos directos de OrderConnection
			fields[strings.ToLower(field.Name)] = struct{}{}
		}
	}

	// También buscar directamente el tipo Order (para casos donde no se usa el patrón de conexión)
	orderType := schema.Types["DeliveryUnitsReport"]
	if orderType != nil {
		collectFields(fields, "", orderType, schema)
	}

	// También buscar tipos relacionados que podrían ser parte de las proyecciones
	for typeName, typeDef := range schema.Types {
		if strings.HasPrefix(typeName, "DeliveryUnitsReport") && typeName != "DeliveryUnitsReport" && typeName != "DeliveryUnitsReportConnection" {
			if isObjectType(typeDef) {
				// Aquí no usamos prefijo porque estos son tipos independientes
				collectFields(fields, "", typeDef, schema)
			}
		}
	}

	// Agregar campos específicos de paginación
	fields["pageinfo"] = struct{}{}
	fields["pageinfo.hasnextpage"] = struct{}{}
	fields["pageinfo.haspreviouspage"] = struct{}{}
	fields["pageinfo.startcursor"] = struct{}{}
	fields["pageinfo.endcursor"] = struct{}{}

	return fields
}

// collectFields recopila recursivamente todos los campos de un tipo
func collectFields(fields map[string]struct{}, prefix string, objectType *ast.Definition, schema *ast.Schema) {
	for _, field := range objectType.Fields {
		// Construir el nombre completo del campo (con prefijo si existe)
		fieldName := field.Name
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		// Agregar el campo al mapa
		fields[strings.ToLower(fieldName)] = struct{}{}

		// Si el campo es un objeto, recopilar sus campos recursivamente
		fieldType := unwrapType(field.Type)
		if typeDef := schema.Types[fieldType]; typeDef != nil && isObjectType(typeDef) {
			collectFields(fields, fieldName, typeDef, schema)
		}
	}
}

// unwrapType desenrolla un tipo no nulo o lista para obtener el tipo base
func unwrapType(t *ast.Type) string {
	if t.Elem != nil {
		return unwrapType(t.Elem)
	}
	return t.NamedType
}

// isObjectType verifica si un tipo es un objeto o una interfaz
func isObjectType(typeDef *ast.Definition) bool {
	return typeDef.Kind == ast.Object || typeDef.Kind == ast.Interface
}

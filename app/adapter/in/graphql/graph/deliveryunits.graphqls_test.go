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
		schemaFields = buildDeliveryUnitsSchemaFieldsMap(schema)
	})

	Describe("Projections Validation", func() {
		It("should ensure all projections exist in the GraphQL schema", func() {
			// Verify that each projection exists in the schema
			for _, path := range projectionPaths {
				// Normalize the path for comparison
				normalizedPath := normalizeDeliveryUnitsProjectionPath(path)

				// Check if the field exists in the schema
				_, exists := schemaFields[normalizedPath]

				// If verbose is enabled, show more information
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

		// This test is optional and only prints information, doesn't fail if projections are missing
		It("should report schema fields without projections (informational only)", func() {
			if !*verbose {
				Skip("Skipping verbose schema field check. Run with -verbose flag to see details.")
			}

			// Create a map of normalized projection paths for quick lookup
			normalizedProjections := make(map[string]string)
			for _, path := range projectionPaths {
				normalizedProjections[normalizeDeliveryUnitsProjectionPath(path)] = path
			}

			// Check each schema field
			missingCount := 0
			for schemaField := range schemaFields {
				if _, exists := normalizedProjections[schemaField]; !exists {
					GinkgoWriter.Printf("Schema field '%s' does not have a corresponding projection\n", schemaField)
					missingCount++
				}
			}

			GinkgoWriter.Printf("Found %d schema fields without corresponding projections\n", missingCount)
			// No expectation here as this is informational only
		})
	})
})

// normalizeDeliveryUnitsProjectionPath normalizes a projection path for comparison with the schema
func normalizeDeliveryUnitsProjectionPath(path string) string {
	// In this example, we just convert to lowercase for case-insensitive comparison
	// You can add more specific logic based on your naming conventions
	return strings.ToLower(path)
}

// buildDeliveryUnitsSchemaFieldsMap builds a map of fields available in the GraphQL schema
func buildDeliveryUnitsSchemaFieldsMap(schema *ast.Schema) map[string]struct{} {
	fields := make(map[string]struct{})

	// Check if schema contains a DeliveryUnitsReport type
	reportType := schema.Types["DeliveryUnitsReport"]
	if reportType != nil {
		collectDeliveryUnitsFields(fields, "", reportType, schema)
	}

	return fields
}

// collectDeliveryUnitsFields recursively collects all fields from a type and its nested types
func collectDeliveryUnitsFields(fields map[string]struct{}, prefix string, typeDef *ast.Definition, schema *ast.Schema) {
	for _, field := range typeDef.Fields {
		fieldPath := prefix + field.Name
		fields[strings.ToLower(fieldPath)] = struct{}{}

		// If the field type is an object type, recursively collect its fields
		fieldType := unwrapDeliveryUnitsType(field.Type)
		if fieldTypeDef := schema.Types[fieldType]; fieldTypeDef != nil && fieldTypeDef.Kind == ast.Object {
			collectDeliveryUnitsFields(fields, fieldPath+".", fieldTypeDef, schema)
		}
	}
}

// unwrapDeliveryUnitsType unwraps a type to get its base name
func unwrapDeliveryUnitsType(t *ast.Type) string {
	if t == nil {
		return ""
	}
	if t.Elem != nil {
		return unwrapDeliveryUnitsType(t.Elem)
	}
	return t.NamedType
}

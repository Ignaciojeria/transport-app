package graph

import (
	_ "embed"
	"strings"
	"transport-app/app/shared/projection/nodes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

//go:embed nodes.graphqls
var nodesSchema string

var _ = Describe("Nodes GraphQL Schema", func() {
	var schema *ast.Schema
	var projectionPaths []string
	var schemaFields map[string]struct{}

	BeforeEach(func() {
		var err error

		// Load both GraphQL schemas
		schema, err = gqlparser.LoadSchema(
			&ast.Source{
				Name:  "deliveryunits.graphqls",
				Input: deliveryunitsSchema,
			},
			&ast.Source{
				Name:  "nodes.graphqls",
				Input: nodesSchema,
			},
		)
		Expect(err).NotTo(HaveOccurred(), "failed to parse schema")

		// Get all projection paths
		projectionPaths = nodes.GetAllProjectionPaths()

		// Build a map to verify field existence in the schema
		schemaFields = buildNodesSchemaFieldsMap(schema)
	})

	Describe("Projections Validation", func() {
		It("should ensure all projections exist in the GraphQL schema", func() {
			// Verify that each projection exists in the schema
			for _, path := range projectionPaths {
				// Normalize the path for comparison
				normalizedPath := normalizeNodesProjectionPath(path)

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
				normalizedProjections[normalizeNodesProjectionPath(path)] = path
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

// normalizeNodesProjectionPath normalizes a projection path for comparison with the schema
func normalizeNodesProjectionPath(path string) string {
	// In this example, we just convert to lowercase for case-insensitive comparison
	// You can add more specific logic based on your naming conventions
	return strings.ToLower(path)
}

// buildNodesSchemaFieldsMap builds a map of fields available in the GraphQL schema
func buildNodesSchemaFieldsMap(schema *ast.Schema) map[string]struct{} {
	fields := make(map[string]struct{})

	// Check if schema contains a NodeConnection type
	connectionType := schema.Types["NodeConnection"]
	if connectionType != nil {
		// If it exists, look for the Location type as part of NodeConnection
		// Typically, NodeConnection would have an "edges" field with NodeEdge that has "node" of type Location
		for _, field := range connectionType.Fields {
			if field.Name == "edges" {
				edgeType := unwrapNodesType(field.Type)
				if edgeTypeDef := schema.Types[edgeType]; edgeTypeDef != nil {
					for _, edgeField := range edgeTypeDef.Fields {
						if edgeField.Name == "node" {
							nodeType := unwrapNodesType(edgeField.Type)
							if nodeTypeDef := schema.Types[nodeType]; nodeTypeDef != nil {
								collectNodesFields(fields, "", nodeTypeDef, schema)
							}
						}
					}
				}
			}
			// Also collect direct fields from NodeConnection
			fields[strings.ToLower(field.Name)] = struct{}{}
		}
	}

	// Also look directly for the Location type (for cases where connection pattern is not used)
	locationType := schema.Types["Location"]
	if locationType != nil {
		collectNodesFields(fields, "", locationType, schema)
	}

	return fields
}

// collectNodesFields recursively collects all fields from a type and its nested types
func collectNodesFields(fields map[string]struct{}, prefix string, typeDef *ast.Definition, schema *ast.Schema) {
	for _, field := range typeDef.Fields {
		fieldPath := prefix + field.Name
		fields[strings.ToLower(fieldPath)] = struct{}{}

		// If the field type is an object type, recursively collect its fields
		fieldType := unwrapNodesType(field.Type)
		if fieldTypeDef := schema.Types[fieldType]; fieldTypeDef != nil && fieldTypeDef.Kind == ast.Object {
			collectNodesFields(fields, fieldPath+".", fieldTypeDef, schema)
		}
	}
}

// unwrapNodesType unwraps a type to get its base name
func unwrapNodesType(t *ast.Type) string {
	if t == nil {
		return ""
	}
	if t.Elem != nil {
		return unwrapNodesType(t.Elem)
	}
	return t.NamedType
}

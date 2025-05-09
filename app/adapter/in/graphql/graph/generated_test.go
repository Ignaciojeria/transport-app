package graph

import (
	"os"
	"time"

	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GraphQL Code Generation", func() {
	var requiredFiles []string
	var schemaFile string
	var schemaModTime time.Time

	BeforeEach(func() {
		// Lista de archivos generados que deben existir
		requiredFiles = []string{
			"generated.go",
			"model/models_gen.go",
		}

		// Archivo de esquema GraphQL
		schemaFile = "orders.graphqls"
	})

	Describe("Code Generation Validation", func() {
		It("should have all required generated files", func() {
			// Verificar que existen los archivos generados
			for _, file := range requiredFiles {
				_, err := os.Stat(file)
				Expect(err).NotTo(HaveOccurred(),
					"Generated file %s does not exist. Run 'go generate ./...' or 'gqlgen generate'", file)
			}
		})

		It("should have generated files more recent than schema", func() {
			// Obtener la información del archivo de esquema
			schemaInfo, err := os.Stat(schemaFile)
			Expect(err).NotTo(HaveOccurred(),
				"Cannot stat schema file: %s", schemaFile)

			schemaModTime = schemaInfo.ModTime()

			// Verificar que los archivos generados son más recientes que el esquema
			for _, file := range requiredFiles {
				fileInfo, err := os.Stat(file)
				if err != nil {
					// Si el archivo no existe, la prueba anterior ya lo detectará
					continue
				}

				// Comprobar que el archivo generado es más reciente que el esquema
				Expect(fileInfo.ModTime()).NotTo(BeTemporally("<", schemaModTime),
					"Generated file %s is older than schema (%s vs %s). Run 'go generate ./...' or 'gqlgen generate'",
					file, fileInfo.ModTime().Format(time.RFC3339), schemaModTime.Format(time.RFC3339))
			}
		})

		It("should have critical types in the generated model file", func() {
			// Leer el archivo de modelo generado
			modelFile := "model/models_gen.go"
			modelContent, err := os.ReadFile(modelFile)
			if err != nil {
				Skip("Cannot read model file: " + err.Error())
				return
			}

			// Lista de tipos críticos que deben existir en el modelo generado
			criticalTypes := []string{
				"type Order struct",
				"type OrderConnection struct",
				// Agrega otros tipos importantes aquí
			}

			// Verificar que los tipos críticos existen en el archivo
			modelContentStr := string(modelContent)
			for _, typeName := range criticalTypes {
				Expect(modelContentStr).To(ContainSubstring(typeName),
					"Critical type '%s' not found in generated model file. Run 'gqlgen generate'", typeName)
			}
		})
	})
})

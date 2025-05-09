package graph

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestGQLGenWasExecuted(t *testing.T) {
	// 1. Verificar que existen los archivos generados
	// Usar rutas relativas desde la ubicación de este test
	requiredFiles := []string{
		"generated.go",
		"model/models_gen.go",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Generated file %s does not exist. Run 'go generate ./...' or 'gqlgen generate'", file)
		}
	}

	// 2. Verificar que los archivos generados sean más recientes que el schema
	schemaFile := "orders.graphqls"
	schemaInfo, err := os.Stat(schemaFile)
	if err != nil {
		t.Fatalf("Cannot stat schema file: %v", err)
	}
	schemaModTime := schemaInfo.ModTime()

	// Verificar cada archivo generado (sólo si se encontró el archivo)
	for _, file := range requiredFiles {
		fileInfo, err := os.Stat(file)
		if err != nil {
			// Ya reportamos el error arriba, no lo duplicamos aquí
			continue
		}

		// Si el archivo generado es más antiguo que el schema, puede que no se haya regenerado
		if fileInfo.ModTime().Before(schemaModTime) {
			t.Errorf("Generated file %s is older than schema (%s vs %s). Run 'go generate ./...' or 'gqlgen generate'",
				file, fileInfo.ModTime().Format(time.RFC3339), schemaModTime.Format(time.RFC3339))
		}
	}

	// 3. Verificar consistencia entre schema y model generado (sólo si existe el archivo)
	modelFile := "model/models_gen.go"
	modelContent, err := os.ReadFile(modelFile)
	if err != nil {
		t.Logf("Cannot read model file: %v", err)
		return
	}

	// Verificar que existen ciertos tipos críticos en el modelo generado
	criticalTypes := []string{
		"type Order struct",
		"type OrderConnection struct",
		// Agrega otros tipos importantes aquí
	}

	modelContentStr := string(modelContent)
	for _, typeName := range criticalTypes {
		if !strings.Contains(modelContentStr, typeName) {
			t.Errorf("Critical type '%s' not found in generated model file. Run 'gqlgen generate'", typeName)
		}
	}
}

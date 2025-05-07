package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func mapDocuments(docs []domain.Document) table.JSONReference {
	// Crear un slice para mapear los documentos
	mapped := make(table.JSONReference, len(docs))

	// Iterar sobre los documentos y mapearlos
	for i, d := range docs {
		mapped[i] = table.Reference{
			Type:  d.Type,
			Value: d.Value,
		}
	}

	return mapped
}

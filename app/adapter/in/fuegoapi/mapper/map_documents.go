package mapper

import "transport-app/app/domain"

func MapDocumentsToDomain(docs []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.Document {
	mapped := make([]domain.Document, len(docs))
	for i, doc := range docs {
		mapped[i] = domain.Document{
			Type:  doc.Type,
			Value: doc.Value,
		}
	}
	return mapped
}

func MapDocumentsFromDomain(docs []domain.Document) []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
} {
	mapped := make([]struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}, len(docs))

	for i, doc := range docs {
		mapped[i] = struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			Type:  doc.Type,
			Value: doc.Value,
		}
	}

	return mapped
}

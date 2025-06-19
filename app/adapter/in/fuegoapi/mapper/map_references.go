package mapper

import "transport-app/app/domain"

func MapReferencesToDomain(refs []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.Reference {
	mapped := make([]domain.Reference, len(refs))
	for i, ref := range refs {
		mapped[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mapped
}

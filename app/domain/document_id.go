package domain

type DocumentID string

func (id DocumentID) IsZero() bool {
	return string(id) == ""
}

func (id DocumentID) Equals(other string) bool {
	return string(id) == other
}

func (id DocumentID) ShouldUpdate(existing string) bool {
	return !id.IsZero() && !id.Equals(existing)
}

package domain

type ReferenceID string

func (r ReferenceID) String() string {
	return string(r)
}

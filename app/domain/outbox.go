package domain

type Outbox struct {
	ID           int64
	CreatedAt    string
	UpdatedAt    string
	Organization Organization
	Attributes   map[string]string
	Payload      []byte
	Status       string
}

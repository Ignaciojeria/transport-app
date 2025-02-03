package domain

type Outbox struct {
	CreatedAt    string
	UpdatedAt    string
	Organization Organization
	Attributes   map[string]string
	Payload      []byte
	Status       string
}

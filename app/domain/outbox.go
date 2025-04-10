package domain

type Outbox struct {
	CreatedAt  string
	UpdatedAt  string
	Attributes map[string]string
	Payload    []byte
	Status     string
}

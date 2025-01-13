package domain

type Outbox struct {
	ID           int64
	ReferenceID  string
	EntityType   string
	EventType    string
	CreatedAt    string
	UpdatedAt    string
	Organization Organization
	Payload      []byte
	Processed    bool
}

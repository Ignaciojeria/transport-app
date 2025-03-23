package domain

import "github.com/google/uuid"

type Contact struct {
	ID           int64
	Organization Organization `json:"organization"`
	FullName     string       `json:"fullName"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	NationalID   string       `json:"nationalID"`
	Documents    []Document   `json:"documents"`
}

func (c Contact) ReferenceID() string {
	var key string

	switch {
	case c.NationalID != "":
		key = c.NationalID
	case c.Email != "":
		key = c.Email
	case c.Phone != "":
		key = c.Phone
	default:
		// Caso extremo: no hay nada identificador. Se genera un UUID para evitar colisiones.
		key = uuid.NewString()
	}

	return Hash(c.Organization, key)
}

func (c Contact) UpdateIfChanged(newContact Contact) Contact {
	// Copiamos la instancia actual

	// Actualizar FullName
	if newContact.ID != 0 {
		c.ID = newContact.ID
	}

	// Actualizar FullName
	if newContact.FullName != "" {
		c.FullName = newContact.FullName
	}

	// Actualizar Email
	if newContact.Email != "" {
		c.Email = newContact.Email
	}

	// Actualizar Phone
	if newContact.Phone != "" {
		c.Phone = newContact.Phone
	}

	// Actualizar NationalID
	if newContact.NationalID != "" {
		c.NationalID = newContact.NationalID
	}

	// Actualizar Documents
	if len(newContact.Documents) > 0 {
		c.Documents = newContact.Documents
	}

	return c
}

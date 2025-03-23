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

func (c Contact) UpdateIfChanged(newContact Contact) (Contact, bool) {
	updated := false

	if newContact.FullName != "" && newContact.FullName != c.FullName {
		c.FullName = newContact.FullName
		updated = true
	}
	if newContact.Email != "" && newContact.Email != c.Email {
		c.Email = newContact.Email
		updated = true
	}
	if newContact.Phone != "" && newContact.Phone != c.Phone {
		c.Phone = newContact.Phone
		updated = true
	}
	if newContact.NationalID != "" && newContact.NationalID != c.NationalID {
		c.NationalID = newContact.NationalID
		updated = true
	}
	if len(newContact.Documents) > 0 {
		c.Documents = newContact.Documents
		updated = true
	}

	return c, updated
}

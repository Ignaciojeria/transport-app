package domain

import "github.com/google/uuid"

type Contact struct {
	ID                       int64
	Organization             Organization
	FullName                 string
	PrimaryEmail             string
	PrimaryPhone             string
	NationalID               string
	Documents                []Document
	AdditionalContactMethods []ContactMethod
}

type ContactMethod struct {
	Type  string `json:"type"`  // Ej: "email", "phone", "whatsapp"
	Value string `json:"value"` // Ej: "ejemplo@correo.com"
}

func (c Contact) DocID() string {
	var key string
	switch {
	case c.PrimaryEmail != "":
		key = c.PrimaryEmail
	case c.PrimaryPhone != "":
		key = c.PrimaryPhone
	case c.NationalID != "":
		key = c.NationalID
	default:
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
	if newContact.PrimaryEmail != "" && newContact.PrimaryEmail != c.PrimaryEmail {
		c.PrimaryEmail = newContact.PrimaryEmail
		updated = true
	}
	if newContact.PrimaryPhone != "" && newContact.PrimaryPhone != c.PrimaryPhone {
		c.PrimaryPhone = newContact.PrimaryPhone
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

package domain

import (
	"context"
)

type Contact struct {
	FullName                 string
	PrimaryEmail             string
	PrimaryPhone             string
	NationalID               string
	Documents                []Document
	AdditionalContactMethods []ContactMethod
}

func (c Contact) Equals(ctx context.Context, other Contact) bool {
	return c.DocID(ctx) == other.DocID(ctx)
}

type ContactMethod struct {
	Type  string `json:"type"`  // Ej: "email", "phone", "whatsapp"
	Value string `json:"value"` // Ej: "ejemplo@correo.com"
}

// UpdateIfChange actualiza los campos de ContactMethod solo si son diferentes
// Devuelve el método de contacto actualizado y un booleano que indica si se realizó algún cambio
func (c ContactMethod) UpdateIfChange(newContact ContactMethod) (ContactMethod, bool) {
	updated := c
	changed := false

	// Verificar y actualizar Type si es diferente y no está vacío
	if newContact.Type != "" && newContact.Type != c.Type {
		updated.Type = newContact.Type
		changed = true
	}

	// Verificar y actualizar Value si es diferente y no está vacío
	if newContact.Value != "" && newContact.Value != c.Value {
		updated.Value = newContact.Value
		changed = true
	}

	return updated, changed
}

func (c Contact) DocID(ctx context.Context) DocumentID {
	var key string
	switch {
	case c.PrimaryEmail != "":
		key = c.PrimaryEmail
	case c.PrimaryPhone != "":
		key = c.PrimaryPhone
	case c.NationalID != "":
		key = c.NationalID
	default:
		key = ""
	}
	return HashByTenant(ctx, key)
}

func (c Contact) UpdateIfChanged(newContact Contact) (Contact, bool) {
	updated := c
	changed := false

	if newContact.FullName != "" && newContact.FullName != c.FullName {
		updated.FullName = newContact.FullName
		changed = true
	}
	if newContact.PrimaryEmail != "" && newContact.PrimaryEmail != c.PrimaryEmail {
		updated.PrimaryEmail = newContact.PrimaryEmail
		changed = true
	}
	if newContact.PrimaryPhone != "" && newContact.PrimaryPhone != c.PrimaryPhone {
		updated.PrimaryPhone = newContact.PrimaryPhone
		changed = true
	}
	if newContact.NationalID != "" && newContact.NationalID != c.NationalID {
		updated.NationalID = newContact.NationalID
		changed = true
	}

	// Actualizar documentos usando UpdateIfChange
	if len(newContact.Documents) > 0 {
		// Crear un mapa de documentos existentes por Type para facilitar la búsqueda
		docMap := make(map[string]int) // Type -> índice en el slice
		for i, doc := range c.Documents {
			docMap[doc.Type] = i
		}

		// Copiar los documentos actuales como base
		updatedDocs := make([]Document, len(c.Documents))
		copy(updatedDocs, c.Documents)

		docsChanged := false

		// Procesar cada nuevo documento
		for _, newDoc := range newContact.Documents {
			// No considerar documentos completamente vacíos
			if newDoc.Type == "" && newDoc.Value == "" {
				continue
			}

			// Buscar por Type para mantener consistencia con otros métodos
			if idx, exists := docMap[newDoc.Type]; exists {
				// Si el documento existe, intentar actualizarlo
				updatedDoc, docChanged := updatedDocs[idx].UpdateIfChange(newDoc)
				if docChanged {
					updatedDocs[idx] = updatedDoc
					docsChanged = true
				}
			} else {
				// Si el documento no existe, agregarlo
				updatedDocs = append(updatedDocs, newDoc)
				docsChanged = true
			}
		}

		if docsChanged {
			updated.Documents = updatedDocs
			changed = true
		}
	}

	// Actualizar métodos de contacto adicionales
	if len(newContact.AdditionalContactMethods) > 0 {
		// Crear un mapa de métodos existentes por Type para facilitar la búsqueda
		methodMap := make(map[string]int) // Type -> índice en el slice
		for i, method := range c.AdditionalContactMethods {
			methodMap[method.Type] = i
		}

		// Copiar los métodos actuales como base
		updatedMethods := make([]ContactMethod, len(c.AdditionalContactMethods))
		copy(updatedMethods, c.AdditionalContactMethods)

		methodsChanged := false

		// Procesar cada nuevo método de contacto
		for _, newMethod := range newContact.AdditionalContactMethods {
			// No considerar métodos completamente vacíos
			if newMethod.Type == "" && newMethod.Value == "" {
				continue
			}

			// Buscar por Type para mantener el comportamiento similar a References
			if idx, exists := methodMap[newMethod.Type]; exists {
				// Si el método existe, intentar actualizarlo
				updatedMethod, methodChanged := updatedMethods[idx].UpdateIfChange(newMethod)
				if methodChanged {
					updatedMethods[idx] = updatedMethod
					methodsChanged = true
				}
			} else {
				// Si el método no existe, agregarlo
				updatedMethods = append(updatedMethods, newMethod)
				methodsChanged = true
			}
		}

		if methodsChanged {
			updated.AdditionalContactMethods = updatedMethods
			changed = true
		}
	}

	return updated, changed
}

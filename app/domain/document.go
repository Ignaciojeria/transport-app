package domain

type Document struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

// UpdateIfChange actualiza los campos de Document solo si son diferentes
// Devuelve el documento actualizado y un booleano que indica si se realizó algún cambio
func (d Document) UpdateIfChange(newDoc Document) (Document, bool) {
	updated := d
	changed := false

	// Verificar y actualizar Type si es diferente y no está vacío
	if newDoc.Type != "" && newDoc.Type != d.Type {
		updated.Type = newDoc.Type
		changed = true
	}

	// Verificar y actualizar Value si es diferente y no está vacío
	if newDoc.Value != "" && newDoc.Value != d.Value {
		updated.Value = newDoc.Value
		changed = true
	}

	return updated, changed
}

// Función auxiliar para comparar arreglos de documentos
func compareDocuments(oldDocs, newDocs []Document) bool {
	if len(oldDocs) != len(newDocs) {
		return false
	}
	for i := range oldDocs {
		if oldDocs[i] != newDocs[i] {
			return false
		}
	}
	return true
}

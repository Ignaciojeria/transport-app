package domain

type NodeType struct {
	Organization Organization
	Value        string
}

// DocID genera un identificador único para NodeType basado en Organization y Value.
// Si Value está vacío, se usa un string vacío como clave.
func (nt NodeType) DocID() DocumentID {
	if nt.Value == "" {
		return ""
	}
	return DocumentID(Hash(nt.Organization, nt.Value))
}

func (nt NodeType) UpdateIfChange(newNodeType NodeType) (NodeType, bool) {
	updated := nt
	changed := false

	// Actualizar Value si es diferente y no está vacío
	if newNodeType.Value != "" && newNodeType.Value != nt.Value {
		updated.Value = newNodeType.Value
		changed = true
	}

	return updated, changed
}

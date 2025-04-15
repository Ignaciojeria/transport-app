package domain

import "context"

type NodeType struct {
	Value string
}

// DocID genera un identificador único para NodeType basado en Organization y Value.
// Si Value está vacío, se usa un string vacío como clave.
func (nt NodeType) DocID(ctx context.Context) DocumentID {
	if nt.Value == "" {
		return ""
	}
	return DocumentID(HashByTenant(ctx, nt.Value))
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

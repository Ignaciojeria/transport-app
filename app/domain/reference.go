package domain

import "context"

type Reference struct {
	Type  string
	Value string
}

func (r Reference) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, r.Type, r.Value)
}

// UpdateIfChange actualiza los campos de Reference solo si son diferentes
// Devuelve la referencia actualizada y un booleano que indica si se realizó algún cambio
func (r Reference) UpdateIfChange(newRef Reference) (Reference, bool) {
	updated := r
	changed := false

	// Verificar y actualizar Type si es diferente y no está vacío
	if newRef.Type != "" && newRef.Type != r.Type {
		updated.Type = newRef.Type
		changed = true
	}

	// Verificar y actualizar Value si es diferente y no está vacío
	if newRef.Value != "" && newRef.Value != r.Value {
		updated.Value = newRef.Value
		changed = true
	}

	return updated, changed
}

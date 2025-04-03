package domain

type NodeInfo struct {
	ID           int64
	ReferenceID  ReferenceID  `json:"referenceID"`
	Organization Organization `json:"organization"`
	Name         string       `json:"name"`
	NodeType     NodeType     `json:"type"`
	Contact      Contact      `json:"contact"`
	References   []Reference  `json:"references"`
	AddressInfo  AddressInfo  `json:"addressInfo"`
	AddressLine2 string       `json:"addressLine2"`
	AddressLine3 string       `json:"addressLine3"`
}

func (n NodeInfo) DocID() DocumentID {
	if n.ReferenceID == "" {
		return ""
	}
	return Hash(n.Organization, string(n.ReferenceID))
}

func (n NodeInfo) UpdateIfChanged(newNode NodeInfo) (NodeInfo, bool) {
	updated := n
	changed := false

	// Actualizar campos básicos
	if newNode.ReferenceID != "" && newNode.ReferenceID != n.ReferenceID {
		updated.ReferenceID = newNode.ReferenceID
		changed = true
	}

	if newNode.Name != "" && newNode.Name != n.Name {
		updated.Name = newNode.Name
		changed = true
	}

	// Actualizar las Referencias usando DocID como identificador
	if len(newNode.References) > 0 {
		// Crear un mapa de referencias existentes por Type para rápida búsqueda
		refMap := make(map[string]int) // Type -> índice en el slice
		for i, ref := range n.References {
			refMap[ref.Type] = i
		}

		// Copiar las referencias actuales
		updatedRefs := make([]Reference, len(n.References))
		copy(updatedRefs, n.References)

		refsChanged := false

		// Procesar cada nueva referencia
		for _, newRef := range newNode.References {
			// No considerar referencias completamente vacías
			if newRef.Type == "" && newRef.Value == "" {
				continue
			}

			// Buscar primero por Type para mantener el comportamiento existente
			if idx, exists := refMap[newRef.Type]; exists {
				// Si la referencia existe, intentar actualizarla
				updatedRef, refChanged := updatedRefs[idx].UpdateIfChange(newRef)
				if refChanged {
					updatedRefs[idx] = updatedRef
					refsChanged = true
				}
			} else {
				// Si la referencia no existe, agregarla
				updatedRefs = append(updatedRefs, newRef)
				refsChanged = true
			}
		}

		if refsChanged {
			updated.References = updatedRefs
			changed = true
		}
	}

	// Actualizar campos AddressLine que se movieron de AddressInfo a NodeInfo
	if newNode.AddressLine2 != "" && newNode.AddressLine2 != n.AddressLine2 {
		updated.AddressLine2 = newNode.AddressLine2
		changed = true
	}

	if newNode.AddressLine3 != "" && newNode.AddressLine3 != n.AddressLine3 {
		updated.AddressLine3 = newNode.AddressLine3
		changed = true
	}

	return updated, changed
}

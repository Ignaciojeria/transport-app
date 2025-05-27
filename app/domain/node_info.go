package domain

import (
	"context"
	"errors"
)

type NodeInfo struct {
	ReferenceID ReferenceID
	Name        string
	NodeType    NodeType
	References  []Reference
	AddressInfo AddressInfo
	Headers     Headers
}

func (n NodeInfo) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, string(n.ReferenceID))
}

func search_node_headers_by_node_doc(ctx context.Context, docID DocumentID) (Headers, error) {
	if docID == "" {
		return Headers{}, errors.New("document ID cannot be empty")
	}

	// Create a new NodeInfo with the document ID
	nodeInfo := NodeInfo{
		ReferenceID: ReferenceID(docID),
	}

	// Get the headers from the node info
	return nodeInfo.Headers, nil
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

	// Actualizar Headers si han cambiado
	if !newNode.Headers.IsEmpty() {
		updatedHeaders, headersChanged := n.Headers.UpdateIfChanged(newNode.Headers)
		if headersChanged {
			updated.Headers = updatedHeaders
			changed = true
		}
	}

	return updated, changed
}

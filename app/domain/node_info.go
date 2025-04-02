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

func (n NodeInfo) DocID() string {
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

	// Actualizar NodeType
	if newNode.NodeType.Value != "" && newNode.NodeType.Value != n.NodeType.Value {
		updated.NodeType.Value = newNode.NodeType.Value
		changed = true
	}

	if newNode.NodeType.ID != 0 && newNode.NodeType.ID != n.NodeType.ID {
		updated.NodeType.ID = newNode.NodeType.ID
		changed = true
	}

	// Actualizar References
	if len(newNode.References) > 0 {
		updated.References = newNode.References
		changed = true
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

	// Reutilizar los métodos UpdateIfChanged para objetos complejos
	// AddressInfo
	updatedAddressInfo, addressChanged := n.AddressInfo.UpdateIfChanged(newNode.AddressInfo)
	if addressChanged {
		updated.AddressInfo = updatedAddressInfo
		changed = true
	}

	// Contact
	updatedContact, contactChanged := n.Contact.UpdateIfChanged(newNode.Contact)
	if contactChanged {
		updated.Contact = updatedContact
		changed = true
	}

	return updated, changed
}

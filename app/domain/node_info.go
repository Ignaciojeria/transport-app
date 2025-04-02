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

func (n NodeInfo) UpdateIfChanged(newNode NodeInfo) NodeInfo {
	// Actualizar ReferenceID
	if newNode.ID != 0 {
		n.ID = newNode.ID
	}
	if newNode.ReferenceID != "" && n.ReferenceID != newNode.ReferenceID {
		n.ReferenceID = newNode.ReferenceID
	}
	// Actualizar Name
	if newNode.Name != "" {
		n.Name = newNode.Name
	}
	// Actualizar Type
	if newNode.NodeType.Value != "" && n.NodeType.Value != newNode.NodeType.Value {
		n.NodeType.Value = newNode.NodeType.Value
	}
	// Actualizar NodeReferences
	if len(newNode.References) > 0 {
		n.References = newNode.References
	}
	if newNode.AddressInfo.ID != 0 {
		n.AddressInfo.ID = newNode.AddressInfo.ID
	}
	if newNode.Contact.ID != 0 {
		n.Contact.ID = newNode.Contact.ID
	}
	if newNode.Contact.ID != 0 {
		n.Contact.ID = newNode.Contact.ID
	}
	if newNode.NodeType.ID != 0 {
		n.NodeType.ID = newNode.NodeType.ID
	}
	n.Organization = newNode.Organization
	//n.AddressInfo = n.AddressInfo.UpdateIfChanged(newNode.AddressInfo)
	//n.NodeType = n.NodeType.UpdateIfChanged(newNode.NodeType)
	return n
}

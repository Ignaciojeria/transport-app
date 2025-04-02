package domain

type Operator struct {
	ID           int64
	ReferenceID  string
	Organization Organization
	OriginNode   NodeInfo
	Contact      Contact `json:"contact"`
	Type         string  `json:"type"`
}

func (o Operator) UpdateIfChanged(newOperator Operator) Operator {
	// Copiamos la instancia actual

	if newOperator.ID != 0 {
		o.ID = newOperator.ID
	}

	// Actualizar campos primitivos solo si tienen valores no vacíos
	if newOperator.ReferenceID != "" {
		o.ReferenceID = newOperator.ReferenceID
	}
	if newOperator.Type != "" {
		o.Type = newOperator.Type
	}

	// Actualizar Organization si tiene valores nuevos
	if newOperator.Organization.ID != 0 {
		o.Organization = newOperator.Organization
	}

	o.OriginNode, _ = o.OriginNode.UpdateIfChanged(newOperator.OriginNode)

	o.Contact, _ = o.Contact.UpdateIfChanged(newOperator.Contact)

	return o
}

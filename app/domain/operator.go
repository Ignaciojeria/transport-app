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
	updatedOperator := o // Copiamos la instancia actual

	// Actualizar campos primitivos solo si tienen valores no vacíos
	if newOperator.ReferenceID != "" {
		updatedOperator.ReferenceID = newOperator.ReferenceID
	}
	if newOperator.Type != "" {
		updatedOperator.Type = newOperator.Type
	}

	// Actualizar Organization si tiene valores nuevos
	if newOperator.Organization.OrganizationCountryID != 0 {
		updatedOperator.Organization = newOperator.Organization
	}

	updatedOperator.OriginNode = updatedOperator.OriginNode.UpdateIfChanged(newOperator.OriginNode)

	updatedOperator.Contact = updatedOperator.Contact.UpdateIfChanged(newOperator.Contact)

	return updatedOperator
}

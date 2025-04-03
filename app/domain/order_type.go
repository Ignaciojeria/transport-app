package domain

type OrderType struct {
	Organization Organization
	Type         string `json:"type"`
	Description  string `json:"description"`
}

func (ot OrderType) DocID() DocumentID {
	return Hash(ot.Organization, ot.Type)
}

func (ot OrderType) UpdateIfChanged(newOrderType OrderType) (OrderType, bool) {
	changed := false

	if newOrderType.Type != "" && newOrderType.Type != ot.Type {
		ot.Type = newOrderType.Type
		changed = true
	}
	if newOrderType.Description != "" && newOrderType.Description != ot.Description {
		ot.Description = newOrderType.Description
		changed = true
	}
	return ot, changed
}

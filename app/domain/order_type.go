package domain

import "context"

type OrderType struct {
	Type        string
	Description string
}

func (ot OrderType) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, ot.Type)
}

func (ot OrderType) UpdateIfChanged(newOrderType OrderType) (OrderType, bool) {
	changed := false

	if newOrderType.Description != "" && newOrderType.Description != ot.Description {
		ot.Description = newOrderType.Description
		changed = true
	}
	return ot, changed
}

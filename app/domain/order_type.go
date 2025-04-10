package domain

import "context"

type OrderType struct {
	Type        string
	Description string
}

func (ot OrderType) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, ot.Type)
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

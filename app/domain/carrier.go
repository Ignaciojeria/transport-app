package domain

import "context"

type Carrier struct {
	Name       string
	NationalID string
}

func (c Carrier) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, c.NationalID)
}

func (c Carrier) UpdateIfChanged(newCarrier Carrier) (Carrier, bool) {
	updated := c
	changed := false

	if newCarrier.Name != "" && newCarrier.Name != c.Name {
		updated.Name = newCarrier.Name
		changed = true
	}

	return updated, changed
}

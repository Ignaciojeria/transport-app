package domain

import "context"

type NonDeliveryReason struct {
	ReferenceID string
	Reason      string
	Details     string
}

func (r NonDeliveryReason) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, r.ReferenceID)
}

func (r NonDeliveryReason) UpdateIfChanged(newR NonDeliveryReason) (NonDeliveryReason, bool) {
	updated := r
	changed := false

	if newR.ReferenceID != "" && newR.ReferenceID != r.ReferenceID {
		updated.ReferenceID = newR.ReferenceID
		changed = true
	}

	if newR.Reason != "" && newR.Reason != r.Reason {
		updated.Reason = newR.Reason
		changed = true
	}

	return updated, changed
}

func (r NonDeliveryReason) IsEmpty() bool {
	return r.ReferenceID == "" && r.Reason == "" && r.Details == ""
}

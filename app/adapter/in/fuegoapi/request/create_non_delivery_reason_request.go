package request

import "transport-app/app/domain"

type CreateNonDeliveryReasonRequest struct {
	ReferenceID string
	Reason      string
}

func (r CreateNonDeliveryReasonRequest) Map() domain.NonDeliveryReason {
	return domain.NonDeliveryReason{
		ReferenceID: r.ReferenceID,
		Reason:      r.Reason,
	}
}

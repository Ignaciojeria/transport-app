package domain

type PlanType struct {
	Organization
	Value string
}

func (pt PlanType) DocID() DocumentID {
	return Hash(pt.Organization, pt.Value)
}

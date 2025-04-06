package domain

type PlanningStatus struct {
	Organization
	Value string
}

func (ps PlanningStatus) DocID() DocumentID {
	return Hash(ps.Organization, ps.Value)
}

package domain

type Carrier struct {
	Organization Organization `json:"organization"`
	Name         string       `json:"name"`
	NationalID   string       `json:"nationalID"`
}

func (c Carrier) DocID() DocumentID {
	return Hash(c.Organization, c.NationalID)
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

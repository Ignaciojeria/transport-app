package domain

type Headers struct {
	Organization Organization
	ID           int64
	Consumer     string `json:"consumer"`
	Commerce     string `json:"commerce"`
}

func (h Headers) ReferenceID() string {
	return Hash(h.Organization, h.Commerce, h.Consumer)
}

func (c Headers) UpdateIfChanged(newHeaders Headers) Headers {
	if newHeaders.ID != 0 {
		c.ID = newHeaders.ID
	}
	if newHeaders.Consumer != "" {
		c.Consumer = newHeaders.Consumer
	}
	if newHeaders.Commerce != "" {
		c.Commerce = newHeaders.Commerce
	}
	if newHeaders.Organization.ID != 0 {
		c.Organization = newHeaders.Organization
	}
	return c
}

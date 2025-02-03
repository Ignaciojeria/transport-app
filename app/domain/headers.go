package domain

type Headers struct {
	Organization Organization
	ID           int64
	Consumer     string `json:"consumer"`
	Commerce     string `json:"commerce"`
}

func (c Headers) UpdateIfChanged(newHeaders Headers) Headers {
	if newHeaders.ID != 0 {
		c.ID = newHeaders.ID
	}
	if newHeaders.Consumer != "" {
		c.Consumer = newHeaders.Consumer
	}
	if newHeaders.Commerce != "" {
		c.Consumer = newHeaders.Commerce
	}
	if newHeaders.Organization.OrganizationCountryID != 0 {
		c.Organization = newHeaders.Organization
	}
	return c
}

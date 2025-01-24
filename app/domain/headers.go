package domain

type Headers struct {
	Organization Organization
	ID           int64
	Consumer     string `json:"consumer"`
	Commerce     string `json:"commerce"`
}

func (c Headers) UpdateIfChanged(newConsumer Headers) Headers {
	if newConsumer.Consumer != "" {
		c.Consumer = newConsumer.Consumer
	}
	if newConsumer.Commerce != "" {
		c.Consumer = newConsumer.Commerce
	}
	if newConsumer.Organization.OrganizationCountryID != 0 {
		c.Organization = newConsumer.Organization
	}
	return c
}

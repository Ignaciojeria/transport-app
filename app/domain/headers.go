package domain

type Headers struct {
	Organization Organization
	Consumer     string `json:"consumer"`
	Commerce     string `json:"commerce"`
}

const ConsumerKey = "consumer"
const CommerceKey = "consumer"

func (h Headers) DocID() DocumentID {
	return Hash(h.Organization, h.Commerce, h.Consumer)
}

func (c Headers) UpdateIfChanged(newHeaders Headers) Headers {
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

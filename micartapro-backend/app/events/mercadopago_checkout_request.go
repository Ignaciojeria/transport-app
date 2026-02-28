package events

type MercadoPagoCheckoutRequest struct {
	Items  []MercadoPagoCheckoutItem `json:"items"`
	Totals struct {
		Subtotal    int    `json:"subtotal"`
		DeliveryFee int    `json:"deliveryFee"`
		Total       int    `json:"total"`
		Currency    string `json:"currency"`
	} `json:"totals"`
	BusinessInfo struct {
		BusinessName string `json:"businessName"`
		Whatsapp     string `json:"whatsapp"`
	} `json:"businessInfo"`
	Fulfillment struct {
		Type string `json:"type"`
	} `json:"fulfillment"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MercadoPagoCheckoutItem struct {
	ProductName string  `json:"productName"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	UnitPrice   int     `json:"unitPrice"`
}

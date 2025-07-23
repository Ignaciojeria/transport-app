package request

import (
	"context"
	"time"
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

// UpsertOrderGroupBy representa el agrupamiento del pedido
type UpsertOrderGroupBy struct {
	Type  string `json:"type" example:"customerOrder"`
	Value string `json:"value" example:"1234567890"`
}

// UpsertOrderTimeRange representa un rango de tiempo
type UpsertOrderTimeRange struct {
	EndTime   string `json:"endTime" example:"09:00"`
	StartTime string `json:"startTime" example:"19:00"`
}

// Map convierte a domain.TimeRange
func (t UpsertOrderTimeRange) Map() domain.TimeRange {
	return domain.TimeRange{
		StartTime: t.StartTime,
		EndTime:   t.EndTime,
	}
}

// UpsertOrderCollectAvailabilityDate representa la fecha de disponibilidad de recolección
type UpsertOrderCollectAvailabilityDate struct {
	Date      string               `json:"date" example:"2025-03-30"`
	TimeRange UpsertOrderTimeRange `json:"timeRange"`
}

// Map convierte a domain.CollectAvailabilityDate
func (c UpsertOrderCollectAvailabilityDate) Map() domain.CollectAvailabilityDate {
	return domain.CollectAvailabilityDate{
		Date:      c.parseDate(c.Date),
		TimeRange: c.TimeRange.Map(),
	}
}

// parseDate convierte string a time.Time
func (c UpsertOrderCollectAvailabilityDate) parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

// UpsertOrderContactMethod representa un método de contacto adicional
type UpsertOrderContactMethod struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Map convierte a domain.ContactMethod
func (c UpsertOrderContactMethod) Map() domain.ContactMethod {
	return domain.ContactMethod{
		Type:  c.Type,
		Value: c.Value,
	}
}

// UpsertOrderDocument representa un documento
type UpsertOrderDocument struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Map convierte a domain.Document
func (d UpsertOrderDocument) Map() domain.Document {
	return domain.Document{
		Type:  d.Type,
		Value: d.Value,
	}
}

// UpsertOrderContact representa la información de contacto
type UpsertOrderContact struct {
	AdditionalContactMethods []UpsertOrderContactMethod `json:"additionalContactMethods"`
	Email                    string                     `json:"email"`
	Phone                    string                     `json:"phone"`
	NationalID               string                     `json:"nationalID"`
	Documents                []UpsertOrderDocument      `json:"documents"`
	FullName                 string                     `json:"fullName"`
}

// Map convierte a domain.Contact
func (c UpsertOrderContact) Map() domain.Contact {
	additionalMethods := make([]domain.ContactMethod, len(c.AdditionalContactMethods))
	for i, method := range c.AdditionalContactMethods {
		additionalMethods[i] = method.Map()
	}

	documents := make([]domain.Document, len(c.Documents))
	for i, doc := range c.Documents {
		documents[i] = doc.Map()
	}

	return domain.Contact{
		FullName:                 c.FullName,
		PrimaryEmail:             c.Email,
		PrimaryPhone:             c.Phone,
		NationalID:               c.NationalID,
		Documents:                documents,
		AdditionalContactMethods: additionalMethods,
	}
}

// UpsertOrderConfidence representa la confianza de las coordenadas
type UpsertOrderConfidence struct {
	Level   float64 `json:"level" example:"0.1"`
	Message string  `json:"message" example:"DISTRICT_CENTROID"`
	Reason  string  `json:"reason" example:"PROVIDER_RESULT_OUT_OF_DISTRICT"`
}

// Map convierte a domain.CoordinatesConfidence
func (c UpsertOrderConfidence) Map() domain.CoordinatesConfidence {
	return domain.CoordinatesConfidence{
		Level:   c.Level,
		Message: c.Message,
		Reason:  c.Reason,
	}
}

// UpsertOrderCoordinates representa las coordenadas
type UpsertOrderCoordinates struct {
	Latitude   float64               `json:"latitude" example:"-33.5147889"`
	Longitude  float64               `json:"longitude" example:"-70.6130425"`
	Source     string                `json:"source" example:"GOOGLE_MAPS"`
	Confidence UpsertOrderConfidence `json:"confidence"`
}

// Map convierte a domain.Coordinates
func (c UpsertOrderCoordinates) Map() domain.Coordinates {
	return domain.Coordinates{
		Point:      orb.Point{c.Longitude, c.Latitude},
		Source:     c.Source,
		Confidence: c.Confidence.Map(),
	}
}

// UpsertOrderPoliticalArea representa el área política
type UpsertOrderPoliticalArea struct {
	Code            string                `json:"code" example:"cl-rm-la-florida"`
	AdminAreaLevel1 string                `json:"adminAreaLevel1" example:"region metropolitana de santiago"`
	AdminAreaLevel2 string                `json:"adminAreaLevel2" example:"santiago"`
	AdminAreaLevel3 string                `json:"adminAreaLevel3" example:"la florida"`
	AdminAreaLevel4 string                `json:"adminAreaLevel4" example:""`
	TimeZone        string                `json:"timeZone" example:"America/Santiago"`
	Confidence      UpsertOrderConfidence `json:"confidence"`
}

// Map convierte a domain.PoliticalArea
func (p UpsertOrderPoliticalArea) Map() domain.PoliticalArea {
	return domain.PoliticalArea{
		Code:            p.Code,
		AdminAreaLevel1: p.AdminAreaLevel1,
		AdminAreaLevel2: p.AdminAreaLevel2,
		AdminAreaLevel3: p.AdminAreaLevel3,
		AdminAreaLevel4: p.AdminAreaLevel4,
		TimeZone:        p.TimeZone,
		Confidence:      p.Confidence.Map(),
	}
}

// UpsertOrderAddressInfo representa la información de dirección
type UpsertOrderAddressInfo struct {
	AddressLine1  string                   `json:"addressLine1" example:"Inglaterra 59"`
	AddressLine2  string                   `json:"addressLine2" example:"Piso 2214"`
	Contact       UpsertOrderContact       `json:"contact"`
	Coordinates   UpsertOrderCoordinates   `json:"coordinates"`
	PoliticalArea UpsertOrderPoliticalArea `json:"politicalArea"`
	ZipCode       string                   `json:"zipCode" example:"7500000"`
}

// Map convierte a domain.AddressInfo
func (a UpsertOrderAddressInfo) Map() domain.AddressInfo {
	return domain.AddressInfo{
		AddressLine1:  a.AddressLine1,
		AddressLine2:  a.AddressLine2,
		Contact:       a.Contact.Map(),
		Coordinates:   a.Coordinates.Map(),
		PoliticalArea: a.PoliticalArea.Map(),
		ZipCode:       a.ZipCode,
	}
}

// UpsertOrderNodeInfo representa la información del nodo
type UpsertOrderNodeInfo struct {
	ReferenceID string `json:"referenceID"`
}

// UpsertOrderDestination representa el destino
type UpsertOrderDestination struct {
	AddressInfo          UpsertOrderAddressInfo `json:"addressInfo"`
	DeliveryInstructions string                 `json:"deliveryInstructions"`
	NodeInfo             UpsertOrderNodeInfo    `json:"nodeInfo"`
}

// UpsertOrderOrderType representa el tipo de pedido
type UpsertOrderOrderType struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

// Map convierte a domain.OrderType
func (o UpsertOrderOrderType) Map() domain.OrderType {
	return domain.OrderType{
		Type:        o.Type,
		Description: o.Description,
	}
}

// UpsertOrderOrigin representa el origen
type UpsertOrderOrigin struct {
	AddressInfo UpsertOrderAddressInfo `json:"addressInfo"`
	NodeInfo    UpsertOrderNodeInfo    `json:"nodeInfo"`
}

// UpsertOrderLabel representa una etiqueta
type UpsertOrderLabel struct {
	Type  string `json:"type" example:"packageCode"`
	Value string `json:"value" example:"uuid"`
}

// Map convierte a domain.Reference
func (l UpsertOrderLabel) Map() domain.Reference {
	return domain.Reference{
		Type:  l.Type,
		Value: l.Value,
	}
}

// UpsertOrderDimensions representa las dimensiones
type UpsertOrderDimensions struct {
	Length int64 `json:"length" example:"100" description:"Length in centimeters (cm)"`
	Height int64 `json:"height" example:"100" description:"Height in centimeters (cm)"`
	Width  int64 `json:"width" example:"100" description:"Width in centimeters (cm)"`
}

// Map convierte a domain.Dimensions
func (d UpsertOrderDimensions) Map() domain.Dimensions {
	return domain.Dimensions{
		Length: d.Length,
		Height: d.Height,
		Width:  d.Width,
		Unit:   "cm",
	}
}

// UpsertOrderItem representa un ítem
type UpsertOrderItem struct {
	Sku         string                `json:"sku" example:"SKU123"`
	Description string                `json:"description" example:"Cama 1 plaza"`
	Dimensions  UpsertOrderDimensions `json:"dimensions"`
	Weight      int64                 `json:"weight" example:"1000" description:"Weight in grams (g)"`
	Quantity    int                   `json:"quantity" example:"1" description:"Quantity in units"`
	Insurance   int64                 `json:"insurance" example:"10000" description:"Insurance value in currency units (CLP, MXN, PEN, etc.) - only integer values accepted"`
}

// Map convierte a domain.Item
func (i UpsertOrderItem) Map() domain.Item {
	return domain.Item{
		Sku:         i.Sku,
		Description: i.Description,
		Dimensions:  i.Dimensions.Map(),
		Weight:      i.Weight,
		Quantity:    i.Quantity,
		Insurance:   i.Insurance,
	}
}

// UpsertOrderDeliveryUnit representa una unidad de entrega
type UpsertOrderDeliveryUnit struct {
	Lpn          string             `json:"lpn" example:"LPN456"`
	SizeCategory string             `json:"sizeCategory" example:"SMALL"`
	Volume       int64              `json:"volume" example:"1000" description:"Volume in cubic centimeters (cm³)"`
	Weight       int64              `json:"weight" example:"1000" description:"Weight in grams (g)"`
	Insurance    int64              `json:"insurance" example:"10000" description:"Insurance value in currency units (CLP, MXN, PEN, CENTS etc.) - only integer values accepted"`
	Skills       []string           `json:"skills"`
	Labels       []UpsertOrderLabel `json:"labels"`
	Items        []UpsertOrderItem  `json:"items"`
}

// Map convierte a domain.DeliveryUnit
func (d UpsertOrderDeliveryUnit) Map() domain.DeliveryUnit {
	// Calculate volume from items if not provided
	volume := d.Volume
	if volume == 0 && len(d.Items) > 0 {
		for _, item := range d.Items {
			itemVolume := item.Dimensions.Length * item.Dimensions.Width * item.Dimensions.Height
			volume += itemVolume * int64(item.Quantity)
		}
	}

	// Calculate weight from items if not provided
	weight := d.Weight
	if weight == 0 && len(d.Items) > 0 {
		for _, item := range d.Items {
			weight += item.Weight * int64(item.Quantity)
		}
	}

	// Map skills
	skills := make([]domain.Skill, len(d.Skills))
	for i, skill := range d.Skills {
		skills[i] = domain.Skill(skill)
	}

	// Map labels
	labels := make([]domain.Reference, len(d.Labels))
	for i, label := range d.Labels {
		labels[i] = label.Map()
	}

	// Map items
	items := make([]domain.Item, len(d.Items))
	for i, item := range d.Items {
		items[i] = item.Map()
	}

	return domain.DeliveryUnit{
		Lpn:          d.Lpn,
		SizeCategory: domain.SizeCategory{Code: d.SizeCategory},
		Volume:       volume,
		Weight:       weight,
		Insurance:    d.Insurance,
		Status:       domain.Status{Status: domain.StatusAvailable},
		Skills:       skills,
		Labels:       labels,
		Items:        items,
	}
}

// UpsertOrderDateRange representa un rango de fechas
type UpsertOrderDateRange struct {
	EndDate   string `json:"endDate" example:"2025-03-30"`
	StartDate string `json:"startDate"  example:"2025-03-28"`
}

// Map convierte a domain.DateRange
func (d UpsertOrderDateRange) Map() domain.DateRange {
	return domain.DateRange{
		StartDate: d.parseDate(d.StartDate),
		EndDate:   d.parseDate(d.EndDate),
	}
}

// parseDate convierte string a time.Time
func (d UpsertOrderDateRange) parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

// UpsertOrderPromisedDate representa la fecha prometida
type UpsertOrderPromisedDate struct {
	DateRange       UpsertOrderDateRange `json:"dateRange"`
	ServiceCategory string               `json:"serviceCategory" example:"REGULAR / SAME DAY"`
	TimeRange       UpsertOrderTimeRange `json:"timeRange"`
}

// Map convierte a domain.PromisedDate
func (p UpsertOrderPromisedDate) Map() domain.PromisedDate {
	return domain.PromisedDate{
		DateRange:       p.DateRange.Map(),
		TimeRange:       p.TimeRange.Map(),
		ServiceCategory: p.ServiceCategory,
	}
}

// UpsertOrderReference representa una referencia
type UpsertOrderReference struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Map convierte a domain.Reference
func (r UpsertOrderReference) Map() domain.Reference {
	return domain.Reference{
		Type:  r.Type,
		Value: r.Value,
	}
}

type UpsertOrderRequest struct {
	ReferenceID             string                             `json:"referenceID" validate:"required" example:"1400001234567890"`
	ExtraFields             map[string]string                  `json:"extraFields"`
	GroupBy                 UpsertOrderGroupBy                 `json:"groupBy"`
	CollectAvailabilityDate UpsertOrderCollectAvailabilityDate `json:"collectAvailabilityDate"`
	Destination             UpsertOrderDestination             `json:"destination"`
	OrderType               UpsertOrderOrderType               `json:"orderType"`
	Origin                  UpsertOrderOrigin                  `json:"origin"`
	DeliveryUnits           []UpsertOrderDeliveryUnit          `json:"deliveryUnits"`
	PromisedDate            UpsertOrderPromisedDate            `json:"promisedDate"`
	References              []UpsertOrderReference             `json:"references"`
}

// Map convierte el request a un objeto de dominio Order
func (req UpsertOrderRequest) Map(ctx context.Context) domain.Order {
	// Map references
	references := make([]domain.Reference, len(req.References))
	for i, ref := range req.References {
		references[i] = ref.Map()
	}

	// Map delivery units
	deliveryUnits := make([]domain.DeliveryUnit, len(req.DeliveryUnits))
	for i, unit := range req.DeliveryUnits {
		deliveryUnits[i] = unit.Map()
	}

	// Map origin node info
	originNodeInfo := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(req.Origin.NodeInfo.ReferenceID),
		AddressInfo: req.Origin.AddressInfo.Map(),
	}

	// Map destination node info
	destinationNodeInfo := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(req.Destination.NodeInfo.ReferenceID),
		AddressInfo: req.Destination.AddressInfo.Map(),
	}

	order := domain.Order{
		ReferenceID:             domain.ReferenceID(req.ReferenceID),
		OrderType:               req.OrderType.Map(),
		References:              references,
		Origin:                  originNodeInfo,
		Destination:             destinationNodeInfo,
		DeliveryUnits:           deliveryUnits,
		CollectAvailabilityDate: req.CollectAvailabilityDate.Map(),
		PromisedDate:            req.PromisedDate.Map(),
		DeliveryInstructions:    req.Destination.DeliveryInstructions,
		ExtraFields:             req.ExtraFields,
		GroupBy: struct {
			Type  string
			Value string
		}{
			Type:  req.GroupBy.Type,
			Value: req.GroupBy.Value,
		},
	}
	order.Headers.SetFromContext(ctx)
	if order.Commerce == "" {
		order.Commerce = "empty"
	}
	if order.Consumer == "" {
		order.Consumer = "empty"
	}
	return order
}

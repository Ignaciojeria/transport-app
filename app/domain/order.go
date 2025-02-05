package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/joomcode/errorx"
)

type Order struct {
	Headers
	ID                      int64
	ReferenceID             ReferenceID             `json:"referenceID"`
	OrderStatus             OrderStatus             `json:"orderStatus"`
	OrderType               OrderType               `json:"orderType"`
	References              []Reference             `json:"references"`
	Origin                  NodeInfo                `json:"origin"`
	Destination             NodeInfo                `json:"destination"`
	Items                   []Item                  `json:"items"`
	Packages                []Package               `json:"packages"`
	CollectAvailabilityDate CollectAvailabilityDate `json:"collectAvailabilityDate"`
	PromisedDate            PromisedDate            `json:"promisedDate"`
	DeliveryInstructions    string                  `json:"deliveryInstructions"`
	TransportRequirements   []Reference             `json:"transportRequirements"`
}

func (o Order) UpdateIfChanged(newOrder Order) Order {
	// Actualizar usando los métodos UpdateIfChanged existentes
	o.Headers = o.Headers.UpdateIfChanged(newOrder.Headers)
	o.Origin.AddressInfo.Contact =
		o.Origin.Contact.UpdateIfChanged(newOrder.Origin.AddressInfo.Contact)
	o.Destination.AddressInfo.Contact =
		o.Origin.Contact.UpdateIfChanged(newOrder.Destination.AddressInfo.Contact)
	o.OrderStatus = o.OrderStatus.UpdateIfChanged(newOrder.OrderStatus)
	o.OrderType = o.OrderType.UpdateIfChanged(newOrder.OrderType)
	o.Origin = o.Origin.UpdateIfChanged(newOrder.Origin)
	o.Destination = o.Destination.UpdateIfChanged(newOrder.Destination)

	// Actualizar referencias - reemplazar directamente si hay nuevas
	if len(newOrder.References) > 0 {
		o.References = newOrder.References
	}

	// Actualizar packages - reemplazar directamente si hay nuevos
	if len(newOrder.Packages) > 0 {
		o.Packages = newOrder.Packages
	}

	// Update basic order information
	if newOrder.ReferenceID != "" {
		o.ReferenceID = newOrder.ReferenceID
	}

	// Update delivery instructions
	if newOrder.DeliveryInstructions != "" {
		o.DeliveryInstructions = newOrder.DeliveryInstructions
	}

	// Update CollectAvailabilityDate individual fields
	if !newOrder.CollectAvailabilityDate.Date.IsZero() {
		o.CollectAvailabilityDate.Date = newOrder.CollectAvailabilityDate.Date
	}
	if newOrder.CollectAvailabilityDate.TimeRange.StartTime != "" {
		o.CollectAvailabilityDate.TimeRange.StartTime = newOrder.CollectAvailabilityDate.TimeRange.StartTime
	}
	if newOrder.CollectAvailabilityDate.TimeRange.EndTime != "" {
		o.CollectAvailabilityDate.TimeRange.EndTime = newOrder.CollectAvailabilityDate.TimeRange.EndTime
	}

	// Update PromisedDate individual fields
	if !newOrder.PromisedDate.DateRange.StartDate.IsZero() {
		o.PromisedDate.DateRange.StartDate = newOrder.PromisedDate.DateRange.StartDate
	}
	if !newOrder.PromisedDate.DateRange.EndDate.IsZero() {
		o.PromisedDate.DateRange.EndDate = newOrder.PromisedDate.DateRange.EndDate
	}
	if newOrder.PromisedDate.TimeRange.StartTime != "" {
		o.PromisedDate.TimeRange.StartTime = newOrder.PromisedDate.TimeRange.StartTime
	}
	if newOrder.PromisedDate.TimeRange.EndTime != "" {
		o.PromisedDate.TimeRange.EndTime = newOrder.PromisedDate.TimeRange.EndTime
	}
	if newOrder.PromisedDate.ServiceCategory != "" {
		o.PromisedDate.ServiceCategory = newOrder.PromisedDate.ServiceCategory
	}

	// Update TransportRequirements if not empty
	if len(newOrder.TransportRequirements) > 0 {
		o.TransportRequirements = newOrder.TransportRequirements
	}

	// Update Items if not empty
	if len(newOrder.Items) > 0 {
		o.Items = newOrder.Items
	}

	// Update Organization if changed
	if newOrder.Organization.OrganizationCountryID != 0 {
		o.Organization = newOrder.Organization
	}

	return o
}

func (o Order) Validate() error {
	// Validar fechas de disponibilidad de recolección
	if err := o.ValidateCollectAvailabilityDate(); err != nil {
		return errorx.Decorate(err, "validation failed for CollectAvailabilityDate")
	}

	// Validar fechas prometidas
	if err := o.ValidatePromisedDate(); err != nil {
		return errorx.Decorate(err, "validation failed for PromisedDate")
	}

	if err := o.ValidatePackages(); err != nil {
		return errorx.Decorate(err, "validation failed for Packages")
	}

	// Validar otras reglas de dominio (si las hay)
	// Por ejemplo, puedes agregar reglas adicionales aquí

	return nil
}

func (o Order) ValidatePromisedDate() error {
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar que EndDate no sea anterior a StartDate cuando ambas están definidas
	if (o.PromisedDate.DateRange.StartDate != time.Time{}) && (o.PromisedDate.DateRange.EndDate != time.Time{}) {
		if o.PromisedDate.DateRange.EndDate.Before(o.PromisedDate.DateRange.StartDate) {
			return errorx.Decorate(
				ErrInvalidDateFormat.New("invalid endDate"),
				"promised delivery endDate cannot be before startDate",
			)
		}
	}

	// Validar formato de horas
	if o.PromisedDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.StartTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid startTime"),
			"promised delivery startTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.StartTime,
		)
	}
	if o.PromisedDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.EndTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid endTime"),
			"promised delivery endTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.EndTime,
		)
	}

	return nil
}

func (o Order) ValidateCollectAvailabilityDate() error {
	// Si la fecha está definida (no es zero value), solo validamos que sea una fecha válida
	// Ya no validamos si está en el pasado
	if (o.CollectAvailabilityDate.Date != time.Time{}) {
		// Aquí podrías agregar otras validaciones específicas si las necesitas
		// pero no la comparación con time.Now()
	}

	// Validar el rango horario
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)
	if o.CollectAvailabilityDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.CollectAvailabilityDate.TimeRange.StartTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid startTime"),
			"collect startTime: %s, expected format hh:mm",
			o.CollectAvailabilityDate.TimeRange.StartTime,
		)
	}
	if o.CollectAvailabilityDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.CollectAvailabilityDate.TimeRange.EndTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid endTime"),
			"collect endTime: %s, expected format hh:mm",
			o.CollectAvailabilityDate.TimeRange.EndTime,
		)
	}

	return nil
}

func (o Order) IsOriginAndDestinationContactEqual() bool {
	originContact := o.Origin.AddressInfo.Contact
	destinationContact := o.Destination.AddressInfo.Contact

	return originContact.FullName == destinationContact.FullName &&
		originContact.Email == destinationContact.Email &&
		originContact.Phone == destinationContact.Phone &&
		originContact.NationalID == destinationContact.NationalID
}

func (o Order) IsOriginAndDestinationAddressEqual() bool {
	originAddress := o.Origin.AddressInfo.RawAddress()
	destinationAddress := o.Destination.AddressInfo.RawAddress()

	return originAddress == destinationAddress
}

func (o Order) IsOriginAndDestinationNodeEqual() bool {
	return o.Origin.ReferenceID == o.Destination.ReferenceID
}

type ReferenceID string

type Reference struct {
	ID    int64
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NodeInfo struct {
	ID           int64
	ReferenceID  ReferenceID  `json:"referenceId"`
	Organization Organization `json:"organization"`
	Name         string       `json:"name"`
	NodeType     NodeType     `json:"type"`
	Contact      Contact      `json:"contact"`
	References   []Reference  `json:"references"`
	AddressInfo  AddressInfo  `json:"addressInfo"`
}

func (n NodeInfo) UpdateIfChanged(newNode NodeInfo) NodeInfo {
	// Actualizar ReferenceID
	if newNode.ID != 0 {
		n.ID = newNode.ID
	}
	if newNode.ReferenceID != "" && n.ReferenceID != newNode.ReferenceID {
		n.ReferenceID = newNode.ReferenceID
	}
	// Actualizar Name
	if newNode.Name != "" {
		n.Name = newNode.Name
	}
	// Actualizar Type
	if newNode.NodeType.Value != "" && n.NodeType.Value != newNode.NodeType.Value {
		n.NodeType.Value = newNode.NodeType.Value
	}
	// Actualizar NodeReferences
	if len(newNode.References) > 0 {
		n.References = newNode.References
	}
	if newNode.AddressInfo.ID != 0 {
		n.AddressInfo.ID = newNode.AddressInfo.ID
	}
	if newNode.Contact.ID != 0 {
		n.Contact.ID = newNode.Contact.ID
	}
	if newNode.Contact.ID != 0 {
		n.Contact.ID = newNode.Contact.ID
	}
	if newNode.NodeType.ID != 0 {
		n.NodeType.ID = newNode.NodeType.ID
	}
	n.Organization = newNode.Organization
	//n.AddressInfo = n.AddressInfo.UpdateIfChanged(newNode.AddressInfo)
	//n.NodeType = n.NodeType.UpdateIfChanged(newNode.NodeType)
	return n
}

type NodeType struct {
	ID           int64
	Organization Organization
	Value        string `json:"type"`
}

func (nt NodeType) UpdateIfChanged(newNodeType NodeType) NodeType {
	if newNodeType.ID != 0 {
		nt.ID = newNodeType.ID
	}
	if newNodeType.Value != "" {
		nt.Value = newNodeType.Value
	}
	if newNodeType.ID != 0 {
		nt.ID = newNodeType.ID
	}
	nt.Organization = newNodeType.Organization
	return nt
}

type Document struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Contact struct {
	ID           int64
	Organization Organization `json:"organization"`
	FullName     string       `json:"fullName"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	NationalID   string       `json:"nationalID"`
	Documents    []Document   `json:"documents"`
}

func (c Contact) UpdateIfChanged(newContact Contact) Contact {
	updatedContact := c // Copiamos la instancia actual

	// Actualizar FullName
	if newContact.ID != 0 {
		updatedContact.ID = newContact.ID
	}

	// Actualizar FullName
	if newContact.FullName != "" {
		updatedContact.FullName = newContact.FullName
	}

	// Actualizar Email
	if newContact.Email != "" {
		updatedContact.Email = newContact.Email
	}

	// Actualizar Phone
	if newContact.Phone != "" {
		updatedContact.Phone = newContact.Phone
	}

	// Actualizar NationalID
	if newContact.NationalID != "" {
		updatedContact.NationalID = newContact.NationalID
	}

	// Actualizar Documents
	if len(newContact.Documents) > 0 {
		updatedContact.Documents = newContact.Documents
	}

	return updatedContact
}

// Función auxiliar para comparar arreglos de documentos
func compareDocuments(oldDocs, newDocs []Document) bool {
	if len(oldDocs) != len(newDocs) {
		return false
	}
	for i := range oldDocs {
		if oldDocs[i] != newDocs[i] {
			return false
		}
	}
	return true
}

type AddressInfo struct {
	ID           int64
	Organization Organization
	Contact      Contact `json:"contact"`
	State        string  `json:"state"`
	County       string  `json:"county"`
	Province     string  `json:"province"`
	District     string  `json:"district"`
	AddressLine1 string  `json:"addressLine1"`
	AddressLine2 string  `json:"addressLine2"`
	AddressLine3 string  `json:"addressLine3"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	ZipCode      string  `json:"zipCode"`
	TimeZone     string  `json:"timeZone"`
}

func (a AddressInfo) UpdateIfChanged(newAddress AddressInfo) AddressInfo {
	if newAddress.ID != 0 {
		a.ID = newAddress.ID
	}
	if newAddress.AddressLine1 != "" {
		a.AddressLine1 = newAddress.AddressLine1
	}
	if newAddress.AddressLine2 != "" {
		a.AddressLine2 = newAddress.AddressLine2
	}
	if newAddress.AddressLine3 != "" {
		a.AddressLine3 = newAddress.AddressLine3
	}
	if newAddress.Latitude != 0 {
		a.Latitude = newAddress.Latitude
	}
	if newAddress.Longitude != 0 {
		a.Longitude = newAddress.Longitude
	}
	if newAddress.State != "" {
		a.State = newAddress.State
	}
	if newAddress.County != "" {
		a.County = newAddress.County
	}
	if newAddress.Province != "" {
		a.Province = newAddress.Province
	}
	if newAddress.District != "" {
		a.District = newAddress.District
	}
	if newAddress.ZipCode != "" {
		a.ZipCode = newAddress.ZipCode
	}
	if newAddress.TimeZone != "" {
		a.TimeZone = newAddress.TimeZone
	}
	return a
}

func (addr AddressInfo) RawAddress() string {
	return concatenateWithCommas(addr.AddressLine1, addr.AddressLine2, addr.AddressLine3)
}

func concatenateWithCommas(values ...string) string {
	result := ""
	for _, value := range values {
		if value != "" {
			if result != "" {
				result += ", "
			}
			result += value
		}
	}
	return result
}

type Quantity struct {
	QuantityNumber int    `json:"quantityNumber"`
	QuantityUnit   string `json:"quantityUnit"`
}

type Insurance struct {
	UnitValue float64 `json:"unitValue"`
	Currency  string  `json:"currency"`
}

type Dimensions struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Depth  float64 `json:"depth"`
	Unit   string  `json:"unit"`
}

type Weight struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Item struct {
	ReferenceID       ReferenceID `json:"referenceId"`
	LogisticCondition string      `json:"logisticCondition"`
	Quantity          Quantity    `json:"quantity"`
	Insurance         Insurance   `json:"insurance"`
	Description       string      `json:"description"`
	Dimensions        Dimensions  `json:"dimensions"`
	Weight            Weight      `json:"weight"`
}

type ItemReference struct {
	ReferenceID ReferenceID `json:"referenceId"`
	Quantity    Quantity    `json:"quantity"`
}

func (o *Order) ValidatePackages() error {
	// Crear un mapa para verificar rápidamente si un ReferenceID pertenece a los ítems de la orden
	itemMap := make(map[ReferenceID]bool)
	for _, item := range o.Items {
		itemMap[item.ReferenceID] = true
	}

	if len(o.Packages) == 1 {
		// Si solo hay un paquete y no tiene referencias de ítems, asignar todos los ítems de la orden
		if len(o.Packages[0].ItemReferences) == 0 {
			for _, item := range o.Items {
				o.Packages[0].ItemReferences = append(o.Packages[0].ItemReferences, ItemReference{
					ReferenceID: item.ReferenceID,
					Quantity:    item.Quantity,
				})
			}
		} else {
			// Validar que todas las referencias del paquete sean válidas
			for _, ref := range o.Packages[0].ItemReferences {
				if !itemMap[ref.ReferenceID] {
					return fmt.Errorf("validation failed: item reference ID '%s' in package is not part of the order items", ref.ReferenceID)
				}
			}
		}
	} else {
		// Si hay más de un paquete, validar que todos los paquetes tengan referencias de ítems
		for _, p := range o.Packages {
			if len(p.ItemReferences) == 0 {
				return errors.New("validation failed: packages with no item references must be explicitly defined when there are multiple packages")
			}

			// Validar que todas las referencias del paquete sean válidas
			for _, ref := range p.ItemReferences {
				if !itemMap[ref.ReferenceID] {
					return fmt.Errorf("validation failed: item reference ID '%s' in package is not part of the order items", ref.ReferenceID)
				}
			}
		}
	}

	return nil
}

type Package struct {
	ID             int64
	Lpn            string `json:"lpn"`
	Organization   Organization
	Dimensions     Dimensions      `json:"dimensions"`
	Weight         Weight          `json:"weight"`
	Insurance      Insurance       `json:"insurance"`
	ItemReferences []ItemReference `json:"itemReferences"`
}

func SearchPackageByLpn(pcks []Package, lpn string) Package {
	for _, pck := range pcks {
		if pck.Lpn == lpn {
			return pck
		}
	}
	return Package{}
}

func (p Package) UpdateIfChanged(newPackage Package) Package {
	// Actualizar Lpn
	if newPackage.Lpn != "" {
		p.Lpn = newPackage.Lpn
	}

	// Actualizar dimensiones si no están vacías
	if newPackage.Dimensions != (Dimensions{}) {
		p.Dimensions = newPackage.Dimensions
	}

	// Actualizar peso si no está vacío
	if newPackage.Weight != (Weight{}) {
		p.Weight = newPackage.Weight
	}

	// Actualizar seguro si no está vacío
	if newPackage.Insurance != (Insurance{}) {
		p.Insurance = newPackage.Insurance
	}

	// Actualizar referencias de ítems
	if len(newPackage.ItemReferences) > 0 {
		p.ItemReferences = newPackage.ItemReferences
	}
	return p
}

// Función auxiliar para comparar arreglos de referencias de ítems
func compareItemReferences(oldRefs, newRefs []ItemReference) bool {
	if len(oldRefs) != len(newRefs) {
		return false
	}
	for i := range oldRefs {
		if oldRefs[i] != newRefs[i] {
			return false
		}
	}
	return true
}

type OrderType struct {
	ID           int64
	Organization Organization
	Type         string `json:"type"`
	Description  string `json:"description"`
}

func (ot OrderType) UpdateIfChanged(newOrderType OrderType) OrderType {
	if newOrderType.ID != 0 {
		ot.ID = newOrderType.ID
	}
	if newOrderType.Type != "" {
		ot.Type = newOrderType.Type
	}
	if newOrderType.Description != "" {
		ot.Description = newOrderType.Description
	}
	if newOrderType.Organization.OrganizationCountryID != 0 {
		ot.Organization = newOrderType.Organization
	}
	return ot
}

type OrderStatus struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (os OrderStatus) UpdateIfChanged(newOrderStatus OrderStatus) OrderStatus {
	// Actualizar ID si es diferente de 0
	if newOrderStatus.ID != 0 {
		os.ID = newOrderStatus.ID
	}

	// Actualizar Status si no está vacío
	if newOrderStatus.Status != "" {
		os.Status = newOrderStatus.Status
	}

	// Actualizar CreatedAt si no es zero value
	if !newOrderStatus.CreatedAt.IsZero() {
		os.CreatedAt = newOrderStatus.CreatedAt
	}

	return os
}

type PromisedDate struct {
	DateRange       DateRange `json:"dateRange"`
	TimeRange       TimeRange `json:"timeRange"`
	ServiceCategory string    `json:"serviceCategory"`
}

type CollectAvailabilityDate struct {
	Date      time.Time `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type OrderCheckout struct {
	Order             Order
	Route             Route
	OrderStatus       OrderStatus
	DeliveredAt       time.Time
	Vehicle           Vehicle
	Recipient         Recipient
	EvicencePhotos    []EvicencePhotos
	Latitude          float32
	Longitude         float32
	NotDeliveryReason NotDeliveryReason
}

type Recipient struct {
	FullName   string
	NationalID string
}

type EvicencePhotos struct {
	ID      int64
	URL     string
	Type    string
	TakenAt time.Time
}

type NotDeliveryReason struct {
	ID          int64
	ReferenceID string
	Detail      string
}

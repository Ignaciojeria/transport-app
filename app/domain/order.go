package domain

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/joomcode/errorx"
)

type Order struct {
	ID                      int64
	ReferenceID             ReferenceID             `json:"referenceID"`
	Organization            Organization            `json:"organization"`
	BusinessIdentifiers     BusinessIdentifiers     `json:"businessIdentifiers"`
	OrderStatus             OrderStatus             `json:"orderStatus"`
	OrderType               OrderType               `json:"orderType"`
	References              []Reference             `json:"references"`
	Origin                  Origin                  `json:"origin"`
	Destination             Destination             `json:"destination"`
	Items                   []Item                  `json:"items"`
	Packages                []Package               `json:"packages"`
	CollectAvailabilityDate CollectAvailabilityDate `json:"collectAvailabilityDate"`
	PromisedDate            PromisedDate            `json:"promisedDate"`
	Visits                  []Visit                 `json:"visits"`
	TransportRequirements   []Reference             `json:"transportRequirements"`
	NeedsUpdate             bool                    `json:"-"`
}

func (o *Order) HydrateOrder(newOrder Order) {
	// Reiniciar el flag global de actualización
	needsUpdate := false

	// Comparar y actualizar campos simples de Order
	if newOrder.ReferenceID != "" && o.ReferenceID != newOrder.ReferenceID {
		o.ReferenceID = newOrder.ReferenceID
		needsUpdate = true
	}

	if newOrder.OrderType.Type != "" && o.OrderType != newOrder.OrderType {
		o.OrderType = newOrder.OrderType
		needsUpdate = true
	}

	// Sobrescribir completamente las referencias
	if len(newOrder.References) > 0 {
		o.References = newOrder.References
		needsUpdate = true
	}

	// Sobrescribir completamente las visitas
	if len(newOrder.Visits) > 0 {
		o.Visits = newOrder.Visits
		needsUpdate = true
	}

	// Sobrescribir completamente los ítems
	if len(newOrder.Items) > 0 {
		o.Items = newOrder.Items
		needsUpdate = true
	}

	// Comparar y actualizar los paquetes
	for i := range o.Packages {
		if i < len(newOrder.Packages) {
			o.Packages[i].UpdateIfChanged(newOrder.Packages[i])
			if o.Packages[i].NeedsUpdate {
				needsUpdate = true
			}
		}
	}

	// Agregar paquetes nuevos
	if len(newOrder.Packages) > len(o.Packages) {
		for _, newPkg := range newOrder.Packages[len(o.Packages):] {
			newPkg.NeedsUpdate = true // Los nuevos paquetes siempre requieren persistencia
			o.Packages = append(o.Packages, newPkg)
			needsUpdate = true
		}
	}

	// Comparar y actualizar origen
	if o.Origin.UpdateIfChanged(newOrder.Origin) {
		needsUpdate = true
	}

	// Comparar y actualizar destino
	if o.Destination.UpdateIfChanged(newOrder.Destination) {
		needsUpdate = true
	}

	// Sobrescribir las promesas de entrega
	if newOrder.PromisedDate.DateRange.StartDate != "" {
		o.PromisedDate = newOrder.PromisedDate
		needsUpdate = true
	}

	// Sobrescribir la fecha de disponibilidad para recolección
	if newOrder.CollectAvailabilityDate.Date != "" {
		o.CollectAvailabilityDate = newOrder.CollectAvailabilityDate
		needsUpdate = true
	}

	// Sobrescribir los identificadores de negocio
	if newOrder.BusinessIdentifiers.Commerce != "" {
		o.BusinessIdentifiers = newOrder.BusinessIdentifiers
		needsUpdate = true
	}

	// Sobrescribir los requerimientos de transporte
	if len(newOrder.TransportRequirements) > 0 {
		o.TransportRequirements = newOrder.TransportRequirements
		needsUpdate = true
	}

	// Actualizar el flag global
	o.NeedsUpdate = needsUpdate
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
	// Validar formato de fecha yyyy-mm-dd
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar el rango de fechas
	if o.PromisedDate.DateRange.StartDate != "" && !dateRegex.MatchString(o.PromisedDate.DateRange.StartDate) {
		return errorx.Decorate(
			ErrInvalidDateFormat.New("invalid startDate"),
			"startDate: %s, expected format yyyy-mm-dd",
			o.PromisedDate.DateRange.StartDate,
		)
	}
	if o.PromisedDate.DateRange.EndDate != "" && !dateRegex.MatchString(o.PromisedDate.DateRange.EndDate) {
		return errorx.Decorate(
			ErrInvalidDateFormat.New("invalid endDate"),
			"endDate: %s, expected format yyyy-mm-dd",
			o.PromisedDate.DateRange.EndDate,
		)
	}

	// Validar los rangos horarios
	if o.PromisedDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.StartTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid startTime"),
			"startTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.StartTime,
		)
	}
	if o.PromisedDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.EndTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid endTime"),
			"endTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.EndTime,
		)
	}

	return nil
}

func (o Order) ValidateCollectAvailabilityDate() error {
	// Validar formato de fecha yyyy-mm-dd
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar la fecha
	if o.CollectAvailabilityDate.Date != "" && !dateRegex.MatchString(o.CollectAvailabilityDate.Date) {
		return errorx.Decorate(
			ErrInvalidDateFormat.New("invalid date"),
			"collect date: %s, expected format yyyy-mm-dd",
			o.CollectAvailabilityDate.Date,
		)
	}

	// Validar el rango horario
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
	return o.Origin.NodeInfo.ReferenceID == o.Destination.NodeInfo.ReferenceID
}

type ReferenceID string

type Reference struct {
	ID    int64
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NodeInfo struct {
	ID          int64
	ReferenceID ReferenceID `json:"referenceId"`
	Name        *string     `json:"name"`
	Type        string      `json:"type"`
	Operator    Operator    `json:"operator"`
	References  []Reference `json:"references"`
}

func (n *NodeInfo) UpdateIfChanged(newNode NodeInfo) bool {
	needsUpdate := false

	if newNode.ReferenceID != "" && n.ReferenceID != newNode.ReferenceID {
		n.ReferenceID = newNode.ReferenceID
		needsUpdate = true
	}

	if newNode.Name != nil && (n.Name == nil || *n.Name != *newNode.Name) {
		n.Name = newNode.Name
		needsUpdate = true
	}

	if newNode.Type != "" && n.Type != newNode.Type {
		n.Type = newNode.Type
		needsUpdate = true
	}

	return needsUpdate
}

type Origin struct {
	ID          int64
	NodeInfo    NodeInfo    `json:"nodeInfo"`
	AddressInfo AddressInfo `json:"addressInfo"`
}

func (o *Origin) UpdateIfChanged(newOrigin Origin) bool {
	needsUpdate := false

	// Comparar y actualizar NodeInfo
	if o.NodeInfo.UpdateIfChanged(newOrigin.NodeInfo) {
		needsUpdate = true
	}

	// Comparar y actualizar AddressInfo
	if o.AddressInfo.UpdateIfChanged(newOrigin.AddressInfo) {
		needsUpdate = true
	}

	return needsUpdate
}

type Document struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Contact struct {
	ID          int64
	FullName    string     `json:"fullName"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	NationalID  string     `json:"nationalID"`
	Documents   []Document `json:"documents"`
	NeedsUpdate bool       `json:"-"`
}

func (c *Contact) UpdateIfChanged(newContact Contact) bool {
	needsUpdate := false

	// Comparar y actualizar FullName
	if newContact.FullName != "" && c.FullName != newContact.FullName {
		c.FullName = newContact.FullName
		needsUpdate = true
	}

	// Comparar y actualizar Email
	if newContact.Email != "" && c.Email != newContact.Email {
		c.Email = newContact.Email
		needsUpdate = true
	}

	// Comparar y actualizar Phone
	if newContact.Phone != "" && c.Phone != newContact.Phone {
		c.Phone = newContact.Phone
		needsUpdate = true
	}

	// Comparar y actualizar NationalID
	if newContact.NationalID != "" && c.NationalID != newContact.NationalID {
		c.NationalID = newContact.NationalID
		needsUpdate = true
	}

	// Comparar y actualizar Documents
	if len(newContact.Documents) > 0 {
		if len(c.Documents) != len(newContact.Documents) || !compareDocuments(c.Documents, newContact.Documents) {
			c.Documents = newContact.Documents
			needsUpdate = true
		}
	}

	return needsUpdate
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
	NeedsUpdate  bool    `json:"-"`
}

func (a *AddressInfo) UpdateIfChanged(newAddress AddressInfo) bool {
	needsUpdate := false

	// Comparar y actualizar AddressLine1
	if newAddress.AddressLine1 != "" && a.AddressLine1 != newAddress.AddressLine1 {
		a.AddressLine1 = newAddress.AddressLine1
		needsUpdate = true
	}

	// Comparar y actualizar AddressLine2
	if newAddress.AddressLine2 != "" && a.AddressLine2 != newAddress.AddressLine2 {
		a.AddressLine2 = newAddress.AddressLine2
		needsUpdate = true
	}

	// Comparar y actualizar AddressLine3
	if newAddress.AddressLine3 != "" && a.AddressLine3 != newAddress.AddressLine3 {
		a.AddressLine3 = newAddress.AddressLine3
		needsUpdate = true
	}

	// Comparar y actualizar Latitude
	if newAddress.Latitude != 0 && a.Latitude != newAddress.Latitude {
		a.Latitude = newAddress.Latitude
		needsUpdate = true
	}

	// Comparar y actualizar Longitude
	if newAddress.Longitude != 0 && a.Longitude != newAddress.Longitude {
		a.Longitude = newAddress.Longitude
		needsUpdate = true
	}

	// Comparar y actualizar State
	if newAddress.State != "" && a.State != newAddress.State {
		a.State = newAddress.State
		needsUpdate = true
	}

	// Comparar y actualizar County
	if newAddress.County != "" && a.County != newAddress.County {
		a.County = newAddress.County
		needsUpdate = true
	}

	// Comparar y actualizar Province
	if newAddress.Province != "" && a.Province != newAddress.Province {
		a.Province = newAddress.Province
		needsUpdate = true
	}

	// Comparar y actualizar District
	if newAddress.District != "" && a.District != newAddress.District {
		a.District = newAddress.District
		needsUpdate = true
	}

	// Comparar y actualizar ZipCode
	if newAddress.ZipCode != "" && a.ZipCode != newAddress.ZipCode {
		a.ZipCode = newAddress.ZipCode
		needsUpdate = true
	}

	// Comparar y actualizar TimeZone
	if newAddress.TimeZone != "" && a.TimeZone != newAddress.TimeZone {
		a.TimeZone = newAddress.TimeZone
		needsUpdate = true
	}

	// Comparar y actualizar Contact
	if a.Contact.UpdateIfChanged(newAddress.Contact) {
		needsUpdate = true
	}

	// Actualizar el flag global
	a.NeedsUpdate = needsUpdate

	return needsUpdate
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

type Operator struct {
	ID      int64
	Contact Contact `json:"contact"`
	Type    string  `json:"type"`
}

type Destination struct {
	ID                   int64
	DeliveryInstructions string      `json:"deliveryInstructions"`
	NodeInfo             NodeInfo    `json:"nodeInfo"`
	AddressInfo          AddressInfo `json:"addressInfo"`
}

func (d *Destination) UpdateIfChanged(newDestination Destination) bool {
	needsUpdate := false

	// Comparar y actualizar DeliveryInstructions
	if newDestination.DeliveryInstructions != "" && d.DeliveryInstructions != newDestination.DeliveryInstructions {
		d.DeliveryInstructions = newDestination.DeliveryInstructions
		needsUpdate = true
	}

	// Comparar y actualizar NodeInfo
	if d.NodeInfo.UpdateIfChanged(newDestination.NodeInfo) {
		needsUpdate = true
	}

	// Comparar y actualizar AddressInfo
	if d.AddressInfo.UpdateIfChanged(newDestination.AddressInfo) {
		needsUpdate = true
	}

	return needsUpdate
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
	Lpn            string          `json:"lpn"`
	Dimensions     Dimensions      `json:"dimensions"`
	Weight         Weight          `json:"weight"`
	Insurance      Insurance       `json:"insurance"`
	ItemReferences []ItemReference `json:"itemReferences"`
	NeedsUpdate    bool            `json:"-"`
}

func (p *Package) UpdateIfChanged(newPackage Package) (updatedPackage Package, needsUpdate bool) {
	// Crear una copia del paquete actual
	updatedPackage = *p
	needsUpdate = false

	// Comparar y actualizar dimensiones
	if newPackage.Dimensions != (Dimensions{}) {
		if p.Dimensions != newPackage.Dimensions {
			updatedPackage.Dimensions = newPackage.Dimensions
			needsUpdate = true
		}
	}

	// Comparar y actualizar peso
	if newPackage.Weight != (Weight{}) {
		if p.Weight != newPackage.Weight {
			updatedPackage.Weight = newPackage.Weight
			needsUpdate = true
		}
	}

	// Comparar y actualizar seguro
	if newPackage.Insurance != (Insurance{}) {
		if p.Insurance != newPackage.Insurance {
			updatedPackage.Insurance = newPackage.Insurance
			needsUpdate = true
		}
	}

	// Comparar y actualizar referencias de ítems
	if len(newPackage.ItemReferences) > 0 {
		if len(p.ItemReferences) != len(newPackage.ItemReferences) || !compareItemReferences(p.ItemReferences, newPackage.ItemReferences) {
			updatedPackage.ItemReferences = newPackage.ItemReferences
			needsUpdate = true
		}
	}

	return updatedPackage, needsUpdate
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
	ID          int64
	Type        string `json:"type"`
	Description string `json:"description"`
}

type OrderStatus struct {
	ID        int64
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type PromisedDate struct {
	DateRange       DateRange `json:"dateRange"`
	TimeRange       TimeRange `json:"timeRange"`
	ServiceCategory string    `json:"serviceCategory"`
}

type CollectAvailabilityDate struct {
	Date      string    `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type DateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type Visit struct {
	Date      string    `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type BusinessIdentifiers struct {
	CommerceID int64
	Commerce   string `json:"commerce"`
	ConsumerID int64
	Consumer   string `json:"consumer"`
}

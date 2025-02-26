package table

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"transport-app/app/domain"

	"github.com/biter777/countries"
	"github.com/paulmach/orb"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`

	ReferenceID string `gorm:"type:varchar(191);default:null;uniqueIndex:idx_reference_organization_country"`

	OrganizationCountryID int64               `gorm:"default:null;uniqueIndex:idx_reference_organization_country"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`

	OrderHeadersID int64        `gorm:"default:null"`
	OrderHeaders   OrderHeaders `gorm:"foreignKey:OrderHeadersID"`

	OrderStatusID int64       `gorm:"default:null"`
	OrderStatus   OrderStatus `gorm:"foreignKey:OrderStatusID"`

	OrderTypeID int64     `gorm:"default:null"`
	OrderType   OrderType `gorm:"foreignKey:OrderTypeID"`

	RouteID *int64 `gorm:"default:null"`
	Route   Route  `gorm:"foreignKey:RouteID"`

	PlanID *int64 `gorm:"default:null"`
	Plan   Plan   `gorm:"foreignKey:PlanID"`

	OrderReferences      []OrderReferences `gorm:"foreignKey:OrderID"`
	DeliveryInstructions string            `gorm:"type:text"`

	// Contacto asociado a la orden
	OriginContactID int64   `gorm:"default:null"`               // Clave foránea al Contact
	OriginContact   Contact `gorm:"foreignKey:OriginContactID"` // Relación con Contact

	// Contacto asociado a la orden
	DestinationContactID int64   `gorm:"default:null"`                    // Clave foránea al Contact
	DestinationContact   Contact `gorm:"foreignKey:DestinationContactID"` // Relación con Contact

	// Dirección de oriden de la orden de compra
	OriginAddressInfoID int64       `gorm:"default:null"`                   // Clave foránea al AddressInfo
	OriginAddressInfo   AddressInfo `gorm:"foreignKey:OriginAddressInfoID"` // Relación con AddressInfo

	// Dirección de destino de la orden de compra
	DestinationAddressInfoID int64       `gorm:"default:null"`                        // Clave foránea al AddressInfo
	DestinationAddressInfo   AddressInfo `gorm:"foreignKey:DestinationAddressInfoID"` // Relación con AddressInfo

	// Nodo de Origen de la orden (en caso de que tenga)
	OriginNodeInfoID int64    `gorm:"default:null"`
	OriginNodeInfo   NodeInfo `gorm:"foreignKey:OriginNodeInfoID"`

	// Nodo de Destino de la orden (en caso de que tenga)
	DestinationNodeInfoID int64    `gorm:"default:null"`
	DestinationNodeInfo   NodeInfo `gorm:"foreignKey:DestinationNodeInfoID"`

	SequenceNumber *int `gorm:"default:null"`

	JSONPlannedData JSONPlannedData `gorm:"type:json"`

	Packages []Package `gorm:"many2many:order_packages"`

	JSONItems JSONItems `gorm:"type:json"`

	CollectAvailabilityDate           time.Time     `gorm:"type:date;default:null"`
	CollectAvailabilityTimeRangeStart string        `gorm:"default:null"`
	CollectAvailabilityTimeRangeEnd   string        `gorm:"default:null"`
	PromisedDateRangeStart            time.Time     `gorm:"type:date;default:null"`
	PromisedDateRangeEnd              time.Time     `gorm:"type:date;default:null"`
	PromisedTimeRangeStart            string        `gorm:"default:null"`
	PromisedTimeRangeEnd              string        `gorm:"default:null"`
	TransportRequirements             JSONReference `gorm:"type:json"`
}

type PlannedData struct {
	JSONPlanLocation          PlanLocation
	JSONPlanCorrectedLocation PlanLocation
	PlanCorrectedDistance     float64
}

type JSONPlannedData PlannedData

// Implementación de sql.Scanner para deserializar JSON desde la BD
func (j *JSONPlannedData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONPlannedData value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Implementación de driver.Valuer para serializar a JSON en la BD
func (j JSONPlannedData) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type CheckoutRejection struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	ReferenceID string `gorm:"type:varchar(191);not null"`
	Reason      string `gorm:"type:varchar(191);not null"`
	Detail      string `gorm:"type:text"`
}

type CheckoutHistory struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`

	OrderID int64 `gorm:"not null;index"`
	Order   Order `gorm:"foreignKey:OrderID"`

	VehicleID int64   `gorm:"not null;index"`
	Vehicle   Vehicle `gorm:"foreignKey:VehicleID"`

	CarrierID int64   `gorm:"not null;index"`
	Carrier   Carrier `gorm:"foreignKey:CarrierID"`

	RouteID int64 `gorm:"not null;index"`
	Route   Route `gorm:"foreignKey:RouteID"`

	OrderStatusID int64       `gorm:"not null;index"`
	OrderStatus   OrderStatus `gorm:"foreignKey:OrderStatusID"`

	CheckoutRejectionID int64             `gorm:"not null;index"`
	CheckoutRejection   CheckoutRejection `gorm:"foreignKey:CheckoutRejectionID"`

	Latitude  float64
	Longitude float64

	RecipientFullName   string
	RecipientNationalID string

	EvidencePhotos JSONEvidencePhotos `gorm:"type:json"`
}

type EvidencePhoto struct {
	URL     string
	Type    string
	TakenAt time.Time
}

// Definimos el tipo para manejar el array de EvidencePhoto como JSON
type JSONEvidencePhotos []EvidencePhoto

// Implementamos los métodos necesarios para el manejo de JSON
func (j *JSONEvidencePhotos) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONEvidencePhotos value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON
func (j JSONEvidencePhotos) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (o Order) Map() domain.Order {
	// Mapear la orden base
	var planID int64
	if o.PlanID != nil {
		planID = *o.PlanID
	}
	var routeID int64
	if o.RouteID != nil {
		routeID = *o.RouteID
	}

	order := domain.Order{
		ID:          o.ID,
		ReferenceID: domain.ReferenceID(o.ReferenceID),
		Headers: domain.Headers{
			Organization: domain.Organization{
				OrganizationCountryID: o.OrganizationCountryID,
			},
		},
		Plan: domain.Plan{
			ID: planID,
			Routes: []domain.Route{
				{
					ID: routeID,
				},
			},
		},
		DeliveryInstructions: o.DeliveryInstructions,
	}

	// Mapear las fechas de disponibilidad de recolección
	order.CollectAvailabilityDate = domain.CollectAvailabilityDate{
		Date: o.CollectAvailabilityDate,
		TimeRange: domain.TimeRange{
			StartTime: o.CollectAvailabilityTimeRangeStart,
			EndTime:   o.CollectAvailabilityTimeRangeEnd,
		},
	}

	// Mapear las fechas prometidas
	order.PromisedDate = domain.PromisedDate{
		DateRange: domain.DateRange{
			StartDate: o.PromisedDateRangeStart,
			EndDate:   o.PromisedDateRangeEnd,
		},
		TimeRange: domain.TimeRange{
			StartTime: o.PromisedTimeRangeStart,
			EndTime:   o.PromisedTimeRangeEnd,
		},
	}

	// Mapear requisitos de transporte
	order.TransportRequirements = o.TransportRequirements.Map()

	// Mapear items
	items := make([]domain.Item, len(o.JSONItems))
	for i, item := range o.JSONItems {
		items[i] = item.Map()
	}
	order.Items = items

	// Mapear Contact IDs
	if o.OriginContactID != 0 {
		order.Origin.Contact.ID = o.OriginContactID
	}
	if o.DestinationContactID != 0 {
		order.Destination.Contact.ID = o.DestinationContactID
	}

	// Mapear AddressInfo IDs
	if o.OriginAddressInfoID != 0 {
		order.Origin.AddressInfo.ID = o.OriginAddressInfoID
	}
	if o.DestinationAddressInfoID != 0 {
		order.Destination.AddressInfo.ID = o.DestinationAddressInfoID
	}

	// Mapear NodeInfo IDs
	if o.OriginNodeInfoID != 0 {
		order.Origin.ID = o.OriginNodeInfoID
	}
	if o.DestinationNodeInfoID != 0 {
		order.Destination.ID = o.DestinationNodeInfoID
	}

	return order
}

type Items struct {
	ReferenceID       string         `gorm:"not null" json:"reference_id"`
	LogisticCondition string         `gorm:"default:null" json:"logistic_condition"`
	QuantityNumber    int            `gorm:"not null" json:"quantity_number"`
	QuantityUnit      string         `gorm:"not null" json:"quantity_unit"`
	JSONInsurance     JSONInsurance  `gorm:"type:json" json:"insurance"`
	Description       string         `gorm:"type:text" json:"description"`
	JSONDimensions    JSONDimensions `gorm:"type:json" json:"dimensions"`
	JSONWeight        JSONWeight     `gorm:"type:json" json:"weight"`
}

type JSONItems []Items

func (i Items) Map() domain.Item {
	return domain.Item{
		ReferenceID:       domain.ReferenceID(i.ReferenceID),
		LogisticCondition: i.LogisticCondition,
		Quantity: domain.Quantity{
			QuantityNumber: i.QuantityNumber,
			QuantityUnit:   i.QuantityUnit,
		},
		Insurance:   i.JSONInsurance.Map(),
		Description: i.Description,
		Dimensions:  i.JSONDimensions.Map(),
		Weight:      i.JSONWeight.Map(),
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONItems) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONItems value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir datos a JSON al guardarlos en la base de datos
func (j JSONItems) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Contact struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_contact_unique"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	FullName              string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_contact_unique"`
	Email                 string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_contact_unique"`
	Phone                 string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_contact_unique"`
	NationalID            string              `gorm:"type:varchar(191);default:null;uniqueIndex:idx_contact_unique"`
	Documents             JSONReference       `gorm:"type:json"`
}

func (c Contact) Map() domain.Contact {
	return domain.Contact{
		ID: c.ID,
		Organization: domain.Organization{
			OrganizationCountryID: c.OrganizationCountryID,
		},
		FullName:   c.FullName,
		Email:      c.Email,
		Phone:      c.Phone,
		NationalID: c.NationalID,
		Documents:  c.Documents.MapDocuments(),
	}
}

type Reference struct {
	Type  string `gorm:"not null" json:"type"`
	Value string `gorm:"not null" json:"value"`
}

type ItemReference struct {
	ReferenceID    string `json:"reference_id"`
	QuantityNumber int    `json:"quantity_number"`
	QuantityUnit   string `json:"quantity_unit"`
}

type JSONItemReferences []ItemReference

func (j JSONItemReferences) Map() []domain.ItemReference {
	mappedReferences := make([]domain.ItemReference, len(j))
	for i, ref := range j {
		mappedReferences[i] = domain.ItemReference{
			ReferenceID: domain.ReferenceID(ref.ReferenceID),
			Quantity: domain.Quantity{
				QuantityNumber: ref.QuantityNumber,
				QuantityUnit:   ref.QuantityUnit,
			},
		}
	}
	return mappedReferences
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONItemReferences) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONItemReferences) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type JSONReference []Reference

func (j JSONReference) Map() []domain.Reference {
	mappedReferences := make([]domain.Reference, len(j))
	for i, ref := range j {
		mappedReferences[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mappedReferences
}

func (j JSONReference) MapDocuments() []domain.Document {
	mappedReferences := make([]domain.Document, len(j))
	for i, ref := range j {
		mappedReferences[i] = domain.Document{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mappedReferences
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONReference) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONReference) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Package struct {
	gorm.Model
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_lpn_org"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ID                    int64               `gorm:"primaryKey"`
	Lpn                   string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_lpn_org"`
	JSONDimensions        JSONDimensions      `gorm:"type:json"`
	JSONWeight            JSONWeight          `gorm:"type:json"`
	JSONInsurance         JSONInsurance       `gorm:"type:json"`
	JSONItemsReferences   JSONItemReferences  `gorm:"type:json"`
}

func (p Package) Map() domain.Package {
	return domain.Package{
		ID:             p.ID,
		Organization:   p.OrganizationCountry.Map(),
		Lpn:            p.Lpn,
		Dimensions:     p.JSONDimensions.Map(),
		Weight:         p.JSONWeight.Map(),
		Insurance:      p.JSONInsurance.Map(),
		ItemReferences: p.JSONItemsReferences.Map(),
	}
}

type OrderPackage struct {
	gorm.Model
	OrderID   int64 `gorm:"not null;index"`
	PackageID int64 `gorm:"not null;index"`
}

type OrderReferences struct {
	gorm.Model
	ID      int64  `gorm:"primaryKey"`
	Type    string `gorm:"not null"`
	Value   string `gorm:"not null"`
	Order   Order  `gorm:"foreignKey:OrderID"`
	OrderID int64  `gorm:"index"`
}

type NodeInfo struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	ReferenceID           string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_reference_organization"`                 // Parte del índice único con OrganizationCountryID
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_reference_organization;uniqueIndex:idx_name_organization"` // Parte de ambos índices únicos
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`                                                  // Relación con OrganizationCountry
	Name                  string              `gorm:"type:varchar(191);default:null;uniqueIndex:idx_name_organization"`                  // Parte del índice único con OrganizationCountryID
	NodeTypeID            *int64              `gorm:"default:null"`
	NodeType              NodeType            `gorm:"foreignKey:NodeTypeID"`
	ContactID             *int64              `gorm:"default:null"`
	Contact               Contact             `gorm:"foreignKey:ContactID"` // Relación con Operator
	AddressID             *int64              `gorm:"default:null"`         // Clave foránea a AddressInfo
	AddressInfo           AddressInfo         `gorm:"foreignKey:AddressID"` // Relación con AddressInfo
	NodeReferences        JSONReference       `gorm:"type:json"`            // Relación con NodeReferences
}

func (n NodeInfo) Map() domain.NodeInfo {
	var contactID, addressID, nodeTypeID int64
	if n.ContactID != nil {
		contactID = *n.ContactID
	}
	if n.AddressID != nil {
		addressID = *n.AddressID
	}
	if n.NodeTypeID != nil {
		nodeTypeID = *n.NodeTypeID
	}
	nodeInfo := domain.NodeInfo{
		ID:          n.ID,
		ReferenceID: domain.ReferenceID(n.ReferenceID),
		Name:        n.Name,
		Organization: domain.Organization{
			OrganizationCountryID: n.OrganizationCountryID,
		},
		NodeType: domain.NodeType{
			ID: nodeTypeID,
		},
		References: n.NodeReferences.Map(),
		Contact: domain.Contact{
			ID: contactID,
		},
		AddressInfo: domain.AddressInfo{
			ID: addressID,
		},
	}
	return nodeInfo
}

type AddressInfo struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_raw_address_organization"`                   // Parte del índice único
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`                                    // Relación con la tabla OrganizationCountry
	RawAddress            string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_raw_address_organization"` // Parte del índice único
	State                 string              `gorm:"default:null"`
	Province              string              `gorm:"default:null"`
	County                string              `gorm:"default:null"`
	District              string              `gorm:"default:null"`
	AddressLine1          string              `gorm:"not null"`
	AddressLine2          string              `gorm:"default:null"`
	AddressLine3          string              `gorm:"default:null"`
	Latitude              float64             `gorm:"default:null"`
	Longitude             float64             `gorm:"default:null"`
	ZipCode               string              `gorm:"default:null"`
	TimeZone              string              `gorm:"default:null"`
}

func (a AddressInfo) Map() domain.AddressInfo {
	return domain.AddressInfo{
		ID:           a.ID,
		Organization: a.OrganizationCountry.Map(),
		State:        a.State,
		Province:     a.Province,
		County:       a.County,
		District:     a.District,
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		AddressLine3: a.AddressLine3,
		Location:     orb.Point{a.Longitude, a.Latitude},
		ZipCode:      a.ZipCode,
		TimeZone:     a.TimeZone,
	}
}

type Insurance struct {
	UnitValue float64 `gorm:"not null" json:"unit_value"`
	Currency  string  `gorm:"not null" json:"currency"`
}

type JSONInsurance Insurance

func (j JSONInsurance) Map() domain.Insurance {
	return domain.Insurance{
		UnitValue: j.UnitValue,
		Currency:  j.Currency,
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONInsurance) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONInsurance) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Dimensions struct {
	Height float64 `gorm:"not null" json:"height"`
	Width  float64 `gorm:"not null" json:"width"`
	Depth  float64 `gorm:"not null" json:"depth"`
	Unit   string  `gorm:"not null" json:"unit"`
}

type JSONDimensions Dimensions

func (j JSONDimensions) Map() domain.Dimensions {
	return domain.Dimensions{
		Height: j.Height,
		Width:  j.Width,
		Depth:  j.Depth,
		Unit:   j.Unit,
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONDimensions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONDimensions) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Weight struct {
	WeightValue float64 `gorm:"not null" json:"weight_value"`
	WeightUnit  string  `gorm:"not null" json:"weight_unit"`
}

type JSONWeight Weight

func (j JSONWeight) Map() domain.Weight {
	return domain.Weight{
		Value: j.WeightValue,
		Unit:  j.WeightUnit,
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONWeight) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONWeight) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type OrderType struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	Type                  string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_type_organization"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_type_organization"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Description           string              `gorm:"type:text"`
}

func (o OrderType) Map() domain.OrderType {
	return domain.OrderType{
		ID:          o.ID,
		Type:        o.Type,
		Description: o.Description,
		Organization: domain.Organization{
			OrganizationCountryID: o.OrganizationCountryID,
		},
	}
}

type OrderStatus struct {
	gorm.Model
	ID     int64  `gorm:"primaryKey"`
	Status string `gorm:"not null"`
}

// Organization tables

type ApiKey struct {
	gorm.Model
	ID             int64        `gorm:"primaryKey"`
	OrganizationID int64        `gorm:"not null;index"`            // ID de la organización asociada
	Organization   Organization `gorm:"foreignKey:OrganizationID"` // Relación con la tabla Organization
	Key            string       `gorm:"not null;unique"`           // Clave única
	Status         string       `gorm:"default:active"`            // Estado: activo, revocado, etc.
}

type Organization struct {
	gorm.Model
	ID        int64                 `gorm:"primaryKey"`
	Email     string                `gorm:"type:varchar(255);not null;unique"`
	Name      string                `gorm:"type:varchar(255);not null;"`
	Countries []OrganizationCountry `gorm:"foreignKey:OrganizationID"` // Relación con países
}

func (o Organization) Map() domain.Organization {
	return domain.Organization{
		ID:    o.ID,
		Name:  o.Name,
		Email: o.Email,
	}
}

type OrganizationCountry struct {
	gorm.Model
	ID             int64        `gorm:"primaryKey"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	OrganizationID int64        `gorm:"not null;uniqueIndex:idx_organization_country"`              // Parte del índice compuesto
	Country        string       `gorm:"type:char(5);not null;uniqueIndex:idx_organization_country"` // Código ISO de 2 caracteres
}

func (o OrganizationCountry) Map() domain.Organization {
	return domain.Organization{
		ID:                    o.OrganizationID,
		OrganizationCountryID: o.ID,
		Country:               countries.ByName(o.Country),
		Name:                  o.Organization.Name,
		Email:                 o.Organization.Email,
	}
}

//Account Tables

type Account struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`

	Type        string `gorm:"default:null"`
	ReferenceID string `gorm:"type:varchar(50);uniqueIndex:idx_account_ref_org"`

	ContactID *int64  `gorm:"default:null;uniqueIndex:idx_organization_contact"`
	Contact   Contact `gorm:"foreignKey:ContactID"`

	AddressInfoID *int64      `gorm:"default:null"`
	AddressInfo   AddressInfo `gorm:"foreignKey:AddressInfoID"`

	IsActive bool `gorm:"default:null;index"`

	OriginNodeInfoID *int64   `gorm:"default:null"`
	OriginNodeInfo   NodeInfo `gorm:"foreignKey:OriginNodeInfoID"`

	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_organization_contact;uniqueIndex:idx_account_ref_org"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
}

func (a Account) MapOperator() domain.Operator {
	// Inicializamos un Operator base
	operator := domain.Operator{
		ID:          a.ID,
		ReferenceID: a.ReferenceID, // Usando el Type como ReferenceID
		Type:        a.Type,
		Organization: domain.Organization{
			OrganizationCountryID: a.OrganizationCountryID,
		},
	}

	// Mapeamos el Contact si existe
	if a.ContactID != nil {
		operator.Contact = a.Contact.Map()
	}

	// Mapeamos el OriginNode si existe
	if a.OriginNodeInfoID != nil {
		operator.OriginNode = a.OriginNodeInfo.Map()
	}

	// Mapeamos el OriginNode si existe
	if a.OriginNodeInfo.AddressID != nil {
		operator.OriginNode.AddressInfo = a.OriginNodeInfo.AddressInfo.Map()
	}

	// Mapeamos el OriginNode si existe
	if a.OriginNodeInfo.NodeTypeID != nil {
		operator.OriginNode.NodeType = a.OriginNodeInfo.NodeType.Map()
	}

	return operator
}

type JSONB json.RawMessage

// Scan implementa la interfaz `sql.Scanner` para deserializar valores JSON
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB([]byte("null"))
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for JSONB")
	}
	*j = JSONB(bytes)
	return nil
}

// Value implementa la interfaz `driver.Valuer` para serializar valores JSON
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Modelo de Empresa de Transporte
type Carrier struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"uniqueIndex:idx_carrier_ref_org;uniqueIndex:idx_carrier_national_org"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ReferenceID           string              `gorm:"type:varchar(50);uniqueIndex:idx_carrier_ref_org"`
	Name                  string              `gorm:"not null"`
	NationalID            string              `gorm:"type:varchar(20);default:null;uniqueIndex:idx_carrier_national_org"`
}

func (c Carrier) Map() domain.Carrier {
	return domain.Carrier{
		ID:          c.ID,
		ReferenceID: c.ReferenceID,
		Name:        c.Name,
		NationalID:  c.NationalID,
		Organization: domain.Organization{
			OrganizationCountryID: c.OrganizationCountryID,
		},
	}
}

// Modelo de Vehículo
type Vehicle struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;index;uniqueIndex:idx_vehicle_ref_org;uniqueIndex:idx_vehicle_plate_org"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ReferenceID           string              `gorm:"type:varchar(50);uniqueIndex:idx_vehicle_ref_org"`
	Plate                 string              `gorm:"type:varchar(20);uniqueIndex:idx_vehicle_plate_org"`
	IsActive              bool
	CertificateDate       string
	VehicleCategoryID     *int64          `gorm:"default null;index"`
	VehicleCategory       VehicleCategory `gorm:"foreignKey:VehicleCategoryID"`
	VehicleHeadersID      int64           `gorm:"not null"`
	VehicleHeaders        VehicleHeaders  `gorm:"foreignKey:VehicleHeadersID"`
	Weight                JSONB           `gorm:"type:json"`          // Tipo JSON para serializar Weight
	Insurance             JSONB           `gorm:"type:json"`          // Tipo JSON para serializar Insurance
	TechnicalReview       JSONB           `gorm:"type:json"`          // Tipo JSON para serializar TechnicalReview
	Dimensions            JSONB           `gorm:"type:json"`          // Tipo JSON para serializar Dimensions
	CarrierID             *int64          `gorm:"default null;index"` // Relación con Carrier
	Carrier               Carrier         `gorm:"foreignKey:CarrierID"`
}

func (v Vehicle) Map() domain.Vehicle {
	return domain.Vehicle{
		ID:              v.ID,
		ReferenceID:     v.ReferenceID,
		Plate:           v.Plate,
		IsActive:        v.IsActive,
		CertificateDate: v.CertificateDate,
		VehicleCategory: domain.VehicleCategory{
			ID:                  v.VehicleCategory.ID,
			Type:                v.VehicleCategory.Type,
			MaxPackagesQuantity: v.VehicleCategory.MaxPackagesQuantity,
		},
		Headers: v.VehicleHeaders.Map(),
		Carrier: v.Carrier.Map(),
	}
}

type VehicleCategory struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_org_country_type"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Type                  string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_org_country_type"`
	MaxPackagesQuantity   int
}

func (vc VehicleCategory) Map() domain.VehicleCategory {
	return domain.VehicleCategory{
		ID:                  vc.ID,
		MaxPackagesQuantity: vc.MaxPackagesQuantity,
		Type:                vc.Type,

		Organization: domain.Organization{OrganizationCountryID: vc.OrganizationCountryID},
	}
}

type OrderHeaders struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	Commerce              string              `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	Consumer              string              `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	OrganizationCountryID int64               `gorm:"not null;index;uniqueIndex:idx_commerce_consumer_org_country"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
}

func (m OrderHeaders) Map() domain.Headers {
	return domain.Headers{
		ID: m.ID,
		Organization: domain.Organization{
			OrganizationCountryID: m.OrganizationCountryID,
		},
		Consumer: m.Consumer,
		Commerce: m.Commerce,
	}
}

type VehicleHeaders struct {
	ID                    int64               `gorm:"primaryKey"`
	Commerce              string              `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	Consumer              string              `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	OrganizationCountryID int64               `gorm:"not null;index;uniqueIndex:idx_commerce_consumer_org_country"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
}

func (m VehicleHeaders) Map() domain.Headers {
	return domain.Headers{
		ID: m.ID,
		Organization: domain.Organization{
			OrganizationCountryID: m.OrganizationCountryID,
		},
		Consumer: m.Consumer,
		Commerce: m.Commerce,
	}
}

type NodeHeaders struct {
	ID                    int64               `gorm:"primaryKey"`
	Commerce              string              `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	Consumer              string              `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	OrganizationCountryID int64               `gorm:"not null;index;uniqueIndex:idx_commerce_consumer_org_country"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
}

func (m NodeHeaders) Map() domain.Headers {
	return domain.Headers{
		ID: m.ID,
		Organization: domain.Organization{
			OrganizationCountryID: m.OrganizationCountryID,
		},
		Consumer: m.Consumer,
		Commerce: m.Commerce,
	}
}

type NodeType struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_name_organization"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Name                  string              `gorm:"type:varchar(191);uniqueIndex:idx_name_organization"`
}

func (n NodeType) Map() domain.NodeType {
	return domain.NodeType{
		ID: n.ID,
		Organization: domain.Organization{
			OrganizationCountryID: n.OrganizationCountryID,
		},
		Value: n.Name,
	}
}

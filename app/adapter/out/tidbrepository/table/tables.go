package table

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"transport-app/app/domain"

	"github.com/biter777/countries"
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

	Packages []Package `gorm:"many2many:order_packages"`

	JSONItems JSONItems `gorm:"type:json"`

	CollectAvailabilityDate           string        `gorm:"default:null"`
	CollectAvailabilityTimeRangeStart string        `gorm:"default:null"`
	CollectAvailabilityTimeRangeEnd   string        `gorm:"default:null"`
	PromisedDateRangeStart            string        `gorm:"default:null"`
	PromisedDateRangeEnd              string        `gorm:"default:null"`
	PromisedTimeRangeStart            string        `gorm:"default:null"`
	PromisedTimeRangeEnd              string        `gorm:"default:null"`
	TransportRequirements             JSONReference `gorm:"type:json"`
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
	NodeTypeID            int64               `gorm:"default:null"`
	NodeType              NodeType            `gorm:"foreignKey:NodeTypeID"`
	ContactID             *int64              `gorm:"default:null"`
	Contact               Contact             `gorm:"foreignKey:ContactID"` // Relación con Operator
	AddressID             *int64              `gorm:"default:null"`         // Clave foránea a AddressInfo
	AddressInfo           AddressInfo         `gorm:"foreignKey:AddressID"` // Relación con AddressInfo
	NodeReferences        JSONReference       `gorm:"type:json"`            // Relación con NodeReferences
}

func (n NodeInfo) Map() domain.NodeInfo {
	var contactID, addressID int64
	if n.ContactID != nil {
		contactID = *n.ContactID
	}
	if n.AddressID != nil {
		addressID = *n.AddressID
	}
	return domain.NodeInfo{
		ID:          n.ID,
		ReferenceID: domain.ReferenceID(n.ReferenceID),
		Organization: domain.Organization{
			OrganizationCountryID: n.OrganizationCountryID,
		},
		Name: n.Name,
		NodeType: domain.NodeType{
			ID:    n.NodeType.ID,
			Value: n.NodeType.Name,
		},
		References: n.NodeReferences.Map(),
		Contact: domain.Contact{
			ID: contactID,
		},
		AddressInfo: domain.AddressInfo{
			ID: addressID,
		},
	}
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
	Latitude              float32             `gorm:"default:null"`
	Longitude             float32             `gorm:"default:null"`
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
		Latitude:     a.Latitude,
		Longitude:    a.Longitude,
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

	ContactID int64   `gorm:"not null;uniqueIndex:idx_organization_contact"`
	Contact   Contact `gorm:"foreignKey:ContactID"`

	IsActive bool `gorm:"not null;index"`

	OriginNodeInfoID int64    `gorm:"default:null"`
	OriginNodeInfo   NodeInfo `gorm:"foreignKey:OriginNodeInfoID"`

	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_organization_contact"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
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
	OrganizationCountryID int64               `gorm:"not null;index"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ReferenceID           string              `gorm:"unique;default:null"`
	Name                  string              `gorm:"not null"`
	NationalID            string              `gorm:"unique;not null"`
	Document              JSONB               `gorm:"type:json" json:"document"` // Tipo JSON para manejar estructuras anidadas
	Vehicles              []Vehicle           `gorm:"foreignKey:CarrierID"`
}

// Modelo de Vehículo
type Vehicle struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;index"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ReferenceID           string              `gorm:"unique;default:null"`
	Plate                 string              `gorm:"not null"`
	IsActive              bool
	CertificateDate       string
	Category              string
	Weight                JSONB   `gorm:"type:json"`      // Tipo JSON para serializar Weight
	Insurance             JSONB   `gorm:"type:json"`      // Tipo JSON para serializar Insurance
	TechnicalReview       JSONB   `gorm:"type:json"`      // Tipo JSON para serializar TechnicalReview
	Dimensions            JSONB   `gorm:"type:json"`      // Tipo JSON para serializar Dimensions
	CarrierID             int64   `gorm:"not null;index"` // Relación con Carrier
	Carrier               Carrier `gorm:"foreignKey:CarrierID"`
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
		Organization: domain.Organization{
			OrganizationCountryID: m.OrganizationCountryID,
		},
		Consumer: m.Consumer,
		Commerce: m.Commerce,
	}
}

type NodeType struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex,length:200"`
}

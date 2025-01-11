package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`

	ReferenceID string `gorm:"type:varchar(191);not null;uniqueIndex:idx_reference_organization_country"`

	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_reference_organization_country"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`

	CommerceID int64    `gorm:"not null"`
	Commerce   Commerce `gorm:"foreignKey:CommerceID"`

	ConsumerID int64    `gorm:"not null"`
	Consumer   Consumer `gorm:"foreignKey:ConsumerID"`

	OrderStatusID int64       `gorm:"not null"`
	OrderStatus   OrderStatus `gorm:"foreignKey:OrderStatusID"`

	OrderTypeID int64     `gorm:"not null"`
	OrderType   OrderType `gorm:"foreignKey:OrderTypeID"`

	OrderReferences      []OrderReferences `gorm:"foreignKey:OrderID"`
	DeliveryInstructions string            `gorm:"type:text"`

	// Contacto asociado a la orden
	OriginContactID int64   `gorm:"not null"`                   // Clave foránea al Contact
	OriginContact   Contact `gorm:"foreignKey:OriginContactID"` // Relación con Contact

	// Contacto asociado a la orden
	DestinationContactID int64   `gorm:"not null"`                        // Clave foránea al Contact
	DestinationContact   Contact `gorm:"foreignKey:DestinationContactID"` // Relación con Contact

	// Dirección de oriden de la orden de compra
	OriginAddressInfoID int64       `gorm:"not null"`                       // Clave foránea al AddressInfo
	OriginAddressInfo   AddressInfo `gorm:"foreignKey:OriginAddressInfoID"` // Relación con AddressInfo

	// Dirección de destino de la orden de compra
	DestinationAddressInfoID int64       `gorm:"not null"`                            // Clave foránea al AddressInfo
	DestinationAddressInfo   AddressInfo `gorm:"foreignKey:DestinationAddressInfoID"` // Relación con AddressInfo

	// Nodo de Origen de la orden (en caso de que tenga)
	OriginNodeInfoID int64    `gorm:"default:null"`
	OriginNodeInfo   NodeInfo `gorm:"foreignKey:OriginNodeInfoID"`

	// Nodo de Destino de la orden (en caso de que tenga)
	DestinationNodeInfoID int64    `gorm:"default:null"`
	DestinationNodeInfo   NodeInfo `gorm:"foreignKey:DestinationNodeInfoID"`

	JSONItems JSONItems `gorm:"type:json"`

	CollectAvailabilityDate           string  `gorm:"default:null"`
	CollectAvailabilityTimeRangeStart string  `gorm:"default:null"`
	CollectAvailabilityTimeRangeEnd   string  `gorm:"default:null"`
	PromisedDateRangeStart            string  `gorm:"default:null"`
	PromisedDateRangeEnd              string  `gorm:"default:null"`
	PromisedTimeRangeStart            string  `gorm:"default:null"`
	PromisedTimeRangeEnd              string  `gorm:"default:null"`
	Visits                            []Visit `gorm:"foreignKey:OrderID"`
	TransportRequirements             []byte  `gorm:"type:json"`
}

type Items struct {
	ReferenceID       string     `gorm:"not null"`
	LogisticCondition string     `gorm:"default:null"`
	Quantity          Quantity   `gorm:"embedded"`
	Insurance         Insurance  `gorm:"embedded"`
	Description       string     `gorm:"type:text"`
	Dimensions        Dimensions `gorm:"embedded"`
	Weight            Weight     `gorm:"embedded"`
}

type JSONItems []Items

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
	NationalID            string              `gorm:"type:varchar(191);default:null"`
	Documents             JSONDocuments       `gorm:"type:json"`
}

type Document struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type JSONDocuments []Document

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONDocuments) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONDocuments) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Package struct {
	gorm.Model
	ID         int64      `gorm:"primaryKey"`
	Lpn        string     `gorm:"type:varchar(191);not null;uniqueIndex"`
	Dimensions Dimensions `gorm:"embedded"`
	Weight     Weight     `gorm:"embedded"`
	Insurance  Insurance  `gorm:"embedded"`
	JSONItems  JSONItems  `gorm:"type:json"`
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
	Name                  *string             `gorm:"type:varchar(191);default:null;uniqueIndex:idx_name_organization"`                  // Parte del índice único con OrganizationCountryID
	Type                  string              `gorm:"not null"`
	OperatorID            int64               `gorm:"default:null"`
	Operator              Operator            `gorm:"foreignKey:OperatorID"` // Relación con Operator
	AddressID             int64               `gorm:"not null"`              // Clave foránea a AddressInfo
	AddressInfo           AddressInfo         `gorm:"foreignKey:AddressID"`  // Relación con AddressInfo
	NodeReferences        []NodeReference     `gorm:"foreignKey:NodeInfoID"` // Relación con NodeReferences
}

type AddressInfo struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_raw_address_organization"`                   // Parte del índice único
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`                                    // Relación con la tabla OrganizationCountry
	RawAddress            string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_raw_address_organization"` // Parte del índice único
	State                 string              `gorm:"default:null"`
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

type NodeReference struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;index"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Type                  string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_type_value_org_node"`
	Value                 string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_type_value_org_node"`
	NodeInfoID            int64               `gorm:"not null;uniqueIndex:idx_type_value_org_node"`
}

type Operator struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	ContactID             int64               `gorm:"not null;uniqueIndex:idx_contact_organization"` // Parte del índice único
	Contact               Contact             `gorm:"foreignKey:ContactID"`                          // Relación con la tabla Contact
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_contact_organization"` // Parte del índice único
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Type                  string              `gorm:"type:varchar(191);not null"`
}

type Quantity struct {
	QuantityNumber int    `gorm:"not null"`
	QuantityUnit   string `gorm:"not null"`
}

type Insurance struct {
	UnitValue int    `gorm:"not null"`
	Currency  string `gorm:"not null"`
}

type Dimensions struct {
	Height float64 `gorm:"not null"`
	Width  float64 `gorm:"not null"`
	Depth  float64 `gorm:"not null"`
	Unit   string  `gorm:"not null"`
}

type Weight struct {
	Value float64 `gorm:"not null"`
	Unit  string  `gorm:"not null"`
}

type OrderType struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	Type                  string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_type_organization"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_type_organization"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Description           string              `gorm:"type:text"`
}

type OrderStatus struct {
	gorm.Model
	ID     int64  `gorm:"primaryKey"`
	Status string `gorm:"not null"`
}

type Visit struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	Order          Order  `gorm:"foreignKey:OrderID"`
	OrderID        int64  `gorm:"not null;index:idx_transport_order_date"` // Part of composite unique index
	Date           string `gorm:"not null;index:idx_transport_order_date"` // Part of composite unique index
	TimeRangeStart string `gorm:"default:null"`
	TimeRangeEnd   string `gorm:"default:null"`
}

type Consumer struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	Name                  string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_organization"` // Índice único compuesto
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_name_organization"`                   // Actualizado para el índice único compuesto
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
}

type Commerce struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	Name                  string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_organization"` // Índice único compuesto
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_name_organization"`                   // Actualizado para el índice único compuesto
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
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

type OrganizationCountry struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	OrganizationID int64  `gorm:"not null;uniqueIndex:idx_organization_country"`              // Parte del índice compuesto
	Country        string `gorm:"type:char(5);not null;uniqueIndex:idx_organization_country"` // Código ISO de 2 caracteres
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

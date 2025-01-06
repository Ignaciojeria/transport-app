package table

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ID                       int64                      `gorm:"primaryKey"`
	ReferenceID              string                     `gorm:"type:varchar(191);not null;uniqueIndex:idx_reference_organization"`
	OrganizationCountryID    int64                      `gorm:"not null;uniqueIndex:idx_reference_org_country"`
	OrganizationCountry      OrganizationCountry        `gorm:"foreignKey:OrganizationCountryID"`
	CommerceID               int64                      `gorm:"not null"`
	Commerce                 Commerce                   `gorm:"foreignKey:CommerceID"`
	ConsumerID               int64                      `gorm:"not null"`
	Consumer                 Consumer                   `gorm:"foreignKey:ConsumerID"`
	OrderStatusID            int64                      `gorm:"not null"`
	OrderStatus              OrderStatus                `gorm:"foreignKey:OrderStatusID"`
	OrderTypeID              int64                      `gorm:"not null"`
	OrderType                OrderType                  `gorm:"foreignKey:OrderTypeID"`
	TransportOrderReferences []TransportOrderReferences `gorm:"foreignKey:TransportOrderID"`
	DeliveryInstructions     string                     `gorm:"type:text"`

	// Contacto asociado a la orden
	ContactID int64   `gorm:"not null"`             // Clave foránea al Contact
	Contact   Contact `gorm:"foreignKey:ContactID"` // Relación con Contact

	// Dirección de la orden de compra
	AddressID   int64       `gorm:"not null"`             // Clave foránea al AddressInfo
	AddressInfo AddressInfo `gorm:"foreignKey:AddressID"` // Relación con AddressInfo

	// Nodo de Origen de la orden (en caso de que tenga)
	OriginNodeInfoID int64    `gorm:"default:null"`
	OriginNodeInfo   NodeInfo `gorm:"foreignKey:OriginNodeInfoID"`

	// Nodo de Destino de la orden (en caso de que tenga)
	DestinationNodeInfoID int64    `gorm:"default:null"`
	DestinationNodeInfo   NodeInfo `gorm:"foreignKey:DestinationNodeInfoID"`

	Items                             []Items                           `gorm:"foreignKey:TransportOrderID"`
	Packages                          []Packages                        `gorm:"foreignKey:TransportOrderID"`
	CollectAvailabilityDate           string                            `gorm:"default:null"`
	CollectAvailabilityTimeRangeStart string                            `gorm:"default:null"`
	CollectAvailabilityTimeRangeEnd   string                            `gorm:"default:null"`
	PromisedDateRangeStart            string                            `gorm:"default:null"`
	PromisedDateRangeEnd              string                            `gorm:"default:null"`
	PromisedTimeRangeStart            string                            `gorm:"default:null"`
	PromisedTimeRangeEnd              string                            `gorm:"default:null"`
	Visit                             Visit                             `gorm:"foreignKey:TransportOrderID"`
	TransportRequirementsReferences   []TransportRequirementsReferences `gorm:"foreignKey:TransportOrderID"`
}

type Contact struct {
	ID         int64  `gorm:"primaryKey"`
	FullName   string `gorm:"not null"`
	Email      string `gorm:"not null"`
	Phone      string `gorm:"not null"`
	NationalID string `gorm:"default:null"`
	Documents  []byte `gorm:"type:json"`
}

type Packages struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	TransportOrderID      int64               `gorm:"not null;uniqueIndex:idx_transport_order_lpn"`
	OrganizationCountryID int64               `gorm:"not null;index"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Lpn                   string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_transport_order_lpn"`
	PackageType           string              `gorm:"type:varchar(191);default:null"`
	Dimensions            Dimensions          `gorm:"embedded"`
	Weight                Weight              `gorm:"embedded"`
	Insurance             Insurance           `gorm:"embedded"`
	ItemReferences        []ItemReferences    `gorm:"foreignKey:PackageID"`
}

type TransportOrderReferences struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	OrganizationCountryID int64               `gorm:"not null;index"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Type                  string              `gorm:"not null"`
	Value                 string              `gorm:"not null"`
	TransportOrderID      int64               `gorm:"index"`
}

type TransportRequirementsReferences struct {
	gorm.Model
	ID               int64  `gorm:"primaryKey"`
	Type             string `gorm:"not null"`
	Value            string `gorm:"not null"`
	TransportOrderID int64  `gorm:"index"`
}

type NodeInfo struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey"`
	ReferenceID           string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_reference_organization"`
	OrganizationCountryID int64               `gorm:"not null;uniqueIndex:idx_reference_org_country"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Name                  string              `gorm:"type:varchar(191);not null;uniqueIndex:idx_nodeinfo_organization_name"` // Único por organización
	Type                  string              `gorm:"not null"`
	OperatorID            int64               `gorm:"not null"`
	Operator              Operator            `gorm:"foreignKey:OperatorID"`
	AddressID             int64               `gorm:"not null"`             // Clave foránea al AddressInfo
	AddressInfo           AddressInfo         `gorm:"foreignKey:AddressID"` // Relación con AddressInfo
	NodeReferences        []NodeReferences    `gorm:"foreignKey:NodeInfoID"`
}

type AddressInfo struct {
	ID           int64   `gorm:"primaryKey"`
	State        string  `gorm:"default:null"`
	County       string  `gorm:"default:null"`
	District     string  `gorm:"default:null"`
	AddressLine1 string  `gorm:"not null"`
	AddressLine2 string  `gorm:"default:null"`
	AddressLine3 string  `gorm:"default:null"`
	RawAddress   string  `gorm:"type:varchar(191);not null;uniqueIndex"`
	Latitude     float64 `gorm:"default:null"`
	Longitude    float64 `gorm:"default:null"`
	ZipCode      string  `gorm:"default:null"`
	TimeZone     string  `gorm:"default:null"`
}

type NodeReferences struct {
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

type Items struct {
	gorm.Model
	ID                int64      `gorm:"primaryKey"`
	ReferenceID       string     `gorm:"not null"`
	LogisticCondition string     `gorm:"default:null"`
	Quantity          Quantity   `gorm:"embedded"`
	Insurance         Insurance  `gorm:"embedded"`
	Description       string     `gorm:"type:text"`
	Dimensions        Dimensions `gorm:"embedded"`
	Weight            Weight     `gorm:"embedded"`
	TransportOrderID  int64      `gorm:"not null"`
}

type ItemReferences struct {
	gorm.Model
	ID          int64    `gorm:"primaryKey"`
	ReferenceID string   `gorm:"not null"`
	Quantity    Quantity `gorm:"embedded"`
	PackageID   int64    `gorm:"not null"` // Foreign key reference
}

type OrderType struct {
	gorm.Model
	ID             int64        `gorm:"primaryKey"`
	Type           string       `gorm:"type:varchar(191);not null;uniqueIndex:idx_type_organization"` // Cambiar a varchar(191)
	OrganizationID int64        `gorm:"not null;uniqueIndex:idx_type_organization"`                   // Parte del índice compuesto
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	Description    string       `gorm:"type:text"` // Puede permanecer como TEXT porque no está en un índice
}

type OrderStatus struct {
	gorm.Model
	ID     int64  `gorm:"primaryKey"`
	Status string `gorm:"not null"`
}

type Visit struct {
	gorm.Model
	ID               int64  `gorm:"primaryKey"`
	TransportOrderID int64  `gorm:"not null;index:idx_transport_order_date"` // Part of composite unique index
	Date             string `gorm:"not null;index:idx_transport_order_date"` // Part of composite unique index
	TimeRangeStart   string `gorm:"default:null"`
	TimeRangeEnd     string `gorm:"default:null"`
}

type Consumer struct {
	gorm.Model
	ID             int64        `gorm:"primaryKey"`
	Name           string       `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_organization"` // Índice único compuesto
	OrganizationID int64        `gorm:"not null;uniqueIndex:idx_name_organization"`                   // Mismo índice único compuesto
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

type Commerce struct {
	gorm.Model
	ID             int64        `gorm:"primaryKey"`
	Name           string       `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_organization"` // Índice único compuesto
	OrganizationID int64        `gorm:"not null;uniqueIndex:idx_name_organization"`                   // Mismo índice único compuesto
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

package table

type TransportOrder struct {
	ID                                int64                             `gorm:"primaryKey"`
	ReferenceID                       string                            `gorm:"type:varchar(191);not null;uniqueIndex:idx_reference_organization"`
	OrganizationID                    int64                             `gorm:"not null;uniqueIndex:idx_reference_organization"`
	Organization                      Organization                      `gorm:"foreignKey:OrganizationID"`
	CommerceID                        int64                             `gorm:"not null"`
	Commerce                          Commerce                          `gorm:"foreignKey:CommerceID"`
	ConsumerID                        int64                             `gorm:"not null"`
	Consumer                          Consumer                          `gorm:"foreignKey:ConsumerID"`
	OrderStatusID                     int                               `gorm:"not null"`
	OrderStatus                       OrderStatus                       `gorm:"foreignKey:OrderStatusID"`
	OrderTypeID                       int64                             `gorm:"not null"`
	OrderType                         OrderType                         `gorm:"foreignKey:OrderTypeID"`
	TransportOrderReferences          []TransportOrderReferences        `gorm:"foreignKey:TransportOrderID"`
	DeliveryInstructions              string                            `gorm:"type:text"`
	OriginID                          int64                             `gorm:"not null"`
	Origin                            Origin                            `gorm:"foreignKey:OriginID"`
	DestinationID                     int64                             `gorm:"not null"`
	Destination                       Destination                       `gorm:"foreignKey:DestinationID"`
	DestinationAddressInfo            AddressInfo                       `gorm:"embedded"`
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

type Packages struct {
	ID               int64            `gorm:"primaryKey"`
	TransportOrderID int64            `gorm:"not null;uniqueIndex:idx_transport_order_lpn"`
	OrganizationID   int64            `gorm:"not null;index"`
	Organization     Organization     `gorm:"foreignKey:OrganizationID"`
	Lpn              string           `gorm:"type:varchar(191);not null;uniqueIndex:idx_transport_order_lpn"`
	PackageType      string           `gorm:"type:varchar(191);default:null"`
	Dimensions       Dimensions       `gorm:"embedded"`
	Weight           Weight           `gorm:"embedded"`
	Insurance        Insurance        `gorm:"embedded"`
	ItemReferences   []ItemReferences `gorm:"foreignKey:PackageID"`
}

type TransportOrderReferences struct {
	ID               int64        `gorm:"primaryKey"`
	OrganizationID   int64        `gorm:"not null;index"`
	Organization     Organization `gorm:"foreignKey:OrganizationID"`
	Type             string       `gorm:"not null"`
	Value            string       `gorm:"not null"`
	TransportOrderID int64        `gorm:"index"`
}

type TransportRequirementsReferences struct {
	ID               int64  `gorm:"primaryKey"`
	Type             string `gorm:"not null"`
	Value            string `gorm:"not null"`
	TransportOrderID int64  `gorm:"index"`
}

type NodeInfo struct {
	ID             int64            `gorm:"primaryKey"`
	OrganizationID int64            `gorm:"not null;index:idx_nodeinfo_organization"` // Índice para organización
	Organization   Organization     `gorm:"foreignKey:OrganizationID"`
	ReferenceID    string           `gorm:"type:varchar(191);not null;uniqueIndex:idx_nodeinfo_organization_referenceid"` // Único por organización
	Name           string           `gorm:"type:varchar(191);not null;uniqueIndex:idx_nodeinfo_organization_name"`        // Único por organización
	Type           string           `gorm:"not null"`
	OperatorID     int64            `gorm:"not null"`
	Operator       Operator         `gorm:"foreignKey:OperatorID"`
	NodeReferences []NodeReferences `gorm:"foreignKey:NodeInfoID"`
}

type NodeReferences struct {
	ID             int64        `gorm:"primaryKey"`
	OrganizationID int64        `gorm:"not null;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	Type           string       `gorm:"not null"`
	Value          string       `gorm:"not null"`
	NodeInfoID     int64        `gorm:"not null;index"` // Clave foránea a NodeInfo
}

type Operator struct {
	ID             int64        `gorm:"primaryKey"`
	OrganizationID int64        `gorm:"not null;index"` // Solo si pertenece a una organización
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	ReferenceID    string       `gorm:"unique;not null"`
	NationalID     string       `gorm:"unique;not null"`
	Type           string       `gorm:"not null"`
	Name           string       `gorm:"not null"`
}

type Contact struct {
	FullName       string `gorm:"not null"`
	ContactMethods []byte `gorm:"type:json"`
	Documents      []byte `gorm:"type:json"`
}

type Origin struct {
	ID            int64       `gorm:"primaryKey"`
	NodeInfoID    int64       `gorm:"not null"`
	NodeInfo      NodeInfo    `gorm:"foreignKey:NodeInfoID"`
	AddressInfoID int64       `gorm:"not null"` // Clave foránea para AddressInfo
	AddressInfo   AddressInfo `gorm:"foreignKey:AddressInfoID"`
}

type Destination struct {
	ID            int64       `gorm:"primaryKey"`
	NodeInfoID    int64       `gorm:"not null"`
	NodeInfo      NodeInfo    `gorm:"foreignKey:NodeInfoID"`
	AddressInfoID int64       `gorm:"not null"` // Clave foránea para AddressInfo
	AddressInfo   AddressInfo `gorm:"foreignKey:AddressInfoID"`
}

type AddressInfo struct {
	ID          int64   `gorm:"primaryKey"`
	Contact     Contact `gorm:"embedded"`
	State       string  `gorm:"not null"`
	County      string  `gorm:"not null"`
	District    string  `gorm:"not null"`
	FullAddress string  `gorm:"not null"`
	Latitude    float64 `gorm:"not null"`
	Longitude   float64 `gorm:"not null"`
	ZipCode     string  `gorm:"not null"`
	TimeZone    string  `gorm:"not null"`
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
	ID          int64    `gorm:"primaryKey"`
	ReferenceID string   `gorm:"not null"`
	Quantity    Quantity `gorm:"embedded"`
	PackageID   int64    `gorm:"not null"` // Foreign key reference
}

type OrderType struct {
	ID             int64        `gorm:"primaryKey"`
	OrganizationID int64        `gorm:"not null;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	Type           string       `gorm:"not null"`
	Description    string       `gorm:"type:text"`
}

type OrderStatus struct {
	ID     int    `gorm:"primaryKey"`
	Status string `gorm:"not null"`
}

type Visit struct {
	ID               int64  `gorm:"primaryKey"`
	TransportOrderID int64  `gorm:"not null;index:idx_transport_order_date"` // Part of composite unique index
	Date             string `gorm:"not null;index:idx_transport_order_date"` // Part of composite unique index
	TimeRangeStart   string `gorm:"default:null"`
	TimeRangeEnd     string `gorm:"default:null"`
}

type Consumer struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);not null;uniqueIndex"`
}

type Commerce struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);not null;uniqueIndex"`
}

type Organization struct {
	ID   int64
	Name string `gorm:"type:varchar(255);not null;uniqueIndex"`
}

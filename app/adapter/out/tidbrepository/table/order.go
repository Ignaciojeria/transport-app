package table

import (
	"time"
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`

	DocumentID string `gorm:"type:char(32);uniqueIndex"`

	ReferenceID string `gorm:"not null"`

	OrganizationID int64        `gorm:"not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`

	OrderHeadersDoc string       `gorm:"type:char(32);index"`
	OrderHeaders    OrderHeaders `gorm:"-"`

	OrderStatusDoc string      `gorm:"type:char(32);index"`
	OrderStatus    OrderStatus `gorm:"-"`

	OrderTypeDoc string    `gorm:"type:char(32);index"`
	OrderType    OrderType `gorm:"-"`

	RouteDoc string `gorm:"type:char(32);index"`
	Route    Route  `gorm:"-"`

	OrderReferences []OrderReferences `gorm:"-"`

	DeliveryInstructions string `gorm:"type:text"`

	// Contacto asociado a la orden
	OriginContactDoc string  `gorm:"type:char(32);index"`
	OriginContact    Contact `gorm:"-"`

	// Contacto asociado a la orden
	DestinationContactDoc string  `gorm:"type:char(32);index"`
	DestinationContact    Contact `gorm:"-"`

	// Dirección de oriden de la orden de compra
	OriginAddressInfoDoc string      `gorm:"type:char(32);index"`
	OriginAddressInfo    AddressInfo `gorm:"-"`

	// Dirección de destino de la orden de compra
	DestinationAddressInfoDoc string      `gorm:"type:char(32);index"`
	DestinationAddressInfo    AddressInfo `gorm:"-"`

	// Nodo de Origen de la orden (en caso de que tenga)
	OriginNodeInfoDoc string   `gorm:"type:char(32);index"`
	OriginNodeInfo    NodeInfo `gorm:"-"`

	// Nodo de Destino de la orden (en caso de que tenga)
	DestinationNodeInfoDoc string   `gorm:"type:char(32);index"`
	DestinationNodeInfo    NodeInfo `gorm:"-"`

	SequenceNumber *int `gorm:"default:null"`

	JSONPlannedData JSONPlannedData `gorm:"type:json"`

	Packages []Package `gorm:"-"`

	AddressLine2                      string        `gorm:"default:null"`
	CollectAvailabilityDate           *time.Time    `gorm:"type:date;default:null"`
	CollectAvailabilityTimeRangeStart string        `gorm:"default:null"`
	CollectAvailabilityTimeRangeEnd   string        `gorm:"default:null"`
	PromisedDateRangeStart            *time.Time    `gorm:"type:date;default:null"`
	PromisedDateRangeEnd              *time.Time    `gorm:"type:date;default:null"`
	PromisedTimeRangeStart            string        `gorm:"default:null"`
	PromisedTimeRangeEnd              string        `gorm:"default:null"`
	ServiceCategory                   string        `gorm:"default:null"`
	TransportRequirements             JSONReference `gorm:"type:json"`
}

func (o Order) Map() domain.Order {
	// Mapear la orden base

	order := domain.Order{
		//	ID:          o.ID,
		AddressLine2: o.AddressLine2,
		ReferenceID:  domain.ReferenceID(o.ReferenceID),
		/*
			Headers: domain.Headers{
				Organization: domain.Organization{
					ID:      o.OrganizationID,
					Country: countries.ByName(o.Organization.Country),
				},
			},
		*/
		//	Plan:                 domain.Plan{},
		//DeliveryInstructions: o.DeliveryInstructions,
	}

	// Mapear las fechas de disponibilidad de recolección
	order.CollectAvailabilityDate = domain.CollectAvailabilityDate{
		Date: safeTime(o.CollectAvailabilityDate),
		TimeRange: domain.TimeRange{
			StartTime: o.CollectAvailabilityTimeRangeStart,
			EndTime:   o.CollectAvailabilityTimeRangeEnd,
		},
	}

	// Mapear las fechas prometidas
	order.PromisedDate = domain.PromisedDate{
		DateRange: domain.DateRange{
			StartDate: safeTime(o.PromisedDateRangeStart),
			EndDate:   safeTime(o.PromisedDateRangeEnd),
		},
		TimeRange: domain.TimeRange{
			StartTime: o.PromisedTimeRangeStart,
			EndTime:   o.PromisedTimeRangeEnd,
		},
	}

	// Mapear requisitos de transporte
	order.TransportRequirements = o.TransportRequirements.Map()

	// Mapear items

	return order
}

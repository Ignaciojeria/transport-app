package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`

	DocumentID string `gorm:"type:char(64);uniqueIndex"`

	ReferenceID string `gorm:"not null"`

	TenantID uuid.UUID `gorm:"not null"`
	Tenant   Tenant    `gorm:"foreignKey:TenantID"`

	OrderHeadersDoc string       `gorm:"type:char(64);index"`
	OrderHeaders    OrderHeaders `gorm:"-"`

	OrderTypeDoc string    `gorm:"type:char(64);index"`
	OrderType    OrderType `gorm:"-"`

	RouteDoc string `gorm:"type:char(64);index"`
	Route    Route  `gorm:"-"`

	OrderReferences []OrderReferences `gorm:"-"`

	DeliveryInstructions string `gorm:"type:text"`

	// Contacto asociado a la orden
	OriginContactDoc string  `gorm:"type:char(64);index"`
	OriginContact    Contact `gorm:"-"`

	// Contacto asociado a la orden
	DestinationContactDoc string  `gorm:"type:char(64);index"`
	DestinationContact    Contact `gorm:"-"`

	// Dirección de oriden de la orden de compra
	OriginAddressInfoDoc string      `gorm:"type:char(64);index"`
	OriginAddressInfo    AddressInfo `gorm:"-"`

	// Dirección de destino de la orden de compra
	DestinationAddressInfoDoc string      `gorm:"type:char(64);index"`
	DestinationAddressInfo    AddressInfo `gorm:"-"`

	// Nodo de Origen de la orden (en caso de que tenga)
	OriginNodeInfoDoc string   `gorm:"type:char(64);index"`
	OriginNodeInfo    NodeInfo `gorm:"-"`

	// Nodo de Destino de la orden (en caso de que tenga)
	DestinationNodeInfoDoc string   `gorm:"type:char(64);index"`
	DestinationNodeInfo    NodeInfo `gorm:"-"`

	SequenceNumber *int `gorm:"default:null"`

	ExtraFields JSONMap `gorm:"type:json"`

	DeliveryUnits []DeliveryUnit `gorm:"-"`

	CollectAvailabilityDate           *time.Time `gorm:"type:date;default:null"`
	CollectAvailabilityTimeRangeStart string     `gorm:"default:null"`
	CollectAvailabilityTimeRangeEnd   string     `gorm:"default:null"`
	PromisedDateRangeStart            *time.Time `gorm:"type:date;default:null"`
	PromisedDateRangeEnd              *time.Time `gorm:"type:date;default:null"`
	PromisedTimeRangeStart            string     `gorm:"default:null"`
	PromisedTimeRangeEnd              string     `gorm:"default:null"`
	ServiceCategory                   string     `gorm:"default:null"`
}

type JSONMap map[string]string

func (m JSONMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = JSONMap{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for JSONMap, got %T", value)
	}
	return json.Unmarshal(bytes, m)
}

func (o Order) Map() domain.Order {
	order := domain.Order{
		ReferenceID: domain.ReferenceID(o.ReferenceID),
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

	return order
}

func safeTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{} // Retorna un time.Time vacío en lugar de nil
	}
	return *t
}

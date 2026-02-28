package supabaserepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"micartapro/app/events"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// ErrTrackingIDConflict es devuelto cuando tracking_id ya existe (colisión).
var ErrTrackingIDConflict = errors.New("tracking_id already exists")

// InsertOrderTracking inserta una fila en la proyección order_tracking.
// aggregate_id = ID del agregado de la orden. fulfillment puede ser nil (PICKUP sin datos extra).
// Retorna ErrTrackingIDConflict si tracking_id ya existe.
type InsertOrderTracking func(ctx context.Context, aggregateID int64, trackingID string, fulfillment *events.OrderFulfillment) error

func init() {
	ioc.Register(NewInsertOrderTracking)
}

func NewInsertOrderTracking(supabase *supabase.Client) InsertOrderTracking {
	return func(ctx context.Context, aggregateID int64, trackingID string, fulfillment *events.OrderFulfillment) error {
		record := map[string]interface{}{
			"aggregate_id": aggregateID,
			"tracking_id":  trackingID,
		}
		if fulfillment != nil {
			if fulfillment.Type != "" {
				record["fulfillment_type"] = fulfillment.Type
			}
			if fulfillment.Contact.FullName != "" {
				record["customer_name"] = fulfillment.Contact.FullName
			}
			if fulfillment.Contact.Phone != "" {
				record["customer_phone"] = fulfillment.Contact.Phone
			}
			if fulfillment.Contact.Email != "" {
				record["customer_email"] = fulfillment.Contact.Email
			}
			if fulfillment.Address.RawAddress != "" {
				record["delivery_address"] = fulfillment.Address.RawAddress
			}
			if fulfillment.Address.DeliveryDetails.Unit != "" {
				record["delivery_unit"] = fulfillment.Address.DeliveryDetails.Unit
			}
			if fulfillment.Address.DeliveryDetails.Notes != "" {
				record["delivery_notes"] = fulfillment.Address.DeliveryDetails.Notes
			}
			lat, lon := fulfillment.Address.Coordinates.Latitude, fulfillment.Address.Coordinates.Longitude
			if lat != 0 || lon != 0 {
				record["coordinates_latitude"] = lat
				record["coordinates_longitude"] = lon
			}
		}
		_, _, err := supabase.From("order_tracking").
			Insert(record, false, "", "", "").
			Execute()
		if err != nil {
			// Código PostgreSQL unique_violation = 23505
			if strings.Contains(err.Error(), "23505") || strings.Contains(strings.ToLower(err.Error()), "unique") {
				return ErrTrackingIDConflict
			}
			return fmt.Errorf("inserting order_tracking: %w", err)
		}
		return nil
	}
}

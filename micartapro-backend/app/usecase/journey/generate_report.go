package journey

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// Timezone por defecto para reportes (Chile). Si es "", se usa UTC.
const defaultReportTimezone = "America/Santiago"

// GenerateJourneyReportXLSX genera un XLSX con las órdenes de la jornada.
// openedAt, closedAt: se convierten al timezone del negocio (America/Santiago) para que
// quien abra el Excel en esa zona vea la hora local. Creado se formatea igual si viene en UTC.
func GenerateJourneyReportXLSX(
	items []ReportOrderItem,
	openedAt, closedAt time.Time,
) ([]byte, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	loc, err := time.LoadLocation(defaultReportTimezone)
	if err != nil {
		loc = time.UTC
	}
	openedLocal := openedAt.In(loc)
	closedLocal := closedAt.In(loc)
	openedStr := openedLocal.Format("02/01/2006 15:04")
	closedStr := closedLocal.Format("02/01/2006 15:04")

	tzLabel := " (Chile)"
	if loc == time.UTC {
		tzLabel = " (UTC)"
	}
	headers := []string{"Orden", "Ítem", "Cantidad", "Unidad", "Estación", "Tipo", "Estado", "Precio", "Costo", "Creado" + tzLabel, "Apertura" + tzLabel, "Cierre" + tzLabel}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	row := 2
	for _, it := range items {
		orderNum := ""
		if it.OrderNumber != nil {
			orderNum = fmt.Sprintf("%d", *it.OrderNumber)
		}
		station := ""
		if it.Station != nil {
			station = *it.Station
		}
		createdStr := it.CreatedAt
		if t, err := time.Parse(time.RFC3339Nano, it.CreatedAt); err == nil {
			createdStr = t.In(loc).Format("02/01/2006 15:04")
		} else if t, err := time.Parse(time.RFC3339, it.CreatedAt); err == nil {
			createdStr = t.In(loc).Format("02/01/2006 15:04")
		}
		vals := []interface{}{orderNum, it.ItemName, it.Quantity, it.Unit, station, it.Fulfillment, it.Status, it.TotalPrice, it.TotalCost, createdStr, openedStr, closedStr}
		for i, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			f.SetCellValue(sheet, cell, v)
		}
		row++
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("writing xlsx: %w", err)
	}
	return buf.Bytes(), nil
}

// ReportOrderItem es un ítem para el reporte (desde order_items_projection).
type ReportOrderItem struct {
	OrderNumber   *int
	ItemName      string
	Quantity      int
	Unit          string
	Station       *string
	Fulfillment   string
	Status        string
	RequestedTime *string
	CreatedAt     string
	TotalPrice    float64
	TotalCost     float64
}

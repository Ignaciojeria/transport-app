package domain

import "time"

type OrderSearchFilters struct {
	Pagination                          Pagination
	Organization                        Organization
	ReferenceIDs                        []string
	Lpns                                []string
	Commerces                           []string `json:"commerce"`
	OrderSearchOperatorDailyPlanFilters OrderSearchOperatorDailyPlanFilters
}

type OrderSearchOperatorDailyPlanFilters struct {
	OperatorReferenceID string
	PlannedDate         time.Time
}

func (f OrderSearchOperatorDailyPlanFilters) GetRouteReferenceID() string {
	// Formatear la fecha como "yyyy-mm-dd"
	formattedDate := f.PlannedDate.Format("2006-01-02")
	// Concatenar el OperatorReferenceID con la fecha formateada
	return formattedDate + "_" + f.OperatorReferenceID
}

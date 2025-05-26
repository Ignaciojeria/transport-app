package projectionresult

type DeliveryUnitsProjectionResult struct {
	ID                                    int64  `json:"id"`
	Channel                               string `json:"channel"`
	OrderReferenceID                      string `json:"order_reference_id"`
	OrderCollectAvailabilityDate          string `json:"order_collect_availability_date"`
	OrderCollectAvailabilityDateStartTime string `json:"order_collect_availability_date_start_time"`
	OrderCollectAvailabilityDateEndTime   string `json:"order_collect_availability_date_end_time"`
	OrderPromisedDateStartDate            string `json:"order_promised_date_start_date"`
	OrderPromisedDateEndDate              string `json:"order_promised_date_end_date"`
	OrderPromisedDateStartTime            string `json:"order_promised_date_start_time"`
	OrderPromisedDateEndTime              string `json:"order_promised_date_end_time"`
	OrderPromisedDateServiceCategory      string `json:"order_promised_date_service_category"`
}

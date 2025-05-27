package projectionresult

type DeliveryUnitsProjectionResult struct {
	ID                                    int64  `json:"id"`
	Channel                               string `json:"channel"`
	Consumer                              string `json:"order_consumer"`
	Commerce                              string `json:"order_commerce"`
	OrderDeliveryInstructions             string `json:"order_delivery_instructions"`
	OrderReferenceID                      string `json:"order_reference_id"`
	OrderCollectAvailabilityDate          string `json:"order_collect_availability_date"`
	OrderCollectAvailabilityDateStartTime string `json:"order_collect_availability_date_start_time"`
	OrderCollectAvailabilityDateEndTime   string `json:"order_collect_availability_date_end_time"`
	OrderPromisedDateStartDate            string `json:"order_promised_date_start_date"`
	OrderPromisedDateEndDate              string `json:"order_promised_date_end_date"`
	OrderPromisedDateStartTime            string `json:"order_promised_date_start_time"`
	OrderPromisedDateEndTime              string `json:"order_promised_date_end_time"`
	OrderPromisedDateServiceCategory      string `json:"order_promised_date_service_category"`

	// LPN and Package Information
	LPN            string `json:"lpn"`
	JSONDimensions string `json:"json_dimensions"`
	JSONWeight     string `json:"json_weight"`
	JSONInsurance  string `json:"json_insurance"`
	JSONItems      string `json:"json_items"`

	// Destination Address Information
	DestinationAddressLine1         string  `json:"destination_address_line1"`
	DestinationAddressLine2         string  `json:"destination_address_line2"`
	DestinationDistrict             string  `json:"destination_district"`
	DestinationLatitude             float64 `json:"destination_latitude"`
	DestinationLongitude            float64 `json:"destination_longitude"`
	DestinationProvince             string  `json:"destination_province"`
	DestinationState                string  `json:"destination_state"`
	DestinationTimeZone             string  `json:"destination_time_zone"`
	DestinationZipCode              string  `json:"destination_zip_code"`
	DestinationRequiresManualReview bool    `json:"destination_requires_manual_review"`

	// Destination Contact Information
	DestinationContactEmail             string `json:"destination_contact_email"`
	DestinationContactFullName          string `json:"destination_contact_full_name"`
	DestinationContactNationalID        string `json:"destination_contact_national_id"`
	DestinationContactPhone             string `json:"destination_contact_phone"`
	DestinationContactDocuments         string `json:"destination_contact_documents"`
	DestinationAdditionalContactMethods string `json:"destination_additional_contact_methods"`
}

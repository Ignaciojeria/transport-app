package projectionresult

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"
	"transport-app/app/adapter/out/tidbrepository/table"
)

// DeliveryUnitSkills is a custom type to handle JSON array scanning
type DeliveryUnitSkills []string

// Value implements the driver.Valuer interface
func (s DeliveryUnitSkills) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *DeliveryUnitSkills) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("type assertion to []byte or string failed")
	}
}

type DeliveryUnitsProjectionResult struct {
	ID                                    int64               `json:"id"`
	Status                                string              `json:"status"`
	Channel                               string              `json:"channel"`
	Consumer                              string              `json:"order_consumer"`
	Commerce                              string              `json:"order_commerce"`
	OrderGroupByType                      string              `json:"order_group_by_type"`
	OrderGroupByValue                     string              `json:"order_group_by_value"`
	OrderDeliveryInstructions             string              `json:"order_delivery_instructions"`
	OrderReferenceID                      string              `json:"order_reference_id"`
	OrderReferences                       table.JSONReference `json:"order_references" gorm:"column:order_references;type:jsonb"`
	OrderCollectAvailabilityDate          string              `json:"order_collect_availability_date"`
	OrderCollectAvailabilityDateStartTime string              `json:"order_collect_availability_date_start_time"`
	OrderCollectAvailabilityDateEndTime   string              `json:"order_collect_availability_date_end_time"`
	OrderPromisedDateStartDate            string              `json:"order_promised_date_start_date"`
	OrderPromisedDateEndDate              string              `json:"order_promised_date_end_date"`
	OrderPromisedDateStartTime            string              `json:"order_promised_date_start_time"`
	OrderPromisedDateEndTime              string              `json:"order_promised_date_end_time"`
	OrderPromisedDateServiceCategory      string              `json:"order_promised_date_service_category"`
	ManualChangePerformedBy               string              `json:"manual_change_performed_by"`
	ManualChangeReason                    string              `json:"manual_change_reason"`
	DeliveryUnitSkills                    DeliveryUnitSkills  `json:"delivery_unit_skills" gorm:"type:jsonb"`

	// Delivery Failure Information
	NonDeliveryReasonReferenceID string `json:"non_delivery_reason_reference_id"`
	NonDeliveryReason            string `json:"non_delivery_reason"`
	NonDeliveryDetail            string `json:"non_delivery_detail"`

	// Evidence Photos Information
	EvidencePhotos table.JSONEvidencePhotos `json:"evidence_photos" gorm:"type:json;"`

	// OrderType Information
	OrderType            string `json:"order_type"`
	OrderTypeDescription string `json:"order_type_description"`

	// LPN and Package Information
	LPN                string               `json:"lpn"`
	JSONDimensions     table.JSONDimensions `json:"json_dimensions"`
	JSONWeight         table.JSONWeight     `json:"json_weight"`
	JSONInsurance      table.JSONInsurance  `json:"json_insurance"`
	JSONItems          table.JSONItems      `json:"json_items"`
	DeliveryUnitLabels table.JSONReference  `json:"delivery_unit_labels" gorm:"column:delivery_unit_labels;type:jsonb"`

	SizeCategory string `json:"size_category"`

	// Origin Address Information
	OriginAddressLine1 string `json:"origin_address_line1"`
	OriginAddressLine2 string `json:"origin_address_line2"`
	OriginTimeZone     string `json:"origin_time_zone"`
	OriginZipCode      string `json:"origin_zip_code"`

	// Origin Political Area Information
	OriginPoliticalAreaCode              string  `json:"origin_political_area_code"`
	OriginPoliticalAreaConfidenceLevel   float64 `json:"origin_political_area_confidence_level"`
	OriginPoliticalAreaConfidenceMessage string  `json:"origin_political_area_confidence_message"`
	OriginPoliticalAreaConfidenceReason  string  `json:"origin_political_area_confidence_reason"`

	// Origin Coordinates Information
	OriginCoordinatesLatitude          float64 `json:"origin_coordinates_latitude"`
	OriginCoordinatesLongitude         float64 `json:"origin_coordinates_longitude"`
	OriginCoordinatesSource            string  `json:"origin_coordinates_source"`
	OriginCoordinatesConfidenceLevel   float64 `json:"origin_coordinates_confidence_level"`
	OriginCoordinatesConfidenceMessage string  `json:"origin_coordinates_confidence_message"`
	OriginCoordinatesConfidenceReason  string  `json:"origin_coordinates_confidence_reason"`

	// Origin Contact Information
	OriginContactEmail             string              `json:"origin_contact_email"`
	OriginContactFullName          string              `json:"origin_contact_full_name"`
	OriginContactNationalID        string              `json:"origin_contact_national_id"`
	OriginContactPhone             string              `json:"origin_contact_phone"`
	OriginContactDocuments         table.JSONReference `json:"origin_contact_documents"`
	OriginAdditionalContactMethods table.JSONReference `json:"origin_additional_contact_methods"`

	// Destination Address Information
	DestinationAddressLine1 string `json:"destination_address_line1"`
	DestinationAddressLine2 string `json:"destination_address_line2"`
	DestinationTimeZone     string `json:"destination_time_zone"`
	DestinationZipCode      string `json:"destination_zip_code"`

	// Destination Political Area Information
	DestinationPoliticalAreaCode              string  `json:"destination_political_area_code"`
	DestinationPoliticalAreaConfidenceLevel   float64 `json:"destination_political_area_confidence_level"`
	DestinationPoliticalAreaConfidenceMessage string  `json:"destination_political_area_confidence_message"`
	DestinationPoliticalAreaConfidenceReason  string  `json:"destination_political_area_confidence_reason"`

	// Destination Coordinates Information
	DestinationCoordinatesLatitude          float64 `json:"destination_coordinates_latitude"`
	DestinationCoordinatesLongitude         float64 `json:"destination_coordinates_longitude"`
	DestinationCoordinatesSource            string  `json:"destination_coordinates_source"`
	DestinationCoordinatesConfidenceLevel   float64 `json:"destination_coordinates_confidence_level"`
	DestinationCoordinatesConfidenceMessage string  `json:"destination_coordinates_confidence_message"`
	DestinationCoordinatesConfidenceReason  string  `json:"destination_coordinates_confidence_reason"`

	// Destination Contact Information
	DestinationContactEmail             string              `json:"destination_contact_email"`
	DestinationContactFullName          string              `json:"destination_contact_full_name"`
	DestinationContactNationalID        string              `json:"destination_contact_national_id"`
	DestinationContactPhone             string              `json:"destination_contact_phone"`
	DestinationContactDocuments         table.JSONReference `json:"destination_contact_documents"`
	DestinationAdditionalContactMethods table.JSONReference `json:"destination_additional_contact_methods"`

	// Extra Fields and Group By
	ExtraFields table.JSONMap `json:"extra_fields" gorm:"column:extra_fields;type:jsonb"`

	DestinationAdminAreaLevel1 string `json:"destination_admin_area_level1" gorm:"column:destination_admin_area_level1"`
	DestinationAdminAreaLevel2 string `json:"destination_admin_area_level2" gorm:"column:destination_admin_area_level2"`
	DestinationAdminAreaLevel3 string `json:"destination_admin_area_level3" gorm:"column:destination_admin_area_level3"`
	DestinationAdminAreaLevel4 string `json:"destination_admin_area_level4" gorm:"column:destination_admin_area_level4"`
	OriginAdminAreaLevel1      string `json:"origin_admin_area_level1" gorm:"column:origin_admin_area_level1"`
	OriginAdminAreaLevel2      string `json:"origin_admin_area_level2" gorm:"column:origin_admin_area_level2"`
	OriginAdminAreaLevel3      string `json:"origin_admin_area_level3" gorm:"column:origin_admin_area_level3"`
	OriginAdminAreaLevel4      string `json:"origin_admin_area_level4" gorm:"column:origin_admin_area_level4"`
}

type DeliveryUnitsProjectionResults []DeliveryUnitsProjectionResult

func (r DeliveryUnitsProjectionResults) Reversed() DeliveryUnitsProjectionResults {
	copied := make(DeliveryUnitsProjectionResults, len(r))
	copy(copied, r)
	slices.Reverse(copied)
	return copied
}

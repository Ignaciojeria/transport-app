package projectionresult

import (
	"slices"
	"transport-app/app/adapter/out/tidbrepository/table"
)

type NodesProjectionResult struct {
	ID int64 `json:"id"`

	NodeName       string              `json:"node_name"`
	NodeType       string              `json:"node_type"`
	NodeReferences table.JSONReference `json:"node_references"`

	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	District     string `json:"district"`
	Province     string `json:"province"`
	State        string `json:"state"`
	TimeZone     string `json:"time_zone"`
	ZipCode      string `json:"zip_code"`

	// Origin Coordinates Information
	CoordinatesLatitude          float64 `json:"coordinates_latitude"`
	CoordinatesLongitude         float64 `json:"coordinates_longitude"`
	CoordinatesSource            string  `json:"coordinates_source"`
	CoordinatesConfidenceLevel   float64 `json:"coordinates_confidence_level"`
	CoordinatesConfidenceMessage string  `json:"coordinates_confidence_message"`
	CoordinatesConfidenceReason  string  `json:"coordinates_confidence_reason"`

	// Origin Contact Information
	ContactEmail             string              `json:"contact_email"`
	ContactFullName          string              `json:"contact_full_name"`
	ContactNationalID        string              `json:"contact_national_id"`
	ContactPhone             string              `json:"contact_phone"`
	ContactDocuments         table.JSONReference `json:"contact_documents"`
	AdditionalContactMethods table.JSONReference `json:"additional_contact_methods"`
}

type NodesProjectionResults []NodesProjectionResult

func (r NodesProjectionResults) Reversed() NodesProjectionResults {
	copied := make(NodesProjectionResults, len(r))
	copy(copied, r)
	slices.Reverse(copied)
	return copied
}

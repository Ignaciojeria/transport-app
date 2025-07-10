package mapper

import (
	"context"
	"strconv"
	"time"
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
)

// formatDate extrae la fecha en formato YYYY-MM-DD de un string en formato RFC3339
func formatDate(dateStr string) *string {
	if dateStr == "" {
		return nil
	}
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return &dateStr
	}
	formatted := date.Format("2006-01-02")
	return &formatted
}

func MapDeliveryUnits(ctx context.Context, deliveryUnits []projectionresult.DeliveryUnitsProjectionResult) []*model.DeliveryUnitsReport {
	deliveryUnitsReport := make([]*model.DeliveryUnitsReport, len(deliveryUnits))

	for i, du := range deliveryUnits {
		report := &model.DeliveryUnitsReport{
			ID:                   strconv.FormatInt(du.ID, 10),
			Commerce:             &du.Commerce,
			Consumer:             &du.Consumer,
			Channel:              &du.Channel,
			Status:               &du.Status,
			DeliveryInstructions: &du.OrderDeliveryInstructions,
			ReferenceID:          du.OrderReferenceID,
			GroupBy: &model.GroupBy{
				Type:  &du.OrderGroupByType,
				Value: &du.OrderGroupByValue,
			},
			CollectAvailabilityDate: &model.CollectAvailabilityDate{
				Date: formatDate(du.OrderCollectAvailabilityDate),
				TimeRange: &model.TimeRange{
					StartTime: &du.OrderCollectAvailabilityDateStartTime,
					EndTime:   &du.OrderCollectAvailabilityDateEndTime,
				},
			},
			Delivery: &model.Delivery{
				Recipient: &model.DeliveryRecipient{
					FullName:   &du.DestinationContactFullName,
					NationalID: &du.DestinationContactNationalID,
				},

				Failure: &model.DeliveryFailure{
					ReferenceID: &du.NonDeliveryReasonReferenceID,
					Reason:      &du.NonDeliveryReason,
					Detail:      &du.NonDeliveryDetail,
				},
				EvidencePhotos: func() []*model.EvidencePhoto {
					if du.EvidencePhotos == nil {
						return []*model.EvidencePhoto{}
					}
					photos := make([]*model.EvidencePhoto, len(du.EvidencePhotos))
					for i, photo := range du.EvidencePhotos {
						takenAt := photo.TakenAt.Format(time.RFC3339)
						photos[i] = &model.EvidencePhoto{
							TakenAt: &takenAt,
							Type:    &photo.Type,
							URL:     &photo.URL,
						}
					}
					return photos
				}(),
			},
			ManualChange: &model.ManualChange{
				PerformedBy: &du.ManualChangePerformedBy,
				Reason:      &du.ManualChangeReason,
			},
			Destination: &model.Location{
				AddressInfo: &model.AddressInfo{
					AddressLine1: &du.DestinationAddressLine1,
					AddressLine2: &du.DestinationAddressLine2,
					ZipCode:      &du.DestinationZipCode,
					PoliticalArea: &model.PoliticalArea{
						Code:            &du.DestinationPoliticalAreaCode,
						AdminAreaLevel1: &du.DestinationAdminAreaLevel1,
						AdminAreaLevel2: &du.DestinationAdminAreaLevel2,
						AdminAreaLevel3: &du.DestinationAdminAreaLevel3,
						AdminAreaLevel4: &du.DestinationAdminAreaLevel4,
						TimeZone:        &du.DestinationTimeZone,
						Confidence: &model.Confidence{
							Level:   &du.DestinationPoliticalAreaConfidenceLevel,
							Message: &du.DestinationPoliticalAreaConfidenceMessage,
							Reason:  &du.DestinationPoliticalAreaConfidenceReason,
						},
					},
					Coordinates: &model.Coordinates{
						Latitude:  &du.DestinationCoordinatesLatitude,
						Longitude: &du.DestinationCoordinatesLongitude,
						Source:    &du.DestinationCoordinatesSource,
						Confidence: &model.Confidence{
							Level:   &du.DestinationCoordinatesConfidenceLevel,
							Message: &du.DestinationCoordinatesConfidenceMessage,
							Reason:  &du.DestinationCoordinatesConfidenceReason,
						},
					},
					Contact: &model.Contact{
						Email:      &du.DestinationContactEmail,
						FullName:   &du.DestinationContactFullName,
						NationalID: &du.DestinationContactNationalID,
						Phone:      &du.DestinationContactPhone,
						Documents: func() []*model.Document {
							if du.DestinationContactDocuments == nil {
								return []*model.Document{}
							}
							documents := make([]*model.Document, len(du.DestinationContactDocuments))
							for j, doc := range du.DestinationContactDocuments {
								documents[j] = &model.Document{
									Type:  &doc.Type,
									Value: &doc.Value,
								}
							}
							return documents
						}(),
						AdditionalContactMethods: func() []*model.ContactMethod {
							if du.DestinationAdditionalContactMethods == nil {
								return []*model.ContactMethod{}
							}
							methods := make([]*model.ContactMethod, len(du.DestinationAdditionalContactMethods))
							for j, method := range du.DestinationAdditionalContactMethods {
								methods[j] = &model.ContactMethod{
									Type:  &method.Type,
									Value: &method.Value,
								}
							}
							return methods
						}(),
					},
				},
			},
			Origin: &model.Location{
				AddressInfo: &model.AddressInfo{
					AddressLine1: &du.OriginAddressLine1,
					AddressLine2: &du.OriginAddressLine2,
					ZipCode:      &du.OriginZipCode,
					PoliticalArea: &model.PoliticalArea{
						Code:            &du.OriginPoliticalAreaCode,
						AdminAreaLevel1: &du.OriginAdminAreaLevel1,
						AdminAreaLevel2: &du.OriginAdminAreaLevel2,
						AdminAreaLevel3: &du.OriginAdminAreaLevel3,
						AdminAreaLevel4: &du.OriginAdminAreaLevel4,
						TimeZone:        &du.OriginTimeZone,
						Confidence: &model.Confidence{
							Level:   &du.OriginPoliticalAreaConfidenceLevel,
							Message: &du.OriginPoliticalAreaConfidenceMessage,
							Reason:  &du.OriginPoliticalAreaConfidenceReason,
						},
					},
					Coordinates: &model.Coordinates{
						Latitude:  &du.OriginCoordinatesLatitude,
						Longitude: &du.OriginCoordinatesLongitude,
						Source:    &du.OriginCoordinatesSource,
						Confidence: &model.Confidence{
							Level:   &du.OriginCoordinatesConfidenceLevel,
							Message: &du.OriginCoordinatesConfidenceMessage,
							Reason:  &du.OriginCoordinatesConfidenceReason,
						},
					},
					Contact: &model.Contact{
						Email:      &du.OriginContactEmail,
						FullName:   &du.OriginContactFullName,
						NationalID: &du.OriginContactNationalID,
						Phone:      &du.OriginContactPhone,
						Documents: func() []*model.Document {
							if du.OriginContactDocuments == nil {
								return []*model.Document{}
							}
							documents := make([]*model.Document, len(du.OriginContactDocuments))
							for j, doc := range du.OriginContactDocuments {
								documents[j] = &model.Document{
									Type:  &doc.Type,
									Value: &doc.Value,
								}
							}
							return documents
						}(),
						AdditionalContactMethods: func() []*model.ContactMethod {
							if du.OriginAdditionalContactMethods == nil {
								return []*model.ContactMethod{}
							}
							methods := make([]*model.ContactMethod, len(du.OriginAdditionalContactMethods))
							for j, method := range du.OriginAdditionalContactMethods {
								methods[j] = &model.ContactMethod{
									Type:  &method.Type,
									Value: &method.Value,
								}
							}
							return methods
						}(),
					},
				},
			},
			OrderType: &model.OrderType{
				Type:        &du.OrderType,
				Description: &du.OrderTypeDescription,
			},
			DeliveryUnit: &model.DeliveryUnit{
				Lpn: &du.LPN,
				Skills: func() []*string {
					if du.DeliveryUnitSkills == nil {
						return []*string{}
					}
					if len(du.DeliveryUnitSkills) == 1 && du.DeliveryUnitSkills[0] == "" {
						return []*string{}
					}
					skills := make([]*string, len(du.DeliveryUnitSkills))
					for i, skill := range du.DeliveryUnitSkills {
						skillCopy := skill
						skills[i] = &skillCopy
					}
					return skills
				}(),
				SizeCategory: &du.SizeCategory,
				Dimensions: &model.Dimension{
					Length: &du.JSONDimensions.Length,
					Width:  &du.JSONDimensions.Width,
					Height: &du.JSONDimensions.Height,
					Unit:   &du.JSONDimensions.Unit,
				},
				Insurance: &model.Insurance{
					UnitValue: &du.JSONInsurance.UnitValue,
					Currency:  &du.JSONInsurance.Currency,
				},
				Items: func() []*model.Item {
					if du.JSONItems == nil {
						return []*model.Item{}
					}
					if len(du.JSONItems) == 1 && du.JSONItems[0].Sku != "" &&
						du.JSONItems[0].Description == "" &&
						du.JSONItems[0].QuantityNumber == 0 &&
						du.JSONItems[0].QuantityUnit == "" {
						return []*model.Item{}
					}
					items := make([]*model.Item, len(du.JSONItems))
					for i, item := range du.JSONItems {
						items[i] = &model.Item{
							Sku:         &item.Sku,
							Description: &item.Description,
							Quantity: &model.Quantity{
								QuantityNumber: &item.QuantityNumber,
								QuantityUnit:   &item.QuantityUnit,
							},
							Dimensions: &model.Dimension{
								Length: &item.JSONDimensions.Length,
								Width:  &item.JSONDimensions.Width,
								Height: &item.JSONDimensions.Height,
								Unit:   &item.JSONDimensions.Unit,
							},
							Insurance: &model.Insurance{
								UnitValue: &item.JSONInsurance.UnitValue,
								Currency:  &item.JSONInsurance.Currency,
							},
							Weight: &model.Weight{
								Unit:  &item.JSONWeight.WeightUnit,
								Value: &item.JSONWeight.WeightValue,
							},
						}
					}
					return items
				}(),
			},
			PromisedDate: &model.PromisedDate{
				ServiceCategory: &du.OrderPromisedDateServiceCategory,
				DateRange: &model.DateRange{
					StartDate: formatDate(du.OrderPromisedDateStartDate),
					EndDate:   formatDate(du.OrderPromisedDateEndDate),
				},
				TimeRange: &model.TimeRange{
					StartTime: &du.OrderPromisedDateStartTime,
					EndTime:   &du.OrderPromisedDateEndTime,
				},
			},
		}

		// Map references if they exist
		if du.OrderReferences != nil {
			references := make([]*model.Reference, len(du.OrderReferences))
			for j, ref := range du.OrderReferences {
				references[j] = &model.Reference{
					Type:  &ref.Type,
					Value: &ref.Value,
				}
			}
			report.References = references
		} else {
			report.References = []*model.Reference{}
		}

		// Map delivery unit labels if they exist
		if du.DeliveryUnitLabels != nil {
			if len(du.DeliveryUnitLabels) == 1 &&
				du.DeliveryUnitLabels[0].Type == "" &&
				du.DeliveryUnitLabels[0].Value == "" {
				report.DeliveryUnit.Labels = []*model.Label{}
			} else {
				labels := make([]*model.Label, len(du.DeliveryUnitLabels))
				for j, label := range du.DeliveryUnitLabels {
					labels[j] = &model.Label{
						Type:  &label.Type,
						Value: &label.Value,
					}
				}
				report.DeliveryUnit.Labels = labels
			}
		} else {
			report.DeliveryUnit.Labels = []*model.Label{}
		}

		// Map extra fields if they exist
		if du.ExtraFields != nil {
			extraFields := make([]*model.KeyValuePair, 0, len(du.ExtraFields))
			for key, value := range du.ExtraFields {
				extraFields = append(extraFields, &model.KeyValuePair{
					Key:   key,
					Value: value,
				})
			}
			report.ExtraFields = extraFields
		} else {
			report.ExtraFields = []*model.KeyValuePair{}
		}

		deliveryUnitsReport[i] = report
	}

	return deliveryUnitsReport
}

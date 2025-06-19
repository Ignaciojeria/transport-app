package mapper

import "transport-app/app/domain"

func MapPromisedDateToDomain(promisedDate struct {
	DateRange struct {
		EndDate   string `json:"endDate" example:"2025-03-30"`
		StartDate string `json:"startDate"  example:"2025-03-28"`
	} `json:"dateRange"`
	ServiceCategory string `json:"serviceCategory" example:"REGULAR / SAME DAY"`
	TimeRange       struct {
		EndTime   string `json:"endTime" example:"21:30"`
		StartTime string `json:"startTime" example:"10:30"`
	} `json:"timeRange"`
}) domain.PromisedDate {
	return domain.PromisedDate{
		DateRange: domain.DateRange{
			StartDate: MapDateStringToTime(promisedDate.DateRange.StartDate),
			EndDate:   MapDateStringToTime(promisedDate.DateRange.EndDate),
		},
		TimeRange: domain.TimeRange{
			StartTime: promisedDate.TimeRange.StartTime,
			EndTime:   promisedDate.TimeRange.EndTime,
		},
		ServiceCategory: promisedDate.ServiceCategory,
	}
}

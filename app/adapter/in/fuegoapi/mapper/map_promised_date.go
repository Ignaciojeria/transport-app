package mapper

import "transport-app/app/domain"

func MapPromisedDateToDomain(promisedDate struct {
	DateRange struct {
		EndDate   string `json:"endDate"`
		StartDate string `json:"startDate"`
	} `json:"dateRange"`
	ServiceCategory string `json:"serviceCategory"`
	TimeRange       struct {
		EndTime   string `json:"endTime"`
		StartTime string `json:"startTime"`
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

func MapPromisedDateFromDomain(promisedDate domain.PromisedDate) struct {
	DateRange struct {
		EndDate   string `json:"endDate"`
		StartDate string `json:"startDate"`
	} `json:"dateRange"`
	ServiceCategory string `json:"serviceCategory"`
	TimeRange       struct {
		EndTime   string `json:"endTime"`
		StartTime string `json:"startTime"`
	} `json:"timeRange"`
} {
	response := struct {
		DateRange struct {
			EndDate   string `json:"endDate"`
			StartDate string `json:"startDate"`
		} `json:"dateRange"`
		ServiceCategory string `json:"serviceCategory"`
		TimeRange       struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	}{
		ServiceCategory: promisedDate.ServiceCategory,
	}
	response.DateRange.StartDate = promisedDate.DateRange.StartDate.Format("2006-01-02")
	response.DateRange.EndDate = promisedDate.DateRange.EndDate.Format("2006-01-02")
	response.TimeRange.StartTime = promisedDate.TimeRange.StartTime
	response.TimeRange.EndTime = promisedDate.TimeRange.EndTime
	return response
}

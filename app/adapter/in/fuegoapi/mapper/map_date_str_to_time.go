package mapper

import "time"

func MapDateStringToTime(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{} // zero value como default
	}
	return date
}

package utils

import "time"

func ConvertStringToTime(date string) time.Time {
	dateRes, err := time.Parse("2006-01-02", date)

	if err != nil {
		return time.Time{}
	}
	return dateRes
}

func ConvertTimeToString(date time.Time) string {
	return date.Format("2006-01-02")
}

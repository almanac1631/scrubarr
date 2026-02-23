package util

import "time"

func MustParseDate(dateStr string) time.Time {
	date, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.UTC)
	if err != nil {
		panic(err)
	}
	return date
}

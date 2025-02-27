package utils

import "time"

func ParseTime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		panic(err)
	}
	return t
}

func ParseTimePtr(timeStr string) *time.Time {
	timeParsed := ParseTime(timeStr)
	return &timeParsed
}

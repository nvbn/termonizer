package utils

import "time"

func QuarterFromTime(t time.Time) int {
	return (int(t.Month())-1)/3 + 1
}

func WeekStart(t time.Time) time.Time {
	weekDay := t.Weekday()
	if weekDay == time.Sunday {
		weekDay = 7
	}
	weekDay -= 1

	return t.AddDate(0, 0, -int(weekDay))
}

func IgnoreTZ(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
}

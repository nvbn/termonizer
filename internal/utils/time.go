package utils

import "time"

func QuarterFromTime(t time.Time) int {
	return (int(t.Month())-1)/3 + 1
}

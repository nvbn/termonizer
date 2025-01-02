package utils

import (
	"testing"
	"time"
)

func TestQuarterFromTime(t *testing.T) {
	dateToQuarter := map[string]int{
		"2024-01-09": 1,
		"2024-04-01": 2,
		"2024-07-01": 3,
		"2024-10-01": 4,
	}

	for date, quarter := range dateToQuarter {
		t.Run(date, func(t *testing.T) {
			q, err := time.Parse("2006-01-02", date)
			if err != nil {
				t.Error("unexpected error:", err)
			}

			if actualQuarter := QuarterFromTime(q); actualQuarter != quarter {
				t.Errorf("expected quarter %d, got %d", quarter, actualQuarter)
			}
		})
	}
}

package model

import (
	"fmt"
	"testing"
	"time"
)

func TestGoals_Title(t *testing.T) {
	makeGoal := func(period Period, start string) Goals {
		parsed, err := time.Parse("2006-01-02", start)
		if err != nil {
			t.Error("unexpected error:", err)
		}

		return Goals{
			Period:  period,
			Content: "",
			Start:   parsed,
			Updated: time.Now(),
		}
	}

	goalsToExpectedTitle := map[Goals]string{
		makeGoal(Year, "2024-12-10"):    "2024",
		makeGoal(Quarter, "2024-12-10"): "2024 Q4",
		makeGoal(Week, "2024-12-10"):    "2024-12-10 (50)",
		makeGoal(Day, "2024-12-10"):     "2024-12-10 (Tuesday)",
	}

	for goal, expectedTitle := range goalsToExpectedTitle {
		t.Run(fmt.Sprintf("%+v", goal), func(t *testing.T) {
			if actualTitle := goal.Title(); actualTitle != expectedTitle {
				t.Errorf("expected title %q, got %q", expectedTitle, actualTitle)
			}
		})
	}
}

func TestGoalsRepository_FindByPeriod_Padding(t *testing.T) {
	r := NewGoalsRepository(func() time.Time {
		return time.Date(2024, 12, 10, 0, 0, 0, 0, time.Local)
	})

	periodToExpectedGoalTitle := map[Period]string{
		Year:    "2024",
		Quarter: "2024 Q4",
		Week:    "2024-12-09 (50)",
		Day:     "2024-12-10 (Tuesday)",
	}

	for period, expectedTitle := range periodToExpectedGoalTitle {
		t.Run(PeriodName(period), func(t *testing.T) {
			actualTitle, err := r.FindByPeriod(period)
			if err != nil {
				t.Error("unexpected error:", err)
			}

			if len(actualTitle) != 1 {
				t.Errorf("expected 1 goal, got %d", len(actualTitle))
			}

			if actualTitle[0].Title() != expectedTitle {
				t.Errorf("expected title %q, got %q", expectedTitle, actualTitle[0].Title())
			}
		})
	}
}

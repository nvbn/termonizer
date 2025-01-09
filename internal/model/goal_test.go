package model

import (
	"fmt"
	"testing"
	"time"
)

func TestGoals_Title(t *testing.T) {
	makeGoal := func(period Period, start string) Goal {
		parsed, err := time.Parse("2006-01-02", start)
		if err != nil {
			t.Error("unexpected error:", err)
		}

		return Goal{
			Period:  period,
			Content: "",
			Start:   parsed,
			Updated: time.Now(),
		}
	}

	goalsToExpectedTitle := map[Goal]string{
		makeGoal(Year, "2024-12-10"):    "2024",
		makeGoal(Quarter, "2024-12-10"): "2024 Q4",
		makeGoal(Week, "2024-12-10"):    "2024-12-10 W50",
		makeGoal(Day, "2024-12-10"):     "2024-12-10 Tuesday",
	}

	for goal, expectedTitle := range goalsToExpectedTitle {
		t.Run(fmt.Sprintf("%+v", goal), func(t *testing.T) {
			if actualTitle := goal.Title(); actualTitle != expectedTitle {
				t.Errorf("expected title %q, got %q", expectedTitle, actualTitle)
			}
		})
	}
}

package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/utils"
	"testing"
	"time"
)

func TestGoal_Title(t *testing.T) {
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
			if actualTitle := goal.FormatStart(); actualTitle != expectedTitle {
				t.Errorf("expected title %q, got %q", expectedTitle, actualTitle)
			}
		})
	}
}

func TestNewGoalForDay(t *testing.T) {
	testDate := time.Date(2023, 10, 2, 15, 30, 0, 0, time.UTC)

	goal := NewGoalForDay(testDate)

	if _, err := uuid.Parse(goal.ID); err != nil || goal.ID == "" {
		t.Errorf("Invalid ID: %v", goal.ID)
	}

	if goal.Period != Day {
		t.Errorf("Expected Period to be Day, got: %v", goal.Period)
	}

	if goal.Content != "" {
		t.Errorf("Expected Content to be an empty string, got: %v", goal.Content)
	}

	if !goal.Start.Equal(testDate) {
		t.Errorf("Expected Start to be %v, got: %v", testDate, goal.Start)
	}

	if !goal.Updated.Equal(testDate) {
		t.Errorf("Expected Updated to be %v, got: %v", testDate, goal.Updated)
	}
}

func TestNewGoalForWeek(t *testing.T) {
	testDate := time.Date(2023, 10, 4, 14, 0, 0, 0, time.UTC) // Wednesday, Oct 4, 2023

	goal := NewGoalForWeek(testDate)

	if _, err := uuid.Parse(goal.ID); err != nil || goal.ID == "" {
		t.Errorf("Invalid ID: %v", goal.ID)
	}

	if goal.Period != Week {
		t.Errorf("Expected Period to be Week, got: %v", goal.Period)
	}

	if goal.Content != "" {
		t.Errorf("Expected Content to be an empty string, got: %v", goal.Content)
	}

	expectedStart := utils.WeekStart(testDate)
	if !goal.Start.Equal(expectedStart) {
		t.Errorf("Expected Start to be %v, got: %v", expectedStart, goal.Start)
	}

	if !goal.Updated.Equal(testDate) {
		t.Errorf("Expected Updated to be %v, got: %v", testDate, goal.Updated)
	}
}

func TestNewGoalForQuarter(t *testing.T) {
	testDate := time.Date(2023, 8, 15, 10, 0, 0, 0, time.Local)
	goal := NewGoalForQuarter(testDate)

	if _, err := uuid.Parse(goal.ID); err != nil || goal.ID == "" {
		t.Errorf("Invalid ID: %v", goal.ID)
	}

	if goal.Period != Quarter {
		t.Errorf("Expected Period to be Quarter, got: %v", goal.Period)
	}

	if goal.Content != "" {
		t.Errorf("Expected Content to be an empty string, got: %v", goal.Content)
	}

	expectedStart := time.Date(2023, 7, 1, 0, 0, 0, 0, time.Local)
	if !goal.Start.Equal(expectedStart) {
		t.Errorf("Expected Start to be %v, got: %v", expectedStart, goal.Start)
	}

	if !goal.Updated.Equal(testDate) {
		t.Errorf("Expected Updated to be %v, got: %v", testDate, goal.Updated)
	}
}

func TestNewGoalForYear(t *testing.T) {
	testDate := time.Date(2023, 5, 20, 12, 0, 0, 0, time.Local)
	goal := NewGoalForYear(testDate)

	if _, err := uuid.Parse(goal.ID); err != nil || goal.ID == "" {
		t.Errorf("Invalid ID: %v", goal.ID)
	}

	if goal.Period != Year {
		t.Errorf("Expected Period to be Year, got: %v", goal.Period)
	}

	if goal.Content != "" {
		t.Errorf("Expected Content to be an empty string, got: %v", goal.Content)
	}

	expectedStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	if !goal.Start.Equal(expectedStart) {
		t.Errorf("Expected Start to be %v, got: %v", expectedStart, goal.Start)
	}

	if !goal.Updated.Equal(testDate) {
		t.Errorf("Expected Updated to be %v, got: %v", testDate, goal.Updated)
	}
}

func TestGoal_CompareStart_Year(t *testing.T) {
	goalYear := Goal{
		Period: Year,
		Start:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if goalYear.CompareStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) != -1 {
		t.Errorf("Expected -1 for Year comparison, got different value")
	}
	if goalYear.CompareStart(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) != 1 {
		t.Errorf("Expected 1 for Year comparison, got different value")
	}
	if goalYear.CompareStart(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)) != 0 {
		t.Errorf("Expected 0 for Year comparison, got different value")
	}
}

func TestGoal_CompareStart_Quarter(t *testing.T) {
	goalQuarter := Goal{
		Period: Quarter,
		Start:  time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
	}
	if goalQuarter.CompareStart(time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)) != -1 {
		t.Errorf("Expected -1 for Quarter comparison, got different value")
	}
	if goalQuarter.CompareStart(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)) != 1 {
		t.Errorf("Expected 1 for Quarter comparison, got different value")
	}
	if goalQuarter.CompareStart(goalQuarter.Start) != 0 {
		t.Errorf("Expected 0 for Quarter comparison, got different value")
	}
}

func TestGoal_CompareStart_Week(t *testing.T) {
	goalWeek := Goal{
		Period: Week,
		Start:  time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC), // Monday
	}
	if goalWeek.CompareStart(time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC)) != -1 {
		t.Errorf("Expected -1 for Week comparison, got different value")
	}
	if goalWeek.CompareStart(time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)) != 1 {
		t.Errorf("Expected 1 for Week comparison, got different value")
	}
	if goalWeek.CompareStart(goalWeek.Start) != 0 {
		t.Errorf("Expected 0 for Week comparison, got different value")
	}
}

func TestGoal_CompareStart_Day(t *testing.T) {
	goalDay := Goal{
		Period: Day,
		Start:  time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
	}
	if goalDay.CompareStart(time.Date(2023, 10, 11, 0, 0, 0, 0, time.UTC)) != -1 {
		t.Errorf("Expected -1 for Day comparison, got different value")
	}
	if goalDay.CompareStart(time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC)) != 1 {
		t.Errorf("Expected 1 for Day comparison, got different value")
	}
	if goalDay.CompareStart(goalDay.Start) != 0 {
		t.Errorf("Expected 0 for Day comparison, got different value")
	}
}

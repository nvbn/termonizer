package repository

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"testing"
	"time"
)

type goalsStorageMock struct{}

func (m *goalsStorageMock) ReadGoalsForPeriod(ctx context.Context, period int) ([]model.Goal, error) {
	return make([]model.Goal, 0), nil
}

func (m *goalsStorageMock) UpdateGoal(ctx context.Context, goals model.Goal) error {
	return nil
}

func (m *goalsStorageMock) CountGoalsForPeriod(ctx context.Context, period int) (int, error) {
	return 0, nil
}

// TODO: make test better
func TestGoalsRepository_FindForPeriod_Padding(t *testing.T) {
	ctx := context.Background()

	r := NewGoalsRepository(
		func() time.Time {
			return time.Date(2024, 12, 10, 0, 0, 0, 0, time.Local)
		},
		&goalsStorageMock{})

	periodToExpectedGoalTitle := map[model.Period][]string{
		model.Year:    {"2025", "2024"},
		model.Quarter: {"2025 Q1", "2024 Q4"},
		model.Week:    {"2024-12-16 W51", "2024-12-09 W50"},
		model.Day:     {"2024-12-11 Wednesday", "2024-12-10 Tuesday"},
	}

	for period, expectedTitle := range periodToExpectedGoalTitle {
		t.Run(model.PeriodName(period), func(t *testing.T) {
			actualTitle, err := r.FindForPeriod(ctx, period)
			if err != nil {
				t.Error("unexpected error:", err)
			}

			if len(actualTitle) != 2 {
				t.Errorf("expected 2 goals, got %d", len(actualTitle))
			}

			if actualTitle[0].FormatStart() != expectedTitle[0] {
				t.Errorf("expected title %q, got %q", expectedTitle[0], actualTitle[0].FormatStart())
			}

			if actualTitle[1].FormatStart() != expectedTitle[1] {
				t.Errorf("expected title %q, got %q", expectedTitle[1], actualTitle[1].FormatStart())
			}
		})
	}
}

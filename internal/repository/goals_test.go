package repository

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"testing"
	"time"
)

type storageMock struct{}

func (m *storageMock) ReadForPeriod(ctx context.Context, period int) ([]model.Goal, error) {
	return make([]model.Goal, 0), nil
}

func (m *storageMock) Update(ctx context.Context, goals model.Goal) error {
	return nil
}

func (m *storageMock) CountForPeriod(ctx context.Context, period int) (int, error) {
	return 0, nil
}

// TODO: make test better
func TestGoalsRepository_FindForPeriod_Padding(t *testing.T) {
	ctx := context.Background()

	r := NewGoalsRepository(
		func() time.Time {
			return time.Date(2024, 12, 10, 0, 0, 0, 0, time.Local)
		},
		&storageMock{})

	periodToExpectedGoalTitle := map[model.Period]string{
		model.Year:    "2024",
		model.Quarter: "2024 Q4",
		model.Week:    "2024-12-09 W50",
		model.Day:     "2024-12-10 Tuesday",
	}

	for period, expectedTitle := range periodToExpectedGoalTitle {
		t.Run(model.PeriodName(period), func(t *testing.T) {
			actualTitle, err := r.FindForPeriod(ctx, period)
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

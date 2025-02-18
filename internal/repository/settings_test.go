package repository

import (
	"context"
	"fmt"
	"github.com/nvbn/termonizer/internal/model"
	"testing"
	"time"
)

type settingsStorageMock struct{}

func (s *settingsStorageMock) ReadSettings(ctx context.Context) ([]model.Setting, error) {
	return []model.Setting{
		{ID: fmt.Sprintf("period_to_amount_%d", model.Week), Value: "12"},
	}, nil
}

func (s *settingsStorageMock) UpdateSetting(ctx context.Context, setting model.Setting) error {
	return nil
}

func TestSettings_GetAmountForPeriod(t *testing.T) {
	ctx := t.Context()

	s, err := NewSettings(ctx, time.Now, &settingsStorageMock{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedWeek := 12
	if s.GetAmountForPeriod(model.Week) != expectedWeek {
		t.Errorf("expected %d, got %d", expectedWeek, s.GetAmountForPeriod(model.Week))
	}

	expectedQuarter := defaultPeriodToAmount[model.Quarter]
	if s.GetAmountForPeriod(model.Quarter) != expectedQuarter {
		t.Errorf("expected %d, got %d", expectedQuarter, s.GetAmountForPeriod(model.Quarter))
	}
}

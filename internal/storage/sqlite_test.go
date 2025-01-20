package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestSQLite_Goals(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(ctx, ":memory:")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	defer s.Close()

	goals, err := s.ReadGoalsForPeriod(ctx, 0)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if len(goals) != 0 {
		t.Errorf("expected 0 goals, got %d", len(goals))
	}

	date, err := time.Parse("2006-01-02", "2024-12-09")
	if err != nil {
		t.Error("unexpected error:", err)
	}

	goals = []model.Goal{
		{
			ID:      uuid.New().String(),
			Period:  0,
			Content: "",
			Start:   utils.IgnoreTZ(date),
			Updated: date,
		},
		{
			ID:      uuid.New().String(),
			Period:  0,
			Content: "content",
			Start:   utils.IgnoreTZ(date),
			Updated: date,
		}}

	for _, goal := range goals {
		if err = s.UpdateGoal(ctx, goal); err != nil {
			t.Error("unexpected error:", err)
		}
	}

	goals, err = s.ReadGoalsForPeriod(ctx, 0)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if !reflect.DeepEqual(goals, goals[:1]) {
		t.Errorf("expected %v, got %v", goals[:1], goals)
	}

	amount, err := s.CountGoalsForPeriod(ctx, 0)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if amount != 1 {
		t.Errorf("expected 1, got %d", amount)
	}
}

func TestSQLite_Settings(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(ctx, ":memory:")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	defer s.Close()

	settings, err := s.ReadSettings(ctx)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if len(settings) != 0 {
		t.Errorf("expected 0 settings, got %d", len(settings))
	}

	setting := model.Setting{
		ID:      uuid.New().String(),
		Value:   "settings value",
		Updated: time.Now().UTC(),
	}

	if err := s.UpdateSetting(ctx, setting); err != nil {
		t.Error("unexpected error:", err)
	}

	settings, err = s.ReadSettings(ctx)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	expected := []model.Setting{setting}
	if !reflect.DeepEqual(expected, settings) {
		t.Errorf("expected %v, got %v", expected, settings)
	}
}

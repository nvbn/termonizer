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

func TestSQLite(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(ctx, ":memory:")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	defer s.Close()

	goals, err := s.ReadForPeriod(ctx, 0)
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
		if err = s.Update(ctx, goal); err != nil {
			t.Error("unexpected error:", err)
		}
	}

	goals, err = s.ReadForPeriod(ctx, 0)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if !reflect.DeepEqual(goals, goals[:1]) {
		t.Errorf("expected %v, got %v", goals[:1], goals)
	}

	amount, err := s.CountForPeriod(ctx, 0)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if amount != 1 {
		t.Errorf("expected 1, got %d", amount)
	}
}

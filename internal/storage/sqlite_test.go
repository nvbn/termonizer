package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/model"
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

	goals, err := s.Read(ctx)
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

	goal := model.Goal{
		ID:      uuid.New().String(),
		Period:  model.Year,
		Content: "",
		Start:   date,
		Updated: date,
	}

	if err = s.Update(ctx, goal); err != nil {
		t.Error("unexpected error:", err)
	}

	goals, err = s.Read(ctx)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if !reflect.DeepEqual(goals, []model.Goal{goal}) {
		t.Errorf("expected %v, got %v", []model.Goal{goal}, goals)
	}
}

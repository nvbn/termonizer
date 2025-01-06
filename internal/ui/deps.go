package ui

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
)

type goalsRepository interface {
	FindForPeriod(ctx context.Context, period model.Period) ([]model.Goal, error)
	CountForPeriod(ctx context.Context, period model.Period) (int, error)
	Update(ctx context.Context, goals model.Goal) error
}

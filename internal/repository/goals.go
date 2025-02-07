package repository

import (
	"context"
	"fmt"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/utils"
	"slices"
	"time"
)

type goalsStorage interface {
	ReadGoalsForPeriod(ctx context.Context, period int) ([]model.Goal, error)
	CountGoalsForPeriod(ctx context.Context, period int) (int, error)
	UpdateGoal(ctx context.Context, goals model.Goal) error
}

type Goals struct {
	timeNow func() time.Time
	storage goalsStorage
}

func NewGoalsRepository(timeNow func() time.Time, storage goalsStorage) *Goals {
	return &Goals{
		timeNow: timeNow,
		storage: storage,
	}
}

func (r *Goals) padYear(goals []model.Goal) []model.Goal {
	now := r.timeNow()
	if len(goals) == 0 || goals[0].Start.Year() < now.Year() {
		goals = slices.Insert(goals, 0, model.NewGoalForYear(now))
	}

	return goals
}

func (r *Goals) padQuarter(goals []model.Goal) []model.Goal {
	lastQuarter := -1
	if len(goals) > 0 {
		lastGoals := goals[0]
		lastQuarter = utils.QuarterFromTime(lastGoals.Start)
	}

	now := r.timeNow()
	currentQuarter := utils.QuarterFromTime(now)
	if lastQuarter < currentQuarter {
		goals = slices.Insert(goals, 0, model.NewGoalForQuarter(now))
	}

	return goals
}

func (r *Goals) padWeek(goals []model.Goal) []model.Goal {
	lastWeek := -1
	if len(goals) > 0 {
		lastGoals := goals[0]
		_, lastWeek = lastGoals.Start.ISOWeek()
	}

	now := r.timeNow()
	_, currentWeek := now.ISOWeek()
	if lastWeek < currentWeek {
		goals = slices.Insert(goals, 0, model.NewGoalForWeek(now))
	}

	return goals
}

func (r *Goals) padDay(goals []model.Goal) []model.Goal {
	now := r.timeNow()
	if len(goals) == 0 || goals[0].Start.Truncate(24*time.Hour).Before(now.Truncate(24*time.Hour)) {
		goals = slices.Insert(goals, 0, model.NewGoalForDay(now))
	}

	return goals
}

func (r *Goals) FindForPeriod(ctx context.Context, period model.Period) ([]model.Goal, error) {
	goals, err := r.storage.ReadGoalsForPeriod(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("unable to read goals: %w", err)
	}

	switch period {
	case model.Year:
		return r.padYear(goals), nil
	case model.Quarter:
		return r.padQuarter(goals), nil
	case model.Week:
		return r.padWeek(goals), nil
	case model.Day:
		return r.padDay(goals), nil
	}

	panic("Unknown period")
}

func (r *Goals) CountForPeriod(ctx context.Context, period model.Period) (int, error) {
	return r.storage.CountGoalsForPeriod(ctx, period)
}

func (r *Goals) Update(ctx context.Context, goal model.Goal) error {
	goal.Updated = r.timeNow()
	return r.storage.UpdateGoal(ctx, goal)
}

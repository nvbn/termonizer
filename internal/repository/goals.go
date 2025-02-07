package repository

import (
	"context"
	"fmt"
	"github.com/nvbn/termonizer/internal/model"
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

func (r *Goals) padYears(goals []model.Goal) []model.Goal {
	now := r.timeNow()
	if len(goals) == 0 || goals[0].CompareStart(now) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForYear(now))
	}

	nowNextYear := now.AddDate(1, 0, 0)
	if goals[0].CompareStart(nowNextYear) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForYear(nowNextYear))
	}

	return goals
}

func (r *Goals) padQuarters(goals []model.Goal) []model.Goal {
	now := r.timeNow()
	if len(goals) == 0 || goals[0].CompareStart(now) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForQuarter(now))
	}

	nowNextQuarter := now.AddDate(0, 3, 0)
	if goals[0].CompareStart(nowNextQuarter) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForQuarter(nowNextQuarter))
	}

	return goals
}

func (r *Goals) padWeeks(goals []model.Goal) []model.Goal {
	now := r.timeNow()
	if len(goals) == 0 || goals[0].CompareStart(now) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForWeek(now))
	}

	nowNextWeek := now.AddDate(0, 0, 7)
	if goals[0].CompareStart(nowNextWeek) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForWeek(nowNextWeek))
	}

	return goals
}

func (r *Goals) padDays(goals []model.Goal) []model.Goal {
	now := r.timeNow()
	if len(goals) == 0 || goals[0].CompareStart(now) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForDay(now))
	}

	nowNextDay := now.AddDate(0, 0, 1)
	if goals[0].CompareStart(nowNextDay) == -1 {
		goals = slices.Insert(goals, 0, model.NewGoalForDay(nowNextDay))
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
		return r.padYears(goals), nil
	case model.Quarter:
		return r.padQuarters(goals), nil
	case model.Week:
		return r.padWeeks(goals), nil
	case model.Day:
		return r.padDays(goals), nil
	default:
		panic("unreachable!")
	}
}

func (r *Goals) CountForPeriod(ctx context.Context, period model.Period) (int, error) {
	return r.storage.CountGoalsForPeriod(ctx, period)
}

func (r *Goals) Update(ctx context.Context, goal model.Goal) error {
	goal.Updated = r.timeNow()
	return r.storage.UpdateGoal(ctx, goal)
}

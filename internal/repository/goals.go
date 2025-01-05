package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/utils"
	"slices"
	"time"
)

type goalsStorage interface {
	ReadForPeriod(ctx context.Context, period int) ([]model.Goal, error)
	CountForPeriod(ctx context.Context, period int) (int, error)
	Update(ctx context.Context, goals model.Goal) error
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
	if len(goals) == 0 || goals[0].Start.Year() < r.timeNow().Year() {
		start := time.Date(r.timeNow().Year(), 1, 1, 0, 0, 0, 0, time.Local)
		goals = slices.Insert(goals, 0, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Year,
			Content: "",
			Start:   start,
			Updated: r.timeNow(),
		})
	}

	return goals
}

func (r *Goals) padQuarter(goals []model.Goal) []model.Goal {
	lastQuarter := -1
	if len(goals) > 0 {
		lastGoals := goals[0]
		lastQuarter = utils.QuarterFromTime(lastGoals.Start)
	}

	currentQuarter := utils.QuarterFromTime(r.timeNow())

	if lastQuarter < currentQuarter {
		currentQuarterStartDate := time.Date(r.timeNow().Year(), time.Month(currentQuarter*3-2), 1, 0, 0, 0, 0, time.Local)
		goals = slices.Insert(goals, 0, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Quarter,
			Content: "",
			Start:   currentQuarterStartDate,
			Updated: r.timeNow(),
		})
	}

	return goals
}

func (r *Goals) padWeek(goals []model.Goal) []model.Goal {
	lastWeek := -1
	if len(goals) > 0 {
		lastGoals := goals[0]
		_, lastWeek = lastGoals.Start.ISOWeek()
	}

	_, currentWeek := r.timeNow().ISOWeek()
	if lastWeek < currentWeek {
		goals = slices.Insert(goals, 0, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Week,
			Content: "",
			Start:   utils.WeekStart(r.timeNow()),
			Updated: r.timeNow(),
		})
	}

	return goals
}

func (r *Goals) padDay(goals []model.Goal) []model.Goal {
	if len(goals) == 0 || goals[0].Start.Day() < r.timeNow().Day() {
		goals = slices.Insert(goals, 0, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Day,
			Content: "",
			Start:   r.timeNow(),
			Updated: r.timeNow(),
		})
	}

	return goals
}

func (r *Goals) FindForPeriod(ctx context.Context, period model.Period) ([]model.Goal, error) {
	goals, err := r.storage.ReadForPeriod(ctx, period)
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
	return r.storage.CountForPeriod(ctx, period)
}

func (r *Goals) Update(ctx context.Context, goals model.Goal) error {
	return r.storage.Update(ctx, goals)
}

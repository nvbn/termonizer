package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/utils"
	"time"
)

type goalsStorage interface {
	Read(ctx context.Context) ([]model.Goal, error)
	Update(ctx context.Context, goals model.Goal) error
}

type Goals struct {
	timeNow  func() time.Time
	storage  goalsStorage
	byPeriod map[model.Period][]model.Goal
}

func NewGoalsRepository(ctx context.Context, timeNow func() time.Time, storage goalsStorage) (*Goals, error) {
	// make it explicit?
	goals, err := storage.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read goals: %w", err)
	}

	byPeriod := make(map[model.Period][]model.Goal)
	for _, goal := range goals {
		byPeriod[goal.Period] = append(byPeriod[goal.Period], goal)
	}

	return &Goals{
		timeNow:  timeNow,
		storage:  storage,
		byPeriod: byPeriod,
	}, nil
}

func (r *Goals) padYear() error {
	if len(r.byPeriod[model.Year]) == 0 || r.byPeriod[model.Year][len(r.byPeriod[model.Year])-1].Start.Year() < r.timeNow().Year() {
		start, err := time.Parse("2006", r.timeNow().Format("2006"))
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}
		r.byPeriod[model.Year] = append(r.byPeriod[model.Year], model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Year,
			Content: "",
			Start:   start,
			Updated: r.timeNow(),
		})
	}

	return nil
}

func (r *Goals) padQuarter() {
	lastQuarter := -1
	if len(r.byPeriod[model.Quarter]) > 0 {
		lastGoals := r.byPeriod[model.Quarter][len(r.byPeriod[model.Quarter])-1]
		lastQuarter = utils.QuarterFromTime(lastGoals.Start)
	}

	currentQuarter := utils.QuarterFromTime(r.timeNow())

	if lastQuarter < currentQuarter {
		currentQuarterStartDate := time.Date(r.timeNow().Year(), time.Month(currentQuarter*3-2), 1, 0, 0, 0, 0, time.Local)
		r.byPeriod[model.Quarter] = append(r.byPeriod[model.Quarter], model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Quarter,
			Content: "",
			Start:   currentQuarterStartDate,
			Updated: r.timeNow(),
		})
	}
}

func (r *Goals) padWeek() {
	lastWeek := -1
	if len(r.byPeriod[model.Week]) > 0 {
		lastGoals := r.byPeriod[model.Week][len(r.byPeriod[model.Week])-1]
		_, lastWeek = lastGoals.Start.ISOWeek()
	}

	_, currentWeek := r.timeNow().ISOWeek()
	if lastWeek < currentWeek {
		r.byPeriod[model.Week] = append(r.byPeriod[model.Week], model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Week,
			Content: "",
			Start:   utils.WeekStart(r.timeNow()),
			Updated: r.timeNow(),
		})
	}
}

func (r *Goals) padDay() {
	if len(r.byPeriod[model.Day]) == 0 || r.byPeriod[model.Day][len(r.byPeriod[model.Day])-1].Start.Day() < r.timeNow().Day() {
		r.byPeriod[model.Day] = append(r.byPeriod[model.Day], model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Day,
			Content: "",
			Start:   r.timeNow(),
			Updated: r.timeNow(),
		})
	}
}

func (r *Goals) FindByPeriod(period model.Period) ([]model.Goal, error) {
	switch period {
	case model.Year:
		if err := r.padYear(); err != nil {
			return nil, err
		}
		return r.byPeriod[period], nil
	case model.Quarter:
		r.padQuarter()
		return r.byPeriod[period], nil
	case model.Week:
		r.padWeek()
		return r.byPeriod[period], nil
	case model.Day:
		r.padDay()
		return r.byPeriod[period], nil
	}

	panic("Unknown period")
}

func (r *Goals) Update(ctx context.Context, goals model.Goal) error {
	return r.storage.Update(ctx, goals)
}

package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/utils"
	"time"
)

type Goals struct {
	ID      string
	Period  Period
	Content string
	Start   time.Time
	Updated time.Time
}

func (g *Goals) Title() string {
	switch g.Period {
	case Year:
		return g.Start.Format("2006")
	case Quarter:
		year := g.Start.Format("2006")
		quarter := utils.QuarterFromTime(g.Start)
		return fmt.Sprintf("%s Q%d", year, quarter)
	case Week:
		date := g.Start.Format("2006-01-02")
		_, weekNumber := g.Start.ISOWeek()
		return fmt.Sprintf("%s (%d)", date, weekNumber)
	case Day:
		date := g.Start.Format("2006-01-02")
		weekDay := g.Start.Weekday()
		return fmt.Sprintf("%s (%s)", date, weekDay)
	}

	panic("Unknown period")
}

type goalsStorage interface {
	Read() (map[Period][]Goals, error)
	Write(map[Period][]Goals) error
}

type GoalsRepository struct {
	timeNow  func() time.Time
	storage  goalsStorage
	byPeriod map[Period][]Goals
}

func NewGoalsRepository(timeNow func() time.Time, storage goalsStorage) (*GoalsRepository, error) {
	goals, err := storage.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read goals: %w", err)
	}

	return &GoalsRepository{
		timeNow:  timeNow,
		storage:  storage,
		byPeriod: goals,
	}, nil
}

func (r *GoalsRepository) padYear() error {
	if len(r.byPeriod[Year]) == 0 || r.byPeriod[Year][len(r.byPeriod[Year])-1].Start.Year() < r.timeNow().Year() {
		start, err := time.Parse("2006", r.timeNow().Format("2006"))
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}
		r.byPeriod[Year] = append(r.byPeriod[Year], Goals{
			ID:      uuid.New().String(),
			Period:  Year,
			Content: "",
			Start:   start,
			Updated: r.timeNow(),
		})
	}

	return nil
}

func (r *GoalsRepository) padQuarter() {
	lastQuarter := -1
	if len(r.byPeriod[Quarter]) > 0 {
		lastGoals := r.byPeriod[Quarter][len(r.byPeriod[Quarter])-1]
		lastQuarter = utils.QuarterFromTime(lastGoals.Start)
	}

	currentQuarter := utils.QuarterFromTime(r.timeNow())

	if lastQuarter < currentQuarter {
		currentQuarterStartDate := time.Date(r.timeNow().Year(), time.Month(currentQuarter*3-2), 1, 0, 0, 0, 0, time.Local)
		r.byPeriod[Quarter] = append(r.byPeriod[Quarter], Goals{
			ID:      uuid.New().String(),
			Period:  Quarter,
			Content: "",
			Start:   currentQuarterStartDate,
			Updated: r.timeNow(),
		})
	}
}

func (r *GoalsRepository) padWeek() {
	lastWeek := -1
	if len(r.byPeriod[Week]) > 0 {
		lastGoals := r.byPeriod[Week][len(r.byPeriod[Week])-1]
		_, lastWeek = lastGoals.Start.ISOWeek()
	}

	_, currentWeek := r.timeNow().ISOWeek()
	if lastWeek < currentWeek {
		weekDay := r.timeNow().Weekday()
		if weekDay == time.Sunday {
			weekDay = 7
		}
		weekDay -= 1

		currentWeekStartDate := r.timeNow().AddDate(0, 0, -int(weekDay))
		r.byPeriod[Week] = append(r.byPeriod[Week], Goals{
			ID:      uuid.New().String(),
			Period:  Week,
			Content: "",
			Start:   currentWeekStartDate,
			Updated: r.timeNow(),
		})
	}
}

func (r *GoalsRepository) padDay() {
	if len(r.byPeriod[Day]) == 0 || r.byPeriod[Day][len(r.byPeriod[Day])-1].Start.Day() < r.timeNow().Day() {
		r.byPeriod[Day] = append(r.byPeriod[Day], Goals{
			ID:      uuid.New().String(),
			Period:  Day,
			Content: "",
			Start:   r.timeNow(),
			Updated: r.timeNow(),
		})
	}
}

func (r *GoalsRepository) FindByPeriod(period Period) ([]Goals, error) {
	switch period {
	case Year:
		if err := r.padYear(); err != nil {
			return nil, err
		}
		return r.byPeriod[period], nil
	case Quarter:
		r.padQuarter()
		return r.byPeriod[period], nil
	case Week:
		r.padWeek()
		return r.byPeriod[period], nil
	case Day:
		r.padDay()
		return r.byPeriod[period], nil
	}

	panic("Unknown period")
}

func (r *GoalsRepository) Sync() error {
	return r.storage.Write(r.byPeriod)
}

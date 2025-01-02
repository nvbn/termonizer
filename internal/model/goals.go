package model

import (
	"fmt"
	"github.com/nvbn/termonizer/internal/utils"
	"time"
)

type Goals struct {
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

type GoalsRepository struct {
	timeNow func() time.Time

	ByPeriod map[Period][]Goals
}

func NewGoalsRepository(timeNow func() time.Time) *GoalsRepository {
	// TODO: read from file!
	return &GoalsRepository{
		timeNow:  timeNow,
		ByPeriod: make(map[Period][]Goals),
	}
}

//func (r *GoalsRepository) allDatesForYear()

func (r *GoalsRepository) padYear() error {
	if len(r.ByPeriod[Year]) == 0 || r.ByPeriod[Year][len(r.ByPeriod[Year])-1].Start.Year() < r.timeNow().Year() {
		start, err := time.Parse("2006", r.timeNow().Format("2006"))
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}
		r.ByPeriod[Year] = append(r.ByPeriod[Year], Goals{
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
	if len(r.ByPeriod[Quarter]) > 0 {
		lastGoals := r.ByPeriod[Quarter][len(r.ByPeriod[Quarter])-1]
		lastQuarter = utils.QuarterFromTime(lastGoals.Start)
	}

	currentQuarter := utils.QuarterFromTime(r.timeNow())

	if lastQuarter < currentQuarter {
		currentQuarterStartDate := time.Date(r.timeNow().Year(), time.Month(currentQuarter*3-2), 1, 0, 0, 0, 0, time.Local)
		r.ByPeriod[Quarter] = append(r.ByPeriod[Quarter], Goals{
			Period:  Quarter,
			Content: "",
			Start:   currentQuarterStartDate,
			Updated: r.timeNow(),
		})
	}
}

func (r *GoalsRepository) padWeek() {
	lastWeek := -1
	if len(r.ByPeriod[Week]) > 0 {
		lastGoals := r.ByPeriod[Week][len(r.ByPeriod[Week])-1]
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
		r.ByPeriod[Week] = append(r.ByPeriod[Week], Goals{
			Period:  Week,
			Content: "",
			Start:   currentWeekStartDate,
			Updated: r.timeNow(),
		})
	}
}

func (r *GoalsRepository) padDay() {
	if len(r.ByPeriod[Day]) == 0 || r.ByPeriod[Day][len(r.ByPeriod[Day])-1].Start.Day() < r.timeNow().Day() {
		r.ByPeriod[Day] = append(r.ByPeriod[Day], Goals{
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
		return r.ByPeriod[period], nil
	case Quarter:
		r.padQuarter()
		return r.ByPeriod[period], nil
	case Week:
		r.padWeek()
		return r.ByPeriod[period], nil
	case Day:
		r.padDay()
		return r.ByPeriod[period], nil
	}

	panic("Unknown period")
}

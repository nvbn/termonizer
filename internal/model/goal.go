package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/utils"
	"time"
)

type Goal struct {
	ID      string
	Period  Period
	Content string
	Start   time.Time
	Updated time.Time
}

func NewGoalForDay(dt time.Time) Goal {
	return Goal{
		ID:      uuid.New().String(),
		Period:  Day,
		Content: "",
		Start:   dt,
		Updated: dt,
	}
}

func NewGoalForWeek(dt time.Time) Goal {
	return Goal{
		ID:      uuid.New().String(),
		Period:  Week,
		Content: "",
		Start:   utils.WeekStart(dt),
		Updated: dt,
	}
}

func NewGoalForQuarter(dt time.Time) Goal {
	currentQuarter := utils.QuarterFromTime(dt)
	currentQuarterStartDate := time.Date(dt.Year(), time.Month(currentQuarter*3-2), 1, 0, 0, 0, 0, time.Local)

	return Goal{
		ID:      uuid.New().String(),
		Period:  Quarter,
		Content: "",
		Start:   currentQuarterStartDate,
		Updated: dt,
	}
}

func NewGoalForYear(dt time.Time) Goal {
	start := time.Date(dt.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	return Goal{
		ID:      uuid.New().String(),
		Period:  Year,
		Content: "",
		Start:   start,
		Updated: dt,
	}
}

func (g *Goal) Title() string {
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
		return fmt.Sprintf("%s W%d", date, weekNumber)
	case Day:
		date := g.Start.Format("2006-01-02")
		weekDay := g.Start.Weekday()
		return fmt.Sprintf("%s %s", date, weekDay)
	}

	panic("Unknown period")
}

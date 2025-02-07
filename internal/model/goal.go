package model

import (
	"cmp"
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

func (g *Goal) FormatStart() string {
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

	panic("unreachable!")
}

// CompareStart look at https://pkg.go.dev/cmp, fuck it's ugly
func (g *Goal) CompareStart(dt time.Time) int {
	switch g.Period {
	case Year:
		return cmp.Compare(g.Start.Year(), dt.Year())
	case Quarter:
		compared := cmp.Compare(g.Start.Year(), dt.Year())
		if compared == 0 {
			return cmp.Compare(utils.QuarterFromTime(g.Start), utils.QuarterFromTime(dt))
		} else {
			return compared
		}
	case Week:
		goalYear, goalWeek := g.Start.ISOWeek()
		dtYear, dtWeek := dt.ISOWeek()
		compared := cmp.Compare(goalYear, dtYear)
		if compared == 0 {
			return cmp.Compare(goalWeek, dtWeek)
		} else {
			return compared
		}
	case Day:
		goalTruncated := g.Start.Truncate(24 * time.Hour)
		dtTruncated := dt.Truncate(24 * time.Hour)
		if goalTruncated.Before(dtTruncated) {
			return -1
		} else if goalTruncated.After(dtTruncated) {
			return 1
		} else {
			return 0
		}
	}

	panic("unreachable!")
}

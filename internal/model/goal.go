package model

import (
	"fmt"
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
		return fmt.Sprintf("%s (%d)", date, weekNumber)
	case Day:
		date := g.Start.Format("2006-01-02")
		weekDay := g.Start.Weekday()
		return fmt.Sprintf("%s (%s)", date, weekDay)
	}

	panic("Unknown period")
}

package model

type Period = int

const (
	Year Period = iota
	Quarter
	Week
	Day
)

var Periods = []Period{Year, Quarter, Week, Day}

func PeriodName(p Period) string {
	switch p {
	case Year:
		return "Year"
	case Quarter:
		return "Quarter"
	case Week:
		return "Week"
	case Day:
		return "Day"
	}

	panic("Unknown period")
}

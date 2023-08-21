package mastercard

import (
	"strconv"
	"strings"
	"time"
)

func TimeFormat(t time.Time, layout string) string {
	switch layout {
	case "j":
		return formatJulianDate(t)
	default:
		return t.Format(timeLayout(layout))
	}
}

func TimeParse(layout, value string) (time.Time, error) {
	switch layout {
	case "j":
		return parseJulianDate(value)
	default:
		return time.Parse(timeLayout(layout), value) //nolint:wrapcheck
	}
}

// timeLayout returns the layout for time.Format given the better known YYYYMMDD type of identifiers
func timeLayout(f string) string {
	r := strings.NewReplacer(
		"YYYY", "2006",
		"YY", "06",
		"MM", "01",
		"DD", "02",
		`hh`, `15`,
		`mm`, `04`,
		`ss`, `05`,
	)

	return r.Replace(f)
}

func formatJulianDate(value time.Time) string {
	return value.Format(`06002`)[1:]
}

func parseJulianDate(value string) (time.Time, error) {
	var decade string
	if today := time.Now().Format(`06`); value[0] <= today[1] {
		decade = today[0:1] // same decade
	} else {
		decade = strconv.Itoa(time.Now().Year() - 10)[2:3] // past decade
	}

	// prepend the decade and parse
	return time.Parse(`06002`, decade+value) //nolint:wrapcheck
}

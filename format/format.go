// Package format provides Russian date formatting helpers shared by 005-bot
// services.
package format

import (
	"fmt"
	"strings"
	"time"
)

//nolint:gochecknoglobals // static lookup table, not mutable state
var russianMonthNames = [13]string{
	time.January:   "января",
	time.February:  "февраля",
	time.March:     "марта",
	time.April:     "апреля",
	time.May:       "мая",
	time.June:      "июня",
	time.July:      "июля",
	time.August:    "августа",
	time.September: "сентября",
	time.October:   "октября",
	time.November:  "ноября",
	time.December:  "декабря",
}

// DateRU renders t in its own offset as "18 августа 14:30" using
// lowercase Russian genitive month names. No UTC/local conversion is applied;
// the timestamp's own wall clock is used.
func DateRU(t time.Time) string {
	return fmt.Sprintf("%d %s %02d:%02d", t.Day(), russianMonthNames[t.Month()], t.Hour(), t.Minute())
}

// Dates renders each timestamp with FormatDateRU and joins the results
// with spaces. An empty or nil slice produces an empty string.
func Dates(dates []time.Time) string {
	parts := make([]string, len(dates))
	for i, d := range dates {
		parts[i] = DateRU(d)
	}
	return strings.Join(parts, " ")
}

package format_test

import (
	"testing"
	"time"

	"github.com/005-bot/apis-go/format"
)

const (
	hourSeconds = 3600
	utcPlus3    = 3 * hourSeconds
	utcMinus5   = -5 * hourSeconds
)

func TestFormatDateRU(t *testing.T) {
	cases := []struct {
		name string
		in   time.Time
		want string
	}{
		{"august afternoon", time.Date(2026, time.August, 18, 14, 30, 0, 0, time.UTC), "18 августа 14:30"},
		{"single digit day and hour", time.Date(2026, time.January, 5, 9, 7, 0, 0, time.UTC), "5 января 09:07"},
		{"december evening", time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC), "31 декабря 23:59"},
		{"zero time", time.Time{}, "1 января 00:00"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := format.DateRU(tc.in); got != tc.want {
				t.Fatalf("FormatDateRU(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatDateRULowercaseMonths(t *testing.T) {
	want := map[time.Month]string{
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
	for m, name := range want {
		expected := "1 " + name + " 00:00"
		if got := format.DateRU(time.Date(2026, m, 1, 0, 0, 0, 0, time.UTC)); got != expected {
			t.Fatalf("month %v: FormatDateRU = %q, want %q", m, got, expected)
		}
	}
}

func TestFormatDateRUOwnOffset(t *testing.T) {
	instant := time.Date(2026, time.August, 18, 11, 30, 0, 0, time.UTC)

	cases := []struct {
		name string
		loc  *time.Location
		want string
	}{
		{"utc", time.UTC, "18 августа 11:30"},
		{"utc+3", time.FixedZone("UTC+3", utcPlus3), "18 августа 14:30"},
		{"utc-5", time.FixedZone("UTC-5", utcMinus5), "18 августа 06:30"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := format.DateRU(instant.In(tc.loc)); got != tc.want {
				t.Fatalf("FormatDateRU in %s = %q, want %q", tc.loc, got, tc.want)
			}
		})
	}
}

func TestFormatDates(t *testing.T) {
	dates := []time.Time{
		time.Date(2026, time.August, 18, 14, 30, 0, 0, time.UTC),
		time.Date(2026, time.August, 18, 20, 0, 0, 0, time.UTC),
	}
	want := "18 августа 14:30 18 августа 20:00"
	if got := format.Dates(dates); got != want {
		t.Fatalf("FormatDates = %q, want %q", got, want)
	}
}

func TestFormatDatesEmpty(t *testing.T) {
	if got := format.Dates(nil); got != "" {
		t.Fatalf("FormatDates(nil) = %q, want empty", got)
	}
	if got := format.Dates([]time.Time{}); got != "" {
		t.Fatalf("FormatDates(empty) = %q, want empty", got)
	}
}

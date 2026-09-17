package cmd

import (
	"fmt"
	"time"
)

// parsedTime is the result of parsing a user-supplied time value.
// HasDate reports whether the input carried its own YYYY-MM-DD date;
// when it did not, Time was stamped onto today in the local timezone.
type parsedTime struct {
	Time    time.Time
	HasDate bool
}

// timeLayout is one accepted input shape. Layouts without a timezone in the
// layout string are parsed in the local timezone.
type timeLayout struct {
	layout  string
	hasDate bool
	hasZone bool
}

// timeLayouts holds every accepted time format, most specific first. The set is
// the product of four parts: an optional YYYY-MM-DD date, a T or space
// separator, a HH:MM or HH:MM:SS time, and a timezone offset that is only
// valid after a T separator.
var timeLayouts = []timeLayout{
	{time.RFC3339, true, true},             // 2026-09-16T14:00:00Z
	{"2006-01-02T15:04Z07:00", true, true}, // 2026-09-16T14:00+02:00
	{"2006-01-02T15:04:05", true, false},   // 2026-09-16T14:00:00
	{"2006-01-02T15:04", true, false},      // 2026-09-16T14:00
	{"2006-01-02 15:04:05", true, false},   // 2026-09-16 14:00:00
	{"2006-01-02 15:04", true, false},      // 2026-09-16 14:00
	{"15:04:05", false, false},             // 14:00:00
	{"15:04", false, false},                // 14:00
}

// timeFormatHelp describes the accepted formats for error messages.
const timeFormatHelp = `Accepted formats:
  2026-09-16T14:00:00Z    date, 'T', time and timezone offset
  2026-09-16 14:00        date and time, local timezone
  14:00                   time only, today, local timezone
The date must be YYYY-MM-DD. The separator may be 'T' or a space, seconds are
optional, and a timezone offset is only allowed after a 'T' separator.`

// parseTime parses a user-supplied time value. Input without timezone
// information is interpreted in the local timezone; input without a date
// resolves to today.
func parseTime(timeStr string) (parsedTime, error) {
	for _, l := range timeLayouts {
		var t time.Time
		var err error

		if l.hasZone {
			t, err = time.Parse(l.layout, timeStr)
		} else {
			t, err = time.ParseInLocation(l.layout, timeStr, time.Local)
		}
		if err != nil {
			continue
		}

		if !l.hasDate {
			t = stampOnDate(time.Now(), t)
		}
		return parsedTime{Time: t, HasDate: l.hasDate}, nil
	}

	return parsedTime{}, fmt.Errorf("invalid time format '%s'.\n%s", timeStr, timeFormatHelp)
}

// stampOnDate returns the clock time of t on the calendar date of date, in the
// local timezone.
func stampOnDate(date, t time.Time) time.Time {
	date = date.Local()
	return time.Date(date.Year(), date.Month(), date.Day(),
		t.Hour(), t.Minute(), t.Second(), 0, time.Local)
}

// isToday reports whether t falls on today's local calendar date.
func isToday(t time.Time) bool {
	now := time.Now().Local()
	t = t.Local()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// formatEntryTime formats a timestamp for display: HH:MM when it falls on
// today, YYYY-MM-DD HH:MM otherwise.
func formatEntryTime(t time.Time) string {
	if isToday(t) {
		return t.Local().Format("15:04")
	}
	return t.Local().Format("2006-01-02 15:04")
}

// formatDuration formats duration in seconds to human-readable format
func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, mins)
}

// formatElapsedTime formats elapsed time from start to now
func formatElapsedTime(start time.Time) string {
	elapsed := int(time.Since(start).Seconds())
	return formatDuration(elapsed)
}

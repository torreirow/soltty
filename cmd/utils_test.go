package cmd

import (
	"strings"
	"testing"
	"time"
)

// todayAt returns the given clock time on today's local calendar date.
func todayAt(hour, min, sec int) time.Time {
	now := time.Now().Local()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, time.Local)
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        time.Time
		wantHasDate bool
		wantLocal   bool
	}{
		{
			name:        "RFC3339 with Z",
			input:       "2026-09-16T14:00:00Z",
			want:        time.Date(2026, 9, 16, 14, 0, 0, 0, time.UTC),
			wantHasDate: true,
		},
		{
			name:        "RFC3339 with offset",
			input:       "2026-09-16T14:00:00+02:00",
			want:        time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
			wantHasDate: true,
		},
		{
			name:        "date, T and time with offset, no seconds",
			input:       "2026-09-16T14:00+02:00",
			want:        time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
			wantHasDate: true,
		},
		{
			name:        "date, T and time with seconds, no zone",
			input:       "2026-09-16T14:00:30",
			want:        time.Date(2026, 9, 16, 14, 0, 30, 0, time.Local),
			wantHasDate: true,
			wantLocal:   true,
		},
		{
			name:        "date, T and time, no zone",
			input:       "2026-09-16T14:00",
			want:        time.Date(2026, 9, 16, 14, 0, 0, 0, time.Local),
			wantHasDate: true,
			wantLocal:   true,
		},
		{
			name:        "date, space and time with seconds",
			input:       "2026-09-16 14:00:30",
			want:        time.Date(2026, 9, 16, 14, 0, 30, 0, time.Local),
			wantHasDate: true,
			wantLocal:   true,
		},
		{
			name:        "date, space and time",
			input:       "2026-09-16 14:00",
			want:        time.Date(2026, 9, 16, 14, 0, 0, 0, time.Local),
			wantHasDate: true,
			wantLocal:   true,
		},
		{
			name:        "time only with seconds",
			input:       "14:00:30",
			want:        todayAt(14, 0, 30),
			wantHasDate: false,
			wantLocal:   true,
		},
		{
			name:        "time only",
			input:       "14:00",
			want:        todayAt(14, 0, 0),
			wantHasDate: false,
			wantLocal:   true,
		},
		{
			name:        "midnight",
			input:       "00:00",
			want:        todayAt(0, 0, 0),
			wantHasDate: false,
			wantLocal:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.input)
			if err != nil {
				t.Fatalf("parseTime(%q) returned error: %v", tt.input, err)
			}
			if !got.Time.Equal(tt.want) {
				t.Errorf("parseTime(%q) = %s, want %s", tt.input, got.Time, tt.want)
			}
			if got.HasDate != tt.wantHasDate {
				t.Errorf("parseTime(%q) HasDate = %v, want %v", tt.input, got.HasDate, tt.wantHasDate)
			}
			if tt.wantLocal && got.Time.Location() != time.Local {
				t.Errorf("parseTime(%q) location = %s, want %s", tt.input, got.Time.Location(), time.Local)
			}
		})
	}
}

func TestParseTimeRejects(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"European date notation", "16-09-2026 14:00"},
		{"US date notation", "09-16-2026 14:00"},
		{"offset after space separator", "2026-09-16 14:00+02:00"},
		{"slash-separated date", "2026/09/16 14:00"},
		{"hour without minutes", "14"},
		{"date without time", "2026-09-16"},
		{"empty string", ""},
		{"not a time at all", "yesterday"},
		{"out of range hour", "25:00"},
		{"out of range month", "2026-13-01 14:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.input)
			if err == nil {
				t.Fatalf("parseTime(%q) = %s, want error", tt.input, got.Time)
			}
		})
	}
}

func TestParseTimeErrorNamesAcceptedFormats(t *testing.T) {
	_, err := parseTime("16-09-2026 14:00")
	if err == nil {
		t.Fatal("expected an error")
	}

	msg := err.Error()
	for _, want := range []string{
		"16-09-2026 14:00",
		"2026-09-16T14:00:00Z",
		"2026-09-16 14:00",
		"14:00",
		"YYYY-MM-DD",
		"seconds are",
		"timezone offset is only allowed after",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message does not mention %q:\n%s", want, msg)
		}
	}
}

func TestStampOnDate(t *testing.T) {
	date := time.Date(2026, 9, 16, 23, 0, 0, 0, time.Local)
	clock := time.Date(0, 1, 1, 1, 30, 0, 0, time.Local)

	got := stampOnDate(date, clock)
	want := time.Date(2026, 9, 16, 1, 30, 0, 0, time.Local)

	if !got.Equal(want) {
		t.Errorf("stampOnDate = %s, want %s", got, want)
	}
}

func TestFormatEntryTime(t *testing.T) {
	today := todayAt(14, 0, 0)
	if got, want := formatEntryTime(today), "14:00"; got != want {
		t.Errorf("formatEntryTime(today) = %q, want %q", got, want)
	}

	other := today.AddDate(0, 0, -3)
	if got, want := formatEntryTime(other), other.Format("2006-01-02 15:04"); got != want {
		t.Errorf("formatEntryTime(other day) = %q, want %q", got, want)
	}
}

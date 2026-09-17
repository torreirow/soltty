// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Wouter van der Toorren

package cmd

import (
	"testing"
	"time"
)

func mustParse(t *testing.T, s string) parsedTime {
	t.Helper()
	pt, err := parseTime(s)
	if err != nil {
		t.Fatalf("parseTime(%q): %v", s, err)
	}
	return pt
}

func TestResolveEntryTimes(t *testing.T) {
	tests := []struct {
		name         string
		start, end   string
		wantStart    time.Time
		wantEnd      time.Time
		wantDuration time.Duration
		wantEndAfter bool
	}{
		{
			name:         "end without date inherits the start date",
			start:        "2026-09-16 09:00",
			end:          "17:00",
			wantStart:    time.Date(2026, 9, 16, 9, 0, 0, 0, time.Local),
			wantEnd:      time.Date(2026, 9, 16, 17, 0, 0, 0, time.Local),
			wantDuration: 8 * time.Hour,
			wantEndAfter: true,
		},
		{
			name:         "end with its own date is left alone",
			start:        "2026-09-16 23:00",
			end:          "2026-09-17 01:00",
			wantStart:    time.Date(2026, 9, 16, 23, 0, 0, 0, time.Local),
			wantEnd:      time.Date(2026, 9, 17, 1, 0, 0, 0, time.Local),
			wantDuration: 2 * time.Hour,
			wantEndAfter: true,
		},
		{
			name:         "crossing midnight without an end date is not rolled over",
			start:        "2026-09-16 23:00",
			end:          "01:00",
			wantStart:    time.Date(2026, 9, 16, 23, 0, 0, 0, time.Local),
			wantEnd:      time.Date(2026, 9, 16, 1, 0, 0, 0, time.Local),
			wantEndAfter: false,
		},
		{
			name:         "both bare times stay on today",
			start:        "09:00",
			end:          "17:00",
			wantStart:    todayAt(9, 0, 0),
			wantEnd:      todayAt(17, 0, 0),
			wantDuration: 8 * time.Hour,
			wantEndAfter: true,
		},
		{
			name:         "same-day typo is not turned into a 23-hour entry",
			start:        "15:00",
			end:          "14:00",
			wantStart:    todayAt(15, 0, 0),
			wantEnd:      todayAt(14, 0, 0),
			wantEndAfter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := resolveEntryTimes(mustParse(t, tt.start), mustParse(t, tt.end))

			if !gotStart.Equal(tt.wantStart) {
				t.Errorf("start = %s, want %s", gotStart, tt.wantStart)
			}
			if !gotEnd.Equal(tt.wantEnd) {
				t.Errorf("end = %s, want %s", gotEnd, tt.wantEnd)
			}
			if got := gotEnd.After(gotStart); got != tt.wantEndAfter {
				t.Errorf("end after start = %v, want %v", got, tt.wantEndAfter)
			}
			if tt.wantEndAfter && gotEnd.Sub(gotStart) != tt.wantDuration {
				t.Errorf("duration = %s, want %s", gotEnd.Sub(gotStart), tt.wantDuration)
			}
		})
	}
}

package cmd

import (
	"testing"
	"time"
)

func TestFormatStartTime(t *testing.T) {
	zurich, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		t.Fatalf("LoadLocation: %v", err)
	}

	tests := []struct {
		name  string
		start string
		date  string
		loc   *time.Location
		want  string
	}{
		{
			// The whole point of the fix: 16:00 in Zurich is 14:00 UTC in summer,
			// not 16:00 UTC. Sending it as UTC put every task two hours late.
			name:  "clock time is read in the given zone, sent as UTC",
			start: "16:00",
			date:  "2026-07-13",
			loc:   zurich,
			want:  "2026-07-13T14:00:00Z",
		},
		{
			name:  "same clock time in winter, when the offset is one hour",
			start: "16:00",
			date:  "2026-01-13",
			loc:   zurich,
			want:  "2026-01-13T15:00:00Z",
		},
		{
			name:  "UTC zone leaves the clock time alone",
			start: "16:00",
			date:  "2026-07-13",
			loc:   time.UTC,
			want:  "2026-07-13T16:00:00Z",
		},
		{
			name:  "a local time before the offset rolls back to the previous UTC day",
			start: "01:00",
			date:  "2026-07-13",
			loc:   zurich,
			want:  "2026-07-12T23:00:00Z",
		},
		{
			name:  "single digit hour",
			start: "9:05",
			date:  "2026-07-13",
			loc:   zurich,
			want:  "2026-07-13T07:05:00Z",
		},
		{
			// An ISO datetime already carries its own offset, so it is not ours to move.
			name:  "ISO datetime passes through untouched",
			start: "2026-07-13T16:00:00+02:00",
			date:  "2026-07-13",
			loc:   zurich,
			want:  "2026-07-13T16:00:00+02:00",
		},
		{
			name:  "empty start stays empty",
			start: "",
			date:  "2026-07-13",
			loc:   zurich,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatStartTime(tt.start, tt.date, tt.loc)
			if err != nil {
				t.Fatalf("formatStartTime(%q, %q) returned error: %v", tt.start, tt.date, err)
			}
			if got != tt.want {
				t.Errorf("formatStartTime(%q, %q) = %q, want %q", tt.start, tt.date, got, tt.want)
			}
		})
	}
}

func TestFormatStartTimeErrors(t *testing.T) {
	tests := []struct {
		name  string
		start string
		date  string
	}{
		{"clock time without a date", "16:00", ""},
		{"unparseable time", "half past four", "2026-07-13"},
		{"hour out of range", "25:00", "2026-07-13"},
		{"minute out of range", "16:75", "2026-07-13"},
		{"unparseable date", "16:00", "13.07.2026"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := formatStartTime(tt.start, tt.date, time.UTC); err == nil {
				t.Errorf("formatStartTime(%q, %q) = nil error, want an error", tt.start, tt.date)
			}
		})
	}
}

func TestResolveLocation(t *testing.T) {
	loc, err := resolveLocation("")
	if err != nil {
		t.Fatalf("resolveLocation(\"\") returned error: %v", err)
	}
	if loc != time.Local {
		t.Errorf("resolveLocation(\"\") = %v, want the local zone", loc)
	}

	loc, err = resolveLocation("Europe/Zurich")
	if err != nil {
		t.Fatalf("resolveLocation(\"Europe/Zurich\") returned error: %v", err)
	}
	if loc.String() != "Europe/Zurich" {
		t.Errorf("resolveLocation(\"Europe/Zurich\") = %v, want Europe/Zurich", loc)
	}

	if _, err := resolveLocation("Mars/Olympus_Mons"); err == nil {
		t.Error("resolveLocation of an unknown zone = nil error, want an error")
	}
}

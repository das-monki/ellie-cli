package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want time.Time
		ok   bool
	}{
		{
			// getTask and byDate return this shape. Reading it as a plain string
			// used to fail silently, so `tasks list` printed no times at all.
			name: "firestore timestamp object",
			raw:  `{"_seconds":1783951500,"_nanoseconds":0}`,
			want: time.Unix(1783951500, 0).UTC(),
			ok:   true,
		},
		{
			name: "firestore timestamp with nanoseconds",
			raw:  `{"_seconds":1783951500,"_nanoseconds":364000000}`,
			want: time.Unix(1783951500, 364000000).UTC(),
			ok:   true,
		},
		{
			// createTask and updateTask echo this shape back.
			name: "ISO datetime string",
			raw:  `"2026-07-13T14:00:00.000Z"`,
			want: time.Date(2026, 7, 13, 14, 0, 0, 0, time.UTC),
			ok:   true,
		},
		{
			name: "date-only string",
			raw:  `"2026-07-13"`,
			want: time.Date(2026, 7, 13, 0, 0, 0, 0, time.UTC),
			ok:   true,
		},
		{"null", `null`, time.Time{}, false},
		{"empty string", `""`, time.Time{}, false},
		{"absent", ``, time.Time{}, false},
		{"unrecognized object", `{"foo":1}`, time.Time{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseTimestamp(json.RawMessage(tt.raw))
			if ok != tt.ok {
				t.Fatalf("ParseTimestamp(%s) ok = %v, want %v", tt.raw, ok, tt.ok)
			}
			if ok && !got.Equal(tt.want) {
				t.Errorf("ParseTimestamp(%s) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestTaskGetStartAndDate(t *testing.T) {
	var task Task
	body := `{
		"id": "a5a9f7e9-17e0-479f-ae53-2dc716e8a9cc",
		"description": "Gym",
		"date": {"_seconds": 1783893600, "_nanoseconds": 0},
		"start": {"_seconds": 1783951500, "_nanoseconds": 0}
	}`
	if err := json.Unmarshal([]byte(body), &task); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	start, ok := task.GetStart()
	if !ok {
		t.Fatal("GetStart() reported no start time on a task that has one")
	}
	if want := time.Unix(1783951500, 0).UTC(); !start.Equal(want) {
		t.Errorf("GetStart() = %v, want %v", start, want)
	}

	date, ok := task.GetDate()
	if !ok {
		t.Fatal("GetDate() reported no date on a task that has one")
	}
	if want := time.Unix(1783893600, 0).UTC(); !date.Equal(want) {
		t.Errorf("GetDate() = %v, want %v", date, want)
	}

	var bare Task
	if err := json.Unmarshal([]byte(`{"id":"x","description":"y"}`), &bare); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if _, ok := bare.GetStart(); ok {
		t.Error("GetStart() reported a start time on a task with none")
	}
}

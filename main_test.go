package main

import (
	"testing"
	"time"
)

func TestParseTimeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "with empty string",
			input:    "",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with invalid string",
			input:    "not-a-time",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with pre-1970 date string",
			input:    "1969-12-31",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with pre-1970 rfc3339 format",
			input:    "1969-12-31T23:59:59Z",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with pre-1970 iso8601 with offset",
			input:    "1969-12-31T18:59:59-05:00",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with just before unix epoch",
			input:    "1969-12-31T23:59:59.999Z",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with negative unix timestamp (seconds)",
			input:    "-1687689000",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with negative unix timestamp (milliseconds)",
			input:    "-1687689000123",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "with unix epoch boundary",
			input:    "1970-01-01T00:00:00Z",
			wantTime: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with only date",
			input:    "2023-06-15",
			wantTime: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with timezone-less datetime",
			input:    "2023-06-15 10:30:00",
			wantTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with seconds-less datetime",
			input:    "2023-06-15 10:30",
			wantTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with rfc3339 format",
			input:    "2025-06-15T10:30:00.123Z",
			wantTime: time.Date(2025, 6, 15, 10, 30, 0, 123000000, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with iso8601 format including timezone",
			input:    "2025-06-15T10:30:00Z",
			wantTime: time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with iso8601 format including offset",
			input:    "2023-06-15T10:30:00-07:00",
			wantTime: time.Date(2023, 6, 15, 17, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "with unix timestamp (seconds)",
			input:    "1687689000",
			wantTime: time.Unix(1687689000, 0),
			wantErr:  false,
		},
		{
			name:     "with unix timestamp (milliseconds)",
			input:    "1687689000123",
			wantTime: time.UnixMilli(1687689000123),
			wantErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := parseTimeInput(test.input)

			if test.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !parsed.Equal(test.wantTime) {
				t.Errorf("got: %v, want: %v", parsed, test.wantTime)
			}
		})
	}
}

package checkin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Schedule bitmask: bit 0 = Sunday, bit 1 = Monday, ..., bit 6 = Saturday.
const (
	schedMon = uint8(1) << time.Monday
	schedTue = uint8(1) << time.Tuesday
	schedWed = uint8(1) << time.Wednesday
	schedThu = uint8(1) << time.Thursday
	schedFri = uint8(1) << time.Friday
	schedSat = uint8(1) << time.Saturday
	schedSun = uint8(1) << time.Sunday

	schedWeekdays = schedMon | schedTue | schedWed | schedThu | schedFri
	schedDaily    = schedWeekdays | schedSat | schedSun
)

// Dates anchored to a known week: Mon 2024-01-08 … Sun 2024-01-14.
func noon(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
}

func TestPrevScheduledDay(t *testing.T) {
	tests := []struct {
		name     string
		from     time.Time
		schedule uint8
		wantDay  int // day-of-month
	}{
		{
			name:     "daily: returns yesterday",
			from:     noon(2024, time.January, 10), // Wednesday
			schedule: schedDaily,
			wantDay:  9, // Tuesday
		},
		{
			name:     "weekdays-only: Monday from skips weekend to Friday",
			from:     noon(2024, time.January, 8), // Monday
			schedule: schedWeekdays,
			wantDay:  5, // previous Friday
		},
		{
			name:     "once-a-week Monday: Thursday from finds Monday",
			from:     noon(2024, time.January, 11), // Thursday
			schedule: schedMon,
			wantDay:  8, // Monday of same week
		},
		{
			name:     "once-a-week Monday: Tuesday from finds Monday same week",
			from:     noon(2024, time.January, 9), // Tuesday
			schedule: schedMon,
			wantDay:  8, // Monday of same week
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := prevScheduledDay(tc.from, tc.schedule)
			assert.Equal(t, tc.wantDay, got.Day(), "unexpected day")
		})
	}
}

func TestSameDay(t *testing.T) {
	tests := []struct {
		name string
		a, b time.Time
		want bool
	}{
		{
			name: "identical times",
			a:    time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			b:    time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "same date, different times",
			a:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC),
			want: true,
		},
		{
			name: "adjacent days",
			a:    time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC),
			b:    time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "month boundary",
			a:    time.Date(2024, 1, 31, 12, 0, 0, 0, time.UTC),
			b:    time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "same month-day, different year",
			a:    time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC),
			b:    time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, sameDay(tc.a, tc.b))
		})
	}
}

package checkin_test

import (
	"context"
	"testing"
	"time"

	"github.com/Uranury/HabitTracker/internal/checkin"
	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/Uranury/HabitTracker/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recentOccurrences returns the n most-recent occurrences of weekday in descending
// order (index 0 = most recent), each at midnight UTC. If today is weekday, it is
// included as index 0.
func recentOccurrences(n int, weekday time.Weekday) []time.Time {
	d := time.Now().UTC()
	for d.Weekday() != weekday {
		d = d.AddDate(0, 0, -1)
	}
	days := make([]time.Time, n)
	for i := range days {
		days[i] = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		d = d.AddDate(0, 0, -7)
	}
	return days
}

func schedFor(weekday time.Weekday) uint8 { return uint8(1) << weekday }

func seedStreakHabit(t *testing.T, repo habit.Repository, userID uuid.UUID, schedule uint8) *habit.Habit {
	t.Helper()
	h := &habit.Habit{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "streak test habit",
		Schedule:  schedule,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	require.NoError(t, repo.Create(context.Background(), h))
	return h
}

func seedCheckin(t *testing.T, repo checkin.Repository, userID, habitID uuid.UUID, date time.Time) {
	t.Helper()
	now := time.Now().UTC()
	c := &checkin.CheckIn{
		ID:        uuid.New(),
		UserID:    userID,
		HabitID:   habitID,
		Status:    checkin.Checked,
		Date:      date,
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, repo.Record(context.Background(), c))
}

func newStreakSvc(t *testing.T) (checkin.Repository, habit.Repository, *checkin.Service, uuid.UUID) {
	t.Helper()
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	habitRepo := habit.NewRepository(db)
	checkinRepo := checkin.NewRepository(db)
	svc := checkin.NewService(checkinRepo, habitRepo)
	return checkinRepo, habitRepo, svc, userID
}

// ── Monday schedule ───────────────────────────────────────────────────────────

func TestGetCurrentStreak_Monday_Consecutive(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Monday))
	for _, d := range recentOccurrences(3, time.Monday) {
		seedCheckin(t, checkinRepo, userID, h.ID, d)
	}

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 3, streak)
}

func TestGetCurrentStreak_Monday_Gap(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Monday))
	days := recentOccurrences(3, time.Monday)
	seedCheckin(t, checkinRepo, userID, h.ID, days[0]) // w0 — checked
	// days[1] intentionally skipped
	seedCheckin(t, checkinRepo, userID, h.ID, days[2]) // w2 — isolated

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 1, streak)
}

func TestGetCurrentStreak_Monday_MostRecentMissed(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Monday))
	days := recentOccurrences(3, time.Monday)
	// skip days[0] (most recent)
	seedCheckin(t, checkinRepo, userID, h.ID, days[1])
	seedCheckin(t, checkinRepo, userID, h.ID, days[2])

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 0, streak)
}

// ── Tuesday schedule ──────────────────────────────────────────────────────────

func TestGetCurrentStreak_Tuesday_Consecutive(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Tuesday))
	for _, d := range recentOccurrences(3, time.Tuesday) {
		seedCheckin(t, checkinRepo, userID, h.ID, d)
	}

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 3, streak)
}

func TestGetCurrentStreak_Tuesday_Gap(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Tuesday))
	days := recentOccurrences(3, time.Tuesday)
	seedCheckin(t, checkinRepo, userID, h.ID, days[0])
	// days[1] skipped
	seedCheckin(t, checkinRepo, userID, h.ID, days[2])

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 1, streak)
}

func TestGetCurrentStreak_Tuesday_MostRecentMissed(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Tuesday))
	days := recentOccurrences(3, time.Tuesday)
	seedCheckin(t, checkinRepo, userID, h.ID, days[1])
	seedCheckin(t, checkinRepo, userID, h.ID, days[2])

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 0, streak)
}

// ── Long run ──────────────────────────────────────────────────────────────────

// TestGetCurrentStreak_LongRun confirms that 14 consecutive Monday check-ins
// produce a streak of 14. This exercises the full loop across ~3 months of history
// and ensures no off-by-one in the prevScheduledDay comparison chain.
func TestGetCurrentStreak_LongRun(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Monday))
	for _, d := range recentOccurrences(14, time.Monday) {
		seedCheckin(t, checkinRepo, userID, h.ID, d)
	}

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 14, streak)
}

// TestGetCurrentStreak_LongRunWithGapInMiddle seeds 14 Mondays but drops one in
// the middle (w7). The streak should equal the run of consecutive weeks before
// the gap, counting from the most recent check-in.
func TestGetCurrentStreak_LongRunWithGapInMiddle(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, schedFor(time.Monday))
	days := recentOccurrences(14, time.Monday)
	for i, d := range days {
		if i == 7 {
			continue // gap at w7
		}
		seedCheckin(t, checkinRepo, userID, h.ID, d)
	}

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 7, streak) // w0–w6 consecutive; w7 gap stops the count
}

// ── Daily schedule ────────────────────────────────────────────────────────────

// recentDays returns the n most-recent calendar days in descending order (index 0
// = today), each at midnight UTC.
func recentDays(n int) []time.Time {
	today := time.Now().UTC()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	days := make([]time.Time, n)
	for i := range days {
		days[i] = today.AddDate(0, 0, -i)
	}
	return days
}

func dailySchedule() uint8 {
	var s uint8
	for d := time.Sunday; d <= time.Saturday; d++ {
		s |= schedFor(d)
	}
	return s
}

// TestGetCurrentStreak_Daily_GapBreaksStreak checks that a gap inside a run of
// daily check-ins stops the streak at the contiguous block closest to today.
// Seeded: today, yesterday, (gap), 3 days ago, 4 days ago → streak 2.
func TestGetCurrentStreak_Daily_GapBreaksStreak(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, dailySchedule())
	days := recentDays(5)
	for i, d := range days {
		if i == 2 {
			continue // gap 2 days ago
		}
		seedCheckin(t, checkinRepo, userID, h.ID, d)
	}

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 2, streak)
}

// TestGetCurrentStreak_Daily_MissedYesterday verifies that a gap on yesterday
// limits the streak to today only even when older check-ins exist.
func TestGetCurrentStreak_Daily_MissedYesterday(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, dailySchedule())
	days := recentDays(4)
	seedCheckin(t, checkinRepo, userID, h.ID, days[0]) // today
	// days[1] = yesterday — intentionally skipped
	seedCheckin(t, checkinRepo, userID, h.ID, days[2])
	seedCheckin(t, checkinRepo, userID, h.ID, days[3])

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 1, streak)
}

// ── Mon–Fri schedule ──────────────────────────────────────────────────────────

// lastNWeekdays returns the n most-recent Mon–Fri days in descending order.
func lastNWeekdays(n int) []time.Time {
	days := make([]time.Time, 0, n)
	d := time.Now().UTC()
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	for len(days) < n {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			days = append(days, d)
		}
		d = d.AddDate(0, 0, -1)
	}
	return days
}

// TestGetCurrentStreak_Daily_TwoSeparateGaps seeds three isolated islands of
// check-ins separated by two distinct gaps. Only the run closest to today counts;
// the loop must stop at the first gap and never reach the second island.
func TestGetCurrentStreak_Daily_TwoSeparateGaps(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, dailySchedule())
	days := recentDays(8)
	// island 1 — days[0]+days[1]
	seedCheckin(t, checkinRepo, userID, h.ID, days[0])
	seedCheckin(t, checkinRepo, userID, h.ID, days[1])
	// gap at days[2]
	// island 2 — days[3]+days[4]
	seedCheckin(t, checkinRepo, userID, h.ID, days[3])
	seedCheckin(t, checkinRepo, userID, h.ID, days[4])
	// gap at days[5]
	// island 3 — days[6]+days[7]
	seedCheckin(t, checkinRepo, userID, h.ID, days[6])
	seedCheckin(t, checkinRepo, userID, h.ID, days[7])

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 2, streak)
}

// TestGetCurrentStreak_GapAtOldestCheckin confirms that when the loop reaches the
// oldest available check-in and it breaks the consecutive run, the accumulated
// streak is returned rather than zero. Running out of rows and hitting a gap are
// distinct outcomes: the former keeps the count, the latter stops it.
func TestGetCurrentStreak_GapAtOldestCheckin(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	h := seedStreakHabit(t, habitRepo, userID, dailySchedule())
	days := recentDays(5)
	seedCheckin(t, checkinRepo, userID, h.ID, days[0])
	seedCheckin(t, checkinRepo, userID, h.ID, days[1])
	// gap at days[2] and days[3]
	seedCheckin(t, checkinRepo, userID, h.ID, days[4]) // oldest row, nothing before it

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 2, streak)
}

// TestGetCurrentStreak_Weekdays_MidweekGap seeds the 5 most-recent Mon–Fri days
// but skips the 3rd (index 2). The 2 days before the gap form the only unbroken
// run from today → streak 2.
func TestGetCurrentStreak_Weekdays_MidweekGap(t *testing.T) {
	checkinRepo, habitRepo, svc, userID := newStreakSvc(t)
	weekdays := schedFor(time.Monday) | schedFor(time.Tuesday) | schedFor(time.Wednesday) |
		schedFor(time.Thursday) | schedFor(time.Friday)
	h := seedStreakHabit(t, habitRepo, userID, weekdays)
	days := lastNWeekdays(5)
	for i, d := range days {
		if i == 2 {
			continue // gap in the middle
		}
		seedCheckin(t, checkinRepo, userID, h.ID, d)
	}

	streak, err := svc.GetCurrentStreak(context.Background(), userID, h.ID, "UTC")
	require.NoError(t, err)
	assert.Equal(t, 2, streak)
}

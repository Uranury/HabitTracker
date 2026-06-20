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

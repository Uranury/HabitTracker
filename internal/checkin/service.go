package checkin

import (
	"context"
	"time"

	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/google/uuid"
)

// NOTE: check-ins are stored as midnight UTC of the user's local date so streak calculation is timezone-correct
// NOTE: one check-in per day is enforced via unique constraint on (user_id, habit_id, date)
// NOTE: no background worker/cron job is needed to fill missed gaps, prevScheduledDay handles this by itself

type Service struct {
	repo       Repository
	habitsRepo habit.Repository
}

func NewService(repo Repository, habitsRepo habit.Repository) *Service {
	return &Service{repo: repo, habitsRepo: habitsRepo}
}

func (svc *Service) CheckIn(ctx context.Context, userID, habitID uuid.UUID, timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	localNow := time.Now().In(loc)
	localDay := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, time.UTC)
	now := time.Now().UTC()
	c := &CheckIn{ID: uuid.New(), UserID: userID, HabitID: habitID, Status: Checked, Date: localDay, CreatedAt: now, UpdatedAt: now}
	return svc.repo.Record(ctx, c)
}

func (svc *Service) GetCheckins(ctx context.Context, userID, habitID uuid.UUID, limit, offset int) ([]*CheckIn, error) {
	_, err := svc.habitsRepo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		return nil, err
	}
	checkins, err := svc.repo.GetByUserAndHabitID(ctx, userID, habitID, limit, offset)
	if err != nil {
		return nil, err
	}
	return checkins, nil
}

func (svc *Service) DeleteCheckin(ctx context.Context, userID, habitID, checkinID uuid.UUID) error {
	return svc.repo.Delete(ctx, userID, habitID, checkinID)
}

func (svc *Service) GetCurrentStreak(ctx context.Context, userID, habitID uuid.UUID, timezone string) (int, error) {
	hbt, err := svc.habitsRepo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		return 0, err
	}
	schedule := hbt.Schedule
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, err
	}
	checkIns, err := svc.repo.GetByUserAndHabitID(ctx, userID, habitID, -1, 0)
	if err != nil {
		return 0, err
	}
	localNow := time.Now().In(loc)
	today := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, time.UTC)
	lastScheduled := prevScheduledDay(today.AddDate(0, 0, 1), schedule)
	if len(checkIns) == 0 || !sameDay(checkIns[0].Date, lastScheduled) {
		return 0, nil
	}

	streak := 0
	for i, checkIn := range checkIns {
		if i > 0 {
			prev := checkIns[i-1]
			expected := prevScheduledDay(prev.Date, schedule)
			if !sameDay(checkIn.Date, expected) {
				break
			}
		}
		streak++
	}
	return streak, nil
}

func prevScheduledDay(from time.Time, schedule uint8) time.Time {
	d := from.AddDate(0, 0, -1)
	for i := 0; i < 7; i++ {
		mask := uint8(1) << uint8(d.Weekday())
		if schedule&mask != 0 {
			return d
		}
		d = d.AddDate(0, 0, -1)
	}
	return time.Time{}
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

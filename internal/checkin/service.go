package checkin

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/google/uuid"
	"time"
)

// TODO: account for user timezones
// TODO: enforce one check in per day
// TODO: no background worker/cron job is needed to fill missed gaps, prevScheduleDay handles this by itself

type Service struct {
	repo       Repository
	habitsRepo habit.Repository
}

func NewService(repo Repository, habitsRepo habit.Repository) *Service {
	return &Service{repo: repo, habitsRepo: habitsRepo}
}

func (svc *Service) GetCurrentStreak(ctx context.Context, userID, habitID uuid.UUID) (int, error) {
	hbt, err := svc.habitsRepo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		return 0, err
	}
	schedule := hbt.Schedule
	checkIns, err := svc.repo.GetByUserAndHabitID(ctx, userID, habitID)
	if err != nil {
		return 0, err
	}
	streak := 0
	for i, checkIn := range checkIns {
		weekday := checkIn.Date.Weekday()
		weekdayMask := uint(1) << weekday
		if weekdayMask&schedule == 0 || checkIn.Status != Checked {
			break
		}
		if i > 0 {
			prev := checkIns[i-1]
			expected := prevScheduledDay(checkIn.Date, schedule)
			if !sameDay(prev.Date, expected) {
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

package habit_test

import (
	"context"
	"testing"
	"time"

	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/Uranury/HabitTracker/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHabitRepository_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := &habit.Habit{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "Morning run",
		Schedule:  0b0111110, // Mon–Fri
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.Create(context.Background(), h)
	require.NoError(t, err)
}

func TestHabitRepository_GetHabitByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	desc := "run before breakfast"
	h := &habit.Habit{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "Morning run",
		Schedule:    0b0111110,
		Description: &desc,
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		UpdatedAt:   time.Now().UTC().Truncate(time.Second),
	}
	require.NoError(t, repo.Create(context.Background(), h))

	got, err := repo.GetHabitByID(context.Background(), userID, h.ID)
	require.NoError(t, err)

	assert.Equal(t, h.ID, got.ID)
	assert.Equal(t, h.Name, got.Name)
	assert.Equal(t, h.Schedule, got.Schedule)
	assert.Equal(t, h.Description, got.Description)
}

func TestHabitRepository_GetHabitByID_WrongUser(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	otherUserID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := &habit.Habit{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "Morning run",
		Schedule:  0b0111110,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	require.NoError(t, repo.Create(context.Background(), h))

	_, err := repo.GetHabitByID(context.Background(), otherUserID, h.ID)
	assert.Error(t, err, "should not return another user's habit")
}

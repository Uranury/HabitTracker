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

// seedHabit inserts a habit and returns it. Use it to set up preconditions.
func seedHabit(t *testing.T, repo habit.Repository, userID uuid.UUID, name string) *habit.Habit {
	t.Helper()
	h := &habit.Habit{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Schedule:  0b0111110, // Mon–Fri
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UpdatedAt: time.Now().UTC().Truncate(time.Second),
	}
	require.NoError(t, repo.Create(context.Background(), h))
	return h
}

// --- Create ---

func TestHabitRepository_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
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
}

func TestHabitRepository_Create_WithDescription(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	desc := "before breakfast"
	h := &habit.Habit{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "Morning run",
		Schedule:    0b0111110,
		Description: &desc,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	require.NoError(t, repo.Create(context.Background(), h))
}

// --- GetHabitByID ---

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

func TestHabitRepository_GetHabitByID_NullDescription(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := seedHabit(t, repo, userID, "Morning run")

	got, err := repo.GetHabitByID(context.Background(), userID, h.ID)
	require.NoError(t, err)
	assert.Nil(t, got.Description)
}

func TestHabitRepository_GetHabitByID_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	_, err := repo.GetHabitByID(context.Background(), userID, uuid.New())
	assert.ErrorIs(t, err, habit.ErrHabitNotFound)
}

func TestHabitRepository_GetHabitByID_WrongUser(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	otherUserID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := seedHabit(t, repo, userID, "Morning run")

	_, err := repo.GetHabitByID(context.Background(), otherUserID, h.ID)
	assert.ErrorIs(t, err, habit.ErrHabitNotFound)
}

// --- GetHabitsByUserID ---

func TestHabitRepository_GetHabitsByUserID(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	seedHabit(t, repo, userID, "Habit A")
	seedHabit(t, repo, userID, "Habit B")

	habits, err := repo.GetHabitsByUserID(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, habits, 2)

	names := []string{habits[0].Name, habits[1].Name}
	assert.ElementsMatch(t, []string{"Habit A", "Habit B"}, names)
}

func TestHabitRepository_GetHabitsByUserID_Empty(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	habits, err := repo.GetHabitsByUserID(context.Background(), userID)
	require.NoError(t, err)
	assert.Empty(t, habits)
}

func TestHabitRepository_GetHabitsByUserID_Isolation(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	otherUserID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	seedHabit(t, repo, userID, "My habit")
	seedHabit(t, repo, otherUserID, "Other user's habit")

	habits, err := repo.GetHabitsByUserID(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, habits, 1)
	assert.Equal(t, "My habit", habits[0].Name)
}

// --- Update ---

func TestHabitRepository_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := seedHabit(t, repo, userID, "Morning run")

	newDesc := "updated description"
	updated := &habit.Habit{
		ID:          h.ID,
		UserID:      userID,
		Name:        "Evening walk",
		Schedule:    0b1000001, // Sun + Sat
		Description: &newDesc,
	}
	require.NoError(t, repo.Update(context.Background(), updated))

	got, err := repo.GetHabitByID(context.Background(), userID, h.ID)
	require.NoError(t, err)
	assert.Equal(t, "Evening walk", got.Name)
	assert.Equal(t, uint8(0b1000001), got.Schedule)
	assert.Equal(t, &newDesc, got.Description)
}

func TestHabitRepository_Update_ClearsDescription(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	desc := "original description"
	h := &habit.Habit{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "Morning run",
		Schedule:    0b0111110,
		Description: &desc,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	require.NoError(t, repo.Create(context.Background(), h))

	updated := &habit.Habit{
		ID:          h.ID,
		UserID:      userID,
		Name:        "Morning run",
		Schedule:    0b0111110,
		Description: nil,
	}
	require.NoError(t, repo.Update(context.Background(), updated))

	got, err := repo.GetHabitByID(context.Background(), userID, h.ID)
	require.NoError(t, err)
	assert.Nil(t, got.Description)
}

func TestHabitRepository_Update_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	err := repo.Update(context.Background(), &habit.Habit{
		ID:     uuid.New(),
		UserID: userID,
		Name:   "Ghost",
	})
	assert.ErrorIs(t, err, habit.ErrHabitNotFound)
}

// --- Delete ---

func TestHabitRepository_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := seedHabit(t, repo, userID, "Morning run")

	require.NoError(t, repo.Delete(context.Background(), userID, h.ID))

	_, err := repo.GetHabitByID(context.Background(), userID, h.ID)
	assert.ErrorIs(t, err, habit.ErrHabitNotFound)
}

func TestHabitRepository_Delete_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	err := repo.Delete(context.Background(), userID, uuid.New())
	assert.ErrorIs(t, err, habit.ErrHabitNotFound)
}

func TestHabitRepository_Delete_WrongUser(t *testing.T) {
	db := testutil.NewTestDB(t)
	userID := testutil.SeedUser(t, db)
	otherUserID := testutil.SeedUser(t, db)
	repo := habit.NewRepository(db)

	h := seedHabit(t, repo, userID, "Morning run")

	err := repo.Delete(context.Background(), otherUserID, h.ID)
	assert.ErrorIs(t, err, habit.ErrHabitNotFound)

	// Verify the habit still exists for the real owner.
	_, err = repo.GetHabitByID(context.Background(), userID, h.ID)
	require.NoError(t, err)
}

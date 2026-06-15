CREATE TABLE IF NOT EXISTS checkins (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id TEXT NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('missed', 'skipped', 'checked')),
    date INTEGER NOT NULL CHECK (date >= 0 AND date <= 255),

    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
# Milestones

## Bugs
- When checking in on an unscheduled day the streak breaks.
---

## M1 — Wire the existing code

- [+] Apply JWT middleware to all protected routes
- [+] Register habit routes (POST, GET list, GET by ID, PATCH, DELETE)
- [+] Register check-in routes (`POST /habits/:id/checkin`, `GET /habits/:id/checkins`, `GET /habits/:id/streak`)
- [+] Fix `GetUserTimeZone()` context key bug
- [+] Enforce one check-in per day (UNIQUE constraint in migration 3)

---

## M2 — Habits CRUD

- [+] `GET /api/habits` — list all habits for the authenticated user
- [+] `GET /api/habits/:id` — get a single habit by ID
- [+] `PATCH /api/habits/:id` — partial update (name, schedule, description); read-then-merge in service
- [+] `DELETE /api/habits/:id` — hard delete; returns 404 if not found

---

## M3 — Check-ins

- [+] Wire check-in service and repository in `infra.go`
- [+] Create `checkin/handler.go`
- [+] `POST /api/habits/:id/checkin` — record today's check-in (Checked); return a proper error on duplicate, not a silent DB conflict
- [+] `GET /api/habits/:id/checkins` — paginated check-in history (`limit`/`offset` query params, default limit 20)
- [+] `GET /api/habits/:id/streak` — current streak
- [ ] `GET /api/habits/:id/streak/best` — best streak (no service method yet; requires full history scan)

---

## M4 — User profile

Fully implemented, expanded beyond the original plan. Routes are under `/api/users/me`.

- [+] `GET /api/users/me` — return id, username, timezone, avatar
- [+] `PATCH /api/users/me/avatar` — update avatar URL
- [+] `PATCH /api/users/me/timezone` — update timezone
- [+] `PATCH /api/users/me/username` — update username
- [+] `PATCH /api/users/me/password` — change password with old-password verification

---

## M5 — API quality

- [ ] Return the created habit in `POST /api/habits` (currently 201 with no body) — client needs a reload to get the ID
- [ ] Return the updated habit in `PATCH /api/habits/:id` (currently 200 with no body) — client needs a reload to reflect changes
- [ ] `POST /api/habits/:id/checkin` — return `{"streak": N}` so the client doesn't need a second request to refresh the streak card
- [ ] Consistent 404 vs 500 — `ErrHabitNotFound` should produce 404, not 500, in all habit and check-in handlers
- [ ] Validate `schedule != 0` on habit create/update — reject at handler level; currently the frontend must guard this itself
- [ ] Standardise error response shape — `{"error": "..."}` is used everywhere but some success responses use `gin.H{}` inconsistently
- [ ] `GET /api/habits/today` — list habits with each habit's today check-in status and current streak in one call; eliminates the N+1 streak requests the frontend currently fires on sidebar load

---

## M6 — Observability & reliability

- [ ] `GET /healthz` — liveness endpoint, no auth, returns uptime
- [ ] Request ID middleware — generate a UUID per request, log it, return as `X-Request-ID` header
- [ ] Graceful shutdown — `server.Shutdown(ctx)` exists but is never called; wire it to `os.Signal` in `main.go`
- [ ] Config validation at startup — fail fast if `JWT_SECRET` is shorter than 32 bytes
- [ ] Startup log — print listen address, DB path, and Go version on start

---

## M7 — Testing gaps

- [ ] Auth flow integration test — signup → login → call a protected route via `httptest`
- [ ] Check-in repository tests — cover `Record`, `GetByUserAndHabitID`, one-per-day UNIQUE constraint
- [ ] Check-in streak tests with real DB — `GetCurrentStreak` through service→repo stack, not just pure functions
- [ ] User repository tests — `Create`, `FindByUsername`, targeted update methods
- [ ] Habit handler tests — HTTP layer via `httptest`, covering auth header, 404 vs 500

---

## Nice to have

- **Best streak across all habits** — `GET /api/users/me/stats` with per-habit personal bests
- **Soft delete for habits** — `archived_at` column instead of hard delete; preserves check-in history
- **Token refresh** — 24-hour tokens with no refresh; add `POST /auth/refresh` or a refresh-token pair
- **PostgreSQL support** — DSN switching exists; needs `?` → `$1` substitution and a pgx/lib-pq driver swap
- **Docker / single-binary packaging** — `CGO_ENABLED=0 go build -o habittracker ./cmd/api` already works; a `FROM scratch` Dockerfile finishes the fast-launch goal
- **OpenAPI spec** — worth generating once M3 is wired

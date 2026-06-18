# Milestones

## Bugs

- **`GET /api/habits/:id` param name mismatch** — route registers `:id` but `GetHabit` handler calls `c.Param("habitID")`; every request gets an empty UUID and hits a DB error. Fix: change to `c.Param("id")`.
- **`DELETE /api/habits/:id` route not registered** — `DeleteHabit` handler is fully implemented and tested; the route is simply missing from `server.go`.
- **`UpdateAvatar` SQL args swapped** — `user/repository.go:76` calls `ExecContext(ctx, query, userID, avatar)` but the query is `SET avatar=? WHERE id=?`; args should be `(avatar, userID)`. Every avatar update currently returns 404.
- **Habit partial update zeroes untouched fields** — `Update` uses `COALESCE(?, col)` but Go zero values (`""`, `0`) are not SQL NULL; updating only description will silently reset `name` to `""` and `schedule` to 0. Fix: pass `sql.NullString` / `sql.NullInt64`, or build a dynamic SET clause.
- **Check-in date not timezone-aware** — `checkin.Service.CheckIn` records `date` as `time.Now()` (UTC wall time). A user in UTC+9 checking in at 00:30 local time gets the check-in filed on the previous UTC day, breaking streak calculation which is timezone-aware.

---

## M1 — Wire the existing code

- [+] Apply JWT middleware to all protected routes
- [+] Register habit routes (POST, GET list, GET by ID, PATCH) — DELETE still missing, see Bugs
- [ ] Register check-in routes (`POST /habits/:id/checkin`, `GET /habits/:id/checkins`, `GET /habits/:id/streak`)
- [+] Fix `GetUserTimeZone()` context key bug
- [+] Enforce one check-in per day (UNIQUE constraint in migration 3)

---

## M2 — Habits CRUD

- [+] `GET /api/habits` — list all habits for the authenticated user
- [+] `GET /api/habits/:id` — handler done; has param name bug (see Bugs)
- [+] `PATCH /api/habits/:id` — update name, schedule, or description; has partial-update bug (see Bugs)
- [ ] `DELETE /api/habits/:id` — handler + service + repo done and tested; route not registered (see Bugs)

---

## M3 — Check-ins

Service and repository exist but nothing is wired. `checkin.Service` and `checkin.Repository` are not instantiated in `infra.go`. No handler file exists.

- [ ] Wire check-in service and repository in `infra.go`
- [ ] Create `checkin/handler.go`
- [ ] `POST /api/habits/:id/checkin` — record today's check-in (Checked); return a proper error on duplicate, not a silent DB conflict
- [ ] `POST /api/habits/:id/skip` — record today's skip (Skipped status)
- [ ] `GET /api/habits/:id/checkins` — paginated check-in history
- [ ] `GET /api/habits/:id/streak` — current streak (`GetCurrentStreak` logic exists in service)
- [ ] `GET /api/habits/:id/streak/best` — best streak (no service method yet; requires full history scan)

---

## M4 — User profile

Fully implemented, expanded beyond the original plan. Routes are under `/api/users/me`.

- [+] `GET /api/users/me` — return id, username, timezone, avatar
- [+] `PATCH /api/users/me/avatar` — update avatar URL; has swapped SQL args bug (see Bugs)
- [+] `PATCH /api/users/me/timezone` — update timezone
- [+] `PATCH /api/users/me/username` — update username
- [+] `PATCH /api/users/me/password` — change password with old-password verification

---

## M5 — API quality

- [ ] Fix the three bugs above (param mismatch, missing DELETE route, avatar args)
- [ ] Fix partial-update bug in `PATCH /api/habits/:id`
- [ ] Return the created resource in `POST /api/habits` (currently 201 with no body)
- [ ] Return the updated resource in `PATCH /api/habits/:id` (currently 200 with no body)
- [ ] Consistent 404 vs 500 — `ErrHabitNotFound` should produce 404, not 500, in habit handlers
- [ ] Validate `schedule != 0` on habit create/update — reject at handler level before hitting the DB
- [ ] Standardise error response shape — `{"error": "..."}` is used everywhere but some handlers use `gin.H{}` for success responses inconsistently

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

- **`GET /api/habits/today`** — all habits + today's check-in status in one call; avoids N+1 on a dashboard
- **Best streak across all habits** — `GET /api/users/me/stats` with per-habit personal bests
- **Soft delete for habits** — `archived_at` column instead of hard delete; preserves check-in history
- **Token refresh** — 24-hour tokens with no refresh; add `POST /auth/refresh` or a refresh-token pair
- **PostgreSQL support** — DSN switching exists; needs `?` → `$1` substitution and a pgx/lib-pq driver swap
- **Docker / single-binary packaging** — `CGO_ENABLED=0 go build -o habittracker ./cmd/api` already works; a `FROM scratch` Dockerfile finishes the fast-launch goal
- **OpenAPI spec** — worth generating once M3 is wired

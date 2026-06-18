# Milestones

## Bugs

---

## M1 — Wire the existing code

Everything in this milestone is already implemented but not connected.

- [+] Apply JWT middleware to all protected routes (`.Use(m.JWTAuth())` in `server.go`)
- [+] Register habit routes (`POST /habits`, `GET /habits`, `GET /habits/:id`, `PUT /habits/:id`, `DELETE /habits/:id`)
- [ ] Register check-in routes (`POST /habits/:id/checkin`, `GET /habits/:id/checkins`, `GET /habits/:id/streak`)
- [+] Fix `GetUserTimeZone()` bug (see Bugs above)
- [+] Enforce one check-in per day in `checkin.Service.CheckIn()`

---

## M2 — Habits CRUD

- [+] `GET /habits` — list all habits for the authenticated user
- [+] `GET /habits/:id` — get a single habit
- [+] `PUT /habits/:id` — update name, schedule, or description
- [+] `DELETE /habits/:id` — delete a habit and its check-ins (cascade handled by FK trigger in migration 4)

---

## M3 — Check-ins

- [ ] `POST /habits/:id/checkin` — record a check-in (Checked status)
- [ ] `POST /habits/:id/skip` — record a skip (Skipped status) for today
- [ ] `GET /habits/:id/checkins` — paginated check-in history
- [ ] `GET /habits/:id/streak` — current streak (logic exists in `checkin.Service.GetCurrentStreak`)
- [ ] `GET /habits/:id/streak/best` — best streak ever (requires scanning full history)

---

## M4 — User profile

- [+] `GET /me` — return current user's profile (id, username, timezone, avatar)
- [+] `PUT /me` — update timezone (and optionally username)
- [+] `POST /me/avatar` — upload avatar image (`UserService.UploadAvatar` exists, no endpoint yet)

---

## Nice to have

- **Stats endpoint** — `GET /habits/:id/stats` returning weekly/monthly completion rate
- **Bulk check-in status** — `GET /habits/today` returning all habits with today's check-in status in one call, to support a dashboard view
- **PostgreSQL support** — DSN switching already exists in `pkg/database/db.go`; would need to replace `?` placeholders with `$1, $2, ...` in queries
- **Tests** — no test files exist; at minimum, unit tests for streak logic (`prevScheduledDay`, `sameDay`) and integration tests for auth flow
- **OpenAPI/Swagger spec** — useful once the API stabilises
- **Docker / single-binary packaging** — fits the original goal of a fast-launch executable; `CGO_ENABLED=0 go build -o habittracker ./cmd/api` already produces a static binary with `modernc.org/sqlite` (pure Go driver, no cgo needed)

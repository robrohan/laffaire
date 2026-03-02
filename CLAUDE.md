# Laffaire — Claude Reference

## Project Overview

A Go web app that lets users create **Events** (calendar groups/categories) and **Entries**
(individual calendar items), then subscribe to them via iCal. Auth is OAuth2 (Google or
compatible). Storage is SQLite (default) or PostgreSQL.

**Naming convention**: "Event" = a named calendar/category container. "Entry" = an individual
calendar event within that container.

---

## Tech Stack

- **Language**: Go
- **Router**: `github.com/gorilla/mux`
- **DB**: SQLite3 (`github.com/mattn/go-sqlite3`) or PostgreSQL (`github.com/lib/pq`)
- **DB access**: `github.com/jmoiron/sqlx` with prepared statements
- **Migrations**: `github.com/rubenv/sql-migrate`
- **Config**: `github.com/ardanlabs/conf` — env vars prefixed `WB_`
- **Templates**: Go `html/template`
- **Auth**: OAuth2 cookie (`WB_AT` = `{uuid}:{md5hash}`)

---

## Key File Locations

| Path | Purpose |
|---|---|
| `cmd/server/main.go` | Entry point, router setup, middleware, auth handlers |
| `internals/env/env.go` | `Env` struct — shared deps injected into handlers |
| `internals/handlers/handlers.go` | HTML page handlers |
| `internals/handlers/api.go` | **JSON API handlers** (new file) |
| `internals/models/models.go` | `User`, `Event`, `Entry` structs |
| `internals/repository/repo.go` | Data access — all DB queries |
| `internals/repository/conn.go` | DB connection and migration runner |
| `internals/ical/ical.go` | iCalendar format generation |
| `migrations/` | SQL schema files |
| `templates/` | Go HTML templates |

---

## Patterns

### Handler pattern
All handlers are higher-order functions returning `http.HandlerFunc`. HTML handlers take
`*template.Template`; API handlers do not.

```go
// HTML handler
func EventsPage(env *env.Env, t *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) { ... }
}

// API handler
func EventsAPI(env *env.Env) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) { ... }
}
```

### JSON response helpers (api.go)
Use `writeJSON` and `writeError` helpers defined in `api.go` for all API responses.

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}
```

### Repository access
Always go through `env.Repo` — never query `env.Db` directly in handlers.

```go
userId, _ := uuid.Parse(env.User.UUID)
events, err := env.Repo.GetEventsByUserId(userId)
```

### Error handling in handlers
Log with `env.Log`, return an HTTP error, and return early. Never panic.

```go
if err != nil {
    env.Log.Error("description", "error", err)
    writeError(w, http.StatusInternalServerError, "description")
    return
}
```

### Routing
- Public routes: registered directly on `router`
- Authenticated routes: registered on `secure` subrouter (`/-/` prefix for HTML, `/api/v1/` for JSON API)
- Secure subrouter uses `LoginVerify` middleware which sets `env.User`

```go
secure.HandleFunc("/events", handlers.EventsPage(env, templates)).Methods("GET")
apiSecure.HandleFunc("/events", handlers.EventsAPI(env)).Methods("GET")
```

### Model struct tags
Models use `db:` tags for sqlx. JSON API requires `json:` tags too.

```go
type Event struct {
    UUID        string `db:"uuid"        json:"id"`
    UserId      string `db:"user_uuid"   json:"-"`
    Title       string `db:"title"       json:"title"`
    Description string `db:"description" json:"description"`
}
```

Note: `UserId` is internal — omit it from API responses with `json:"-"`.

### Request body decoding (API POST/PUT)
```go
var input SomeStruct
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    writeError(w, http.StatusBadRequest, "invalid request body")
    return
}
```

### Auth middleware behaviour
`LoginVerify` sets `env.User` on each request. API routes reuse the same middleware.
On failure it currently redirects to `/login` — for API routes it should return 401 JSON.
A separate `APILoginVerify` middleware handles this cleanly.

---

## Testing

Test file: `internals/handlers/api_test.go` (package `handlers_test`)
Run all tests: `go test ./...`

### Test DB setup
Use an in-memory SQLite with `migrate.MemoryMigrationSource` — avoids needing the
`migrations/` directory on disk when running tests from a sub-package.

```go
func newTestDB(t *testing.T) *sqlx.DB {
    db, _ := sqlx.Open("sqlite3", ":memory:")
    src := &migrate.MemoryMigrationSource{
        Migrations: []*migrate.Migration{
            {Id: "000000", Up: []string{`CREATE TABLE IF NOT EXISTS users (...)`}},
            {Id: "000001", Up: []string{`CREATE TABLE IF NOT EXISTS event (...)`,
                                        `CREATE TABLE IF NOT EXISTS entry (...)`}},
        },
    }
    migrate.Exec(db.DB, "sqlite3", src, migrate.Up)
    t.Cleanup(func() { db.Close() })
    return db
}
```

### Auth in tests
Auth middleware (`APILoginVerify`) lives in `main.go` and is not used in handler tests.
Instead, set `env.User` directly on the `*env.Env` before making requests. This mirrors
exactly what the middleware does in production.

```go
e := &env.Env{..., User: &models.User{UUID: testUserUUID}}
```

To test ownership (403) cases, create a second env with a different UUID and a separate
router pointing at the same DB:

```go
otherE := newTestEnv(t, db, otherUserUUID)
otherRouter := apiRouter(otherE)
```

### Test structure
- One `*sqlx.DB` per top-level test function (`TestAPIEvents`, `TestAPIEntries`)
- Sub-tests run **sequentially** (no `t.Parallel()`) and share state — later sub-tests
  depend on IDs captured by earlier ones (e.g. create → get → update → delete)
- Use `httptest.NewRecorder()` + `router.ServeHTTP()` — no real HTTP server needed

### Helper functions
| Helper | Purpose |
|---|---|
| `newTestDB(t)` | In-memory SQLite with migrations applied |
| `newTestEnv(t, db, userUUID)` | `*env.Env` with real repo, silenced logger, set user |
| `apiRouter(e)` | Mux router wired with all `/api/v1/` routes |
| `do(t, router, method, path, body)` | Fire a request, return `*httptest.ResponseRecorder` |
| `mustDecode[T](t, body)` | Generic JSON decode helper, fails test on error |

### What to test per endpoint
- Happy path (correct status code + response body fields)
- Validation errors (missing required fields → 400, invalid UUID → 400)
- Ownership enforcement (other user's resource → 403)
- State changes are visible (create then list, delete then list)

---

## API Plan

Routes live under `/api/v1/`, protected by `APILoginVerify` (returns JSON 401, not redirect).

- [x] Add `json:` tags to `Event`, `Entry`, and `User` models
- [x] Add `DeleteEvent` to repository (only `DeleteEntry` exists today)
- [x] Create `internals/handlers/api.go` with `writeJSON`/`writeError` helpers
- [x] Implement `GET /api/v1/events` — list authenticated user's events
- [x] Implement `POST /api/v1/events` — create an event
- [x] Implement `GET /api/v1/events/{id}` — get a single event
- [x] Implement `PUT /api/v1/events/{id}` — update an event
- [x] Implement `DELETE /api/v1/events/{id}` — delete an event
- [x] Implement `GET /api/v1/events/{id}/entries` — list entries for an event
- [x] Implement `POST /api/v1/entries` — create an entry
- [x] Implement `GET /api/v1/entries/{id}` — get a single entry
- [x] Implement `PUT /api/v1/entries/{id}` — update an entry
- [x] Implement `DELETE /api/v1/entries/{id}` — delete an entry
- [x] Register all API routes and `APILoginVerify` middleware in `main.go`
- [x] Add tests for API event handlers (`TestAPIEvents` — 13 sub-tests)
- [x] Add tests for API entry handlers (`TestAPIEntries` — 17 sub-tests)
- [x] Update CLAUDE.md with testing patterns

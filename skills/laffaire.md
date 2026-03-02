# Laffaire

Use this skill to manage projects and calendar entries in a Laffaire instance via its JSON API.

In Laffaire, a **project** is called an **Event** (a named calendar group), and an individual
calendar item within it is called an **Entry**.

## Installation

Copy this file to your Claude Code commands directory:

```sh
cp skills/laffaire.md ~/.claude/commands/laffaire.md
```

Then invoke it with `/laffaire` in any Claude Code session.

---

## Usage

When the user invokes `/laffaire`, ask them what they want to do:

1. **List projects** — list all their events
2. **Create a project** — create a new named event/calendar group
3. **List entries** — list calendar entries in a project
4. **Add an entry** — add a calendar entry to an existing project
5. **Delete a project or entry**

Before making any API call, check if the user has set these values in the conversation or
environment. If not, ask for them:

- `LAFFAIRE_URL` — base URL of the Laffaire instance, e.g. `https://example.com`
- `LAFFAIRE_TOKEN` — a Bearer token created from the `/-/tokens` page of the UI

---

## API Reference

All routes are under `/api/v1/`. All requests must include:

```
Authorization: Bearer <token>
Content-Type: application/json
```

### Projects (Events)

**List all projects**
```
GET /api/v1/events
```
Returns an array of `{ id, title, description }`.

**Create a project**
```
POST /api/v1/events
{ "title": "My Project", "description": "Optional description" }
```
Returns `201` with `{ id, title, description }`. The `id` is needed to add entries.

**Get a project**
```
GET /api/v1/events/{id}
```

**Update a project**
```
PUT /api/v1/events/{id}
{ "title": "New Title", "description": "New description" }
```

**Delete a project**
```
DELETE /api/v1/events/{id}
```
Returns `204 No Content`.

---

### Entries (Calendar Items)

**List entries in a project**
```
GET /api/v1/events/{project_id}/entries
```
Returns an array of entry objects.

**Create an entry**
```
POST /api/v1/entries
{
  "event_id":     "<project id>",
  "subject":      "Team standup",
  "start_date":   "2026-03-10",
  "start_time":   "09:00",
  "end_date":     "2026-03-10",
  "end_time":     "09:30",
  "all_day_event": false,
  "description":  "Daily sync",
  "location":     "Zoom",
  "private":      false
}
```
`event_id` and `subject` are required. Returns `201` with the created entry.

**Get an entry**
```
GET /api/v1/entries/{id}
```

**Update an entry**
```
PUT /api/v1/entries/{id}
{ "subject": "...", "start_date": "...", ... }
```
`subject` is required.

**Delete an entry**
```
DELETE /api/v1/entries/{id}
```
Returns `204 No Content`.

---

## Instructions for Claude

- Use `curl` to make API calls unless the user prefers another tool.
- Always confirm the action with the user before deleting anything.
- When creating an entry, if the user does not supply dates/times, ask before defaulting.
- If the API returns a non-2xx status, show the user the error JSON and suggest next steps.
- After creating a project, offer to immediately add entries to it.

### Example curl — create a project

```sh
curl -s -X POST "$LAFFAIRE_URL/api/v1/events" \
  -H "Authorization: Bearer $LAFFAIRE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Q2 Roadmap", "description": "Planning items for Q2"}'
```

### Example curl — add an entry

```sh
curl -s -X POST "$LAFFAIRE_URL/api/v1/entries" \
  -H "Authorization: Bearer $LAFFAIRE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id":   "<project id from above>",
    "subject":    "Kickoff meeting",
    "start_date": "2026-03-10",
    "start_time": "10:00",
    "end_date":   "2026-03-10",
    "end_time":   "11:00"
  }'
```

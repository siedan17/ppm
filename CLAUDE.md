# PPM — Personal Project Management

Local web app for managing customer projects, tasks, meetings, and contacts with enforced field discipline.

## Stack

- **Backend**: Go with `net/http` (no framework)
- **Database**: SQLite via `modernc.org/sqlite` (pure Go, no CGo)
- **Migrations**: goose (SQL-based up/down migrations)
- **Query Layer**: sqlc (SQL → type-safe Go code, no ORM)
- **Frontend**: HTML + CSS + JS, server-side rendered via `html/template`
- **Markdown**: `github.com/yuin/goldmark` (GFM)
- **AI (future)**: `github.com/anthropics/anthropic-sdk-go` — interface defined, NoOp stub active
- **Single binary**: Static files, templates, and migrations embedded via `//go:embed`

## Running

```bash
go run ./cmd/ppm --port 8080 --db ppm.db
make run    # same as above
make build  # produces ./ppm binary
```

## Architecture

Clean Architecture with strict dependency direction: Domain ← Application ← Infrastructure/Delivery.

- **Domain Layer (core)**: Pure Go — entities, value objects, business rules, repository interfaces. Depends on nothing.
- **Application Layer (use cases)**: Use cases / interactors / application services. Depends on domain layer and domain-defined interfaces only. No database/sql, sqlc, HTTP imports.
- **Infrastructure Layer (outer)**: SQLite driver, sqlc-generated code, goose migrations, repository implementations, external services. Depends on domain interfaces.
- **Delivery Layer (interface)**: HTTP handlers, HTML templates, htmx, CLI commands. Depends on application layer. Never depended on by inner layers.

## Project Structure

```
cmd/ppm/main.go              Entry point, routing, dependency wiring
embed.go                     //go:embed directives for templates, static, migrations
internal/
  domain/                    Entities, value objects, business rules, repository interfaces (pure Go)
  application/               Use cases / application services (depends on domain only)
  infrastructure/
    persistence/             sqlc-generated code, repository implementations
    database/                SQLite connection, goose migration runner
    services/                External service adapters (AI, etc.)
  delivery/
    http/                    HTTP handlers, middleware, routing
    render/                  Template rendering + goldmark markdown
migrations/                  goose SQL migrations (up/down)
queries/                     sqlc SQL query files
sqlc.yaml                    sqlc configuration
templates/
  layout.html                Base HTML with nav/footer
  partials/                  nav, flash, task_row, task_list, dashboard_tasks
  pages/                     dashboard, project_*, task_*, meeting_*, person_*
static/
  css/style.css              Priority colors, status chips, overdue highlights
  js/htmx.min.js             Vendored htmx 2.0.4
```

## Key Patterns

- **Method override**: POST with `_method=PUT|PATCH|DELETE` for HTML forms
- **htmx partials**: Task filters swap `<tbody>` via `/partials/task-list`, dashboard overdue via `/partials/dashboard-tasks`
- **Inline status toggle**: Task status `<select>` fires `hx-patch="/tasks/{id}/status"`, returns updated `<tr>`
- **Per-page templates**: Each page gets its own template clone (layout + partials + page) to avoid Go template namespace conflicts
- **Meeting→Task workflow**: "Create Task from Meeting" pre-fills project_id and meeting_id

## Database

SQLite with WAL mode + foreign keys. Migrations managed by goose (up/down SQL files). Trigger-based `updated_at` auto-update on all main tables. All queries written in SQL and compiled to Go via sqlc.

## AI Integration (prepared, not implemented)

`services/agent.go` defines `AgentService` interface with:
- `ExtractTasksFromMeeting(ctx, notes, project) → []ExtractedTask`
- `SummarizeProjectStatus(ctx, projectID) → string`

Default `NoOpAgent` stub returns nil. Future: implement with Anthropic Go SDK tool calling, or build custom agentic loop.

## HTTP Routes

```
GET  /                          Dashboard
GET  /projects, /tasks, /meetings, /people   List pages
GET  /*/new                     Create forms
POST /*                         Create
GET  /*/{id}                    Detail (projects, meetings)
GET  /*/{id}/edit               Edit forms
PUT  /*/{id}                    Update
DELETE /*/{id}                  Delete
PATCH /tasks/{id}/status        Inline status toggle
POST /projects/{id}/people      Link contact
POST /meetings/{id}/participants  Add participant
GET  /meetings/{id}/create-task   Task form from meeting
GET  /meetings/{id}/export        Download .md
GET  /projects/{id}/export/meetings  Download all meetings .md
```

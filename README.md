# Movies API — Actor & Genre module (WIP)

This covers the **Actor** and **Genre** entities of the `movies-api` project. 

## Stack

- Go, `net/http` (built-in router, method-based patterns: `GET /path`, `POST /path`, etc.)
- SQLite via `database/sql` + `github.com/mattn/go-sqlite3`
- Raw SQL, no ORM

## Architecture

```
handler → service → repository → SQL → SQLite
```

- **handler** — parses the HTTP request (path/query params, body), calls the service, writes the HTTP response
- **service** — thin layer, passes calls through to the repository
- **repository** — all the SQL lives here: `SQLiteActorRepository` and `SQLiteGenreRepository`, each with its own constructor (`NewSQLiteActorRepository(db)`, etc.) and interface (`ActorRepository`, `GenreRepository`)

Repo → service → handler wiring happens in `routes.go`, inside `RegisterRoutes(mux, database)` — `main.go` only initializes the database and starts the server.

## Actor — implemented functionality

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/actors` | Create an actor. Body: `name`, `birthdate` (`YYYY-MM-DD`), optional `movieIds: []int` to link to existing movies right away |
| `GET` | `/api/actors` | List all actors. Without `movies`, filmography is not fetched (avoids unnecessary JOINs on every request) |
| `GET` | `/api/actors?page={page}&size={size}` | Paginated list of all actors. Both params are required together; response includes `page`, `size`, and `total` (total actor count, regardless of pagination) |
| `GET` | `/api/actors?movies=true` | List all actors, each including their filmography |
| `GET` | `/api/actors?name={name}` | Search by name, partial case-insensitive match (`LIKE %name%`). Always includes filmography for each match |
| `GET` | `/api/actors/{id}` | Actor by id, including filmography |
| `PATCH` | `/api/actors/{id}` | Partial update: `name`, `birthdate`, `movieIds` — any field can be omitted and won't be changed |
| `DELETE` | `/api/actors/{id}` | Delete an actor. If they have movies — 400 with a clear message. `?force=true` deletes the links in `movie_actors` first, then the actor |
| `DELETE` | `/api/actors/deleteconnection/{id}` | Remove specific movie links from an actor. Body: `{"movieIds": [1, 2]}` — only the listed links are removed, the actor and any other links stay intact |

### Implementation notes

- **Actor filmography** (`GetMovies`) uses `JOIN movies ON movies.id = movie_actors.movie_id WHERE movie_actors.actor_id = ?` — a single query instead of a loop of individual ones
- **Pagination** (`GetAll`) uses `LIMIT ? OFFSET ?` in SQL (`OFFSET = page * size`), plus a separate `SELECT COUNT(*)` for `total`. Requesting a page beyond the data just returns an empty `results` list with `200 OK` (no special-casing needed — SQL handles it naturally)
- **PATCH** takes a separate `ActorPatchRequest` struct with pointer fields (`*string`, etc.) so "field not provided" can be distinguished from "field provided empty"
- **Creating movie links** (`Create`/`Update` with `movieIds`) is wrapped in a transaction (`tx.Begin()` / `tx.Commit()` / `defer tx.Rollback()`) — if a link insert fails, everything rolls back, including the created/updated actor
- Link insertion uses `INSERT OR IGNORE INTO movie_actors(movie_id, actor_id) VALUES (?, ?)`, protected against duplicates by the composite `PRIMARY KEY(movie_id, actor_id)`
- **Force-delete**: first `DELETE FROM movie_actors WHERE actor_id = ?`, then `DELETE FROM actors WHERE id = ?` — order matters because foreign keys are enabled
- **Removing specific links** (`DeleteConnection`) builds a single `DELETE ... WHERE actor_id = ? AND movie_id IN (?, ?, ...)` query with a dynamically sized placeholder list, instead of one query per movie id

## Genre — implemented functionality

Same CRUD set and same architecture as Actor:

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/genres` | Create a genre |
| `GET` | `/api/genres` | List all genres |
| `GET` | `/api/genres?movies=true` | List all genres, each including their filmography |
| `GET` | `/api/genres/{id}` | Genre by id |
| `PATCH` | `/api/genres/{id}` | Partial update of the name |
| `DELETE` | `/api/genres/{id}` | Delete a genre. If it has linked movies — 400, `?force=true` to force delete |
| `DELETE` | `/api/genres/deleteconnection/{id}` | Remove specific movie links from a genre. Body: `{"movieIds": [1, 2]}` — same pattern as the Actor equivalent |

## Input validation

Both `entity.Actor` and `entity.Genre` have a `Validate()` method (no external library — hand-written checks), called from the service layer before create/update. Errors from multiple fields are combined via `errors.Join`, so a single call can report several problems at once (e.g. empty name and a birth date in the future).

## Errors

- Sentinel errors (`entity.ErrNotFound`, `entity.ErrInvalidContent`) are wrapped with `%w` in repository/service errors, so handlers can check the cause via `errors.Is(...)` and pick the right HTTP status (404, 400) instead of always returning 500
- All handlers return `*customerrors.HttpError` instead of writing to the response directly; a single `HttpErrorHandler.ServeHTTP` (shared with the Movie module) logs the underlying error and writes the response in one place
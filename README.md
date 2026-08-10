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

## Movie — implemented functionality

| Method   | Path                                 | Description                                                   |
| -------- | ------------------------------------ | ------------------------------------------------------------- |
| `POST`   | `/api/movie`                         | Create a movie                                                |
| `GET`    | `/api/movie`                         | List all movies                                               |
| `GET`    | `/api/movie?genre={genreId}`         | Retrieve movies filtered by genre                             |
| `GET`    | `/api/movie?year={releaseYear}`      | Retrieve movies filtered by release year                      |
| `GET`    | `/api/movie?actor={actorId}`         | Retrieve movies that the specified actor has starred in       |
| `GET`    | `/api/movie?page={page}&size={size}` | Retrieve movies with pagination                               |
| `GET`    | `/api/movie/search?title={title}`    | Search movies by title using a case-insensitive partial match |
| `GET`    | `/api/movie/{id}`                    | Retrieve a movie by ID                                        |
| `GET`    | `/api/movie/{id}/actors`             | Retrieve all actors starring in a movie                       |
| `PATCH`  | `/api/movie/{id}`                    | Partially update an existing movie                            |
| `DELETE` | `/api/movie/{id}`                    | Delete a movie                                                |

### Implementation notes

* `GET /api/movie` supports filtering by `actor`, `genre`, and `year` query parameters.
* `GET /api/movie` also supports pagination using `page` and `size`.
* When no query parameters are provided, all movies are returned.
* Pagination responses include `Page`, `Size`, and `Movies`.
* Movie title search is available through `/api/movie/search?title={title}`.
* `GET /api/movie/{id}/actors` retrieves all actors associated with a specific movie.
* Movie creation, partial updates, and deletion are implemented.
* Movie routes follow the `handler → service → repository → SQL → SQLite` architecture.
* All movie handlers are wrapped with `customerrors.HttpErrorHandler`.
* The movie repository is initialized with `NewSQLiteMovieRepository(database)`, the service with `NewMovieService(movieRepo)`, and the handler with `NewMovieHandler(movieService)`.


## Actor — implemented functionality

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/actors` | Create an actor. Body: `name`, `birthdate` (ISO 8601, `YYYY-MM-DD` or full RFC3339), optional `movieIds: []int` to link to existing movies right away |
| `GET` | `/api/actors` | List all actors. Without `movies`, filmography is not fetched (avoids unnecessary JOINs on every request) |
| `GET` | `/api/actors?movies=true` | List all actors, each including their filmography |
| `GET` | `/api/actors?name={name}` | Search by name, partial case-insensitive match (`LIKE %name%`). Always includes filmography for each match |
| `GET` | `/api/actors/{id}` | Actor by id, including filmography |
| `PATCH` | `/api/actors/{id}` | Partial update: `name`, `birthdate`, `movieIds` — any field can be omitted and won't be changed |
| `DELETE` | `/api/actors/{id}` | Delete an actor. If they have movies — 400 with a clear message. `?force=true` deletes the links in `movie_actors` first, then the actor |
| `DELETE` | `/api/actors/deleteconnection/{id}` | Remove specific movie links from an actor. Body: `{"movieIds": [1, 2]}` — only the listed links are removed, the actor and any other links stay intact |

### Implementation notes

- **Actor filmography** (`GetMovies`) uses `JOIN movies ON movies.id = movie_actors.movie_id WHERE movie_actors.actor_id = ?` — a single query instead of a loop of individual ones
- **Dates** are stored in the DB as `TEXT` in `YYYY-MM-DD` format, converted both ways via `time.Parse("2006-01-02", ...)` / `.Format("2006-01-02")`
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
| `GET` | `/api/genres/{id}` | Genre by id |
| `PATCH` | `/api/genres/{id}` | Partial update of the name |
| `DELETE` | `/api/genres/{id}` | Delete a genre. If it has linked movies — 400, `?force=true` to force delete |
| `DELETE` | `/api/genres/deleteconnection/{id}` | Remove specific movie links from a genre. Body: `{"movieIds": [1, 2]}` — same pattern as the Actor equivalent |

## Input validation

Both `entity.Actor` and `entity.Genre` have a `Validate()` method (no external library — hand-written checks), called from the service layer before create/update. Errors from multiple fields are combined via `errors.Join`, so a single call can report several problems at once (e.g. empty name and a birth date in the future).

## Known limitations / TODO

- The standard JSON decoder currently only accepts dates in RFC3339 format (`2000-09-08T00:00:00Z`) on `POST`; the simple `YYYY-MM-DD` format for incoming JSON isn't supported yet (needs a custom `UnmarshalJSON`/`MarshalJSON` for the date type)
- Pagination (`page`/`size`) for list endpoints isn't implemented yet
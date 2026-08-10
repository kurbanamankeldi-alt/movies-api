# Movies API — Actor & Genre module

Covers the **Actor** and **Genre** parts of the `movies-api` project (Movie is implemented separately by a teammate).

## Setup

```
go mod tidy
go run main.go
```

The server starts on `http://localhost:8081` and seeds the database with sample data on first run.

## Actor endpoints

### `POST /api/actors` — create an actor

**Body**
- Required: `name` (string), `birthdate` (string, `YYYY-MM-DD`)
- Optional: `movieIds` (`[]int`) — link to existing movies right away

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

**Body**
- Required: `name` (string), `birthdate` (string, `YYYY-MM-DD`)
- Optional: `movieIds` (`[]int`) — link to existing movies right away

```json
{
  "name": "Tom Hardy",
  "birthdate": "1977-09-15",
  "movieIds": [1, 5]
}
```

Returns the created actor, including its `id` and `version`.

### `GET /api/actors` — list actors

**Query params** (all optional, combine freely)
- `movies=true` — include each actor's filmography
- `page`, `size` — paginate; **must be provided together**. Response is `{"actors": [...], "page": ..., "size": ..., "total": ...}`
- `name={text}` — search by name instead of listing everyone (partial, case-insensitive). Always includes filmography

### `GET /api/actors/{id}` — get one actor

Always includes filmography.

### `PATCH /api/actors/{id}` — update an actor

**Body**
- **Required: `version` (int)** — the actor's current version (get it from a prior `GET`); the request fails without it
- Optional: `name`, `birthdate`, `movieIds` — omit any field to leave it unchanged

```json
{
  "version": 1,
  "name": "Tom Hardy Jr."
}
```

If `version` doesn't match the actor's current version in the database (someone else updated it in the meantime), the request fails with `409 Conflict`. Re-fetch the actor and try again with the new version.

### `DELETE /api/actors/{id}` — delete an actor

**Query params**
- `force=true` (optional) — deletes the actor even if they're linked to movies (removes the links too). Without it, deleting an actor who has movies fails with `400`

### `DELETE /api/actors/deleteconnection/{id}` — unlink specific movies

**Body**
- Required: `movieIds` (`[]int`) — only these links are removed; the actor and any other links stay intact

```json
{
  "movieIds": [1, 5]
}
```

## Genre endpoints

Same shape as Actor, minus pagination and version checks.

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

**Body**
- Required: `name` (string)

```json
{
  "name": "Neo-noir"
}
```

### `GET /api/genres` — list genres

**Query params** (optional)
- `movies=true` — include each genre's associated movies

### `GET /api/genres/{id}` — get one genre

Always includes associated movies.

### `PATCH /api/genres/{id}` — update a genre

**Body**
- Optional: `name`

```json
{
  "name": "Neo-noir Thriller"
}
```

### `DELETE /api/genres/{id}` — delete a genre

**Query params**
- `force=true` (optional) — same behavior as Actor: without it, deleting a genre with movies fails with `400`

### `DELETE /api/genres/deleteconnection/{id}` — unlink specific movies

**Body**
- Required: `movieIds` (`[]int`)

```json
{
  "movieIds": [1, 5]
}
```

## What to expect

- `200 OK` — successful read
- `201 Created` — successful create
- `204 No Content` — successful update to nothing returned / successful delete
- `400 Bad Request` — invalid input (missing required field, bad id, bad date, etc.) — the response body explains what's wrong
- `404 Not Found` — the actor/genre id doesn't exist
- `409 Conflict` — an actor was updated by someone else since you last fetched it (see `PATCH` above)
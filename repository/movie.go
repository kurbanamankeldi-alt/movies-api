package repository

import (
	"database/sql"
    "github.com/kurbanamankeldi-alt/movies-api/entity"
    "github.com/kurbanamankeldi-alt/movies-api/customerrors"
    sqlite3 "github.com/mattn/go-sqlite3"
    "time"
    "strings"
    "errors"
    "fmt"
)

type SQLiteMovieRepository struct {
    db *sql.DB
}

func NewSQLiteMovieRepository(db *sql.DB) *SQLiteMovieRepository{
    return &SQLiteMovieRepository{db: db}
}

type MovieRepository interface {
    FindAll() ([]*entity.Movie, error)
    FindWithPagination(page, size int) ([]*entity.Movie, error)
    FindById(id int) (*entity.Movie, error)
    FindByGenre(genreId int) ([]*entity.Movie, error)	   
    FindByYear(year int) ([]*entity.Movie, error)
    FindByActor(actorId int) ([]*entity.Movie, error) 
    FindActors(id int) ([]entity.Actor, error)
    Create(movie *entity.Movie) (int64, error)
    Update(id int, newData *entity.Movie) (int64, error)
    Delete(id int) (int64, error)
    //extra
    FindByExactTitle(title string) (*entity.Movie, error) 
    FindByTitleContains(title string) ([]*entity.Movie, error)
}

func (r *SQLiteMovieRepository) FindAll() ([]*entity.Movie, error) {
    queryMoviesTable := `SELECT * FROM movies ORDER BY Id`

    rows, err := r.db.Query(queryMoviesTable)
    if err != nil {
        return nil, fmt.Errorf("%w: select movies: %w", customerrors.ErrDB, err)
    }
    defer rows.Close()

    movies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            return nil, fmt.Errorf("%w: scan movie row: %w", customerrors.ErrDB, err)
        }
        movies = append(movies, m)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: iterate movie rows: %w", customerrors.ErrDB, err)
    }

    if err := r.populateRelations(movies); err != nil {
        return nil, err
    }

    return movies, nil

}

func (r *SQLiteMovieRepository) FindWithPagination(page, size int) ([]*entity.Movie, error) {
    offset := page * size

    queryMoviesTable := `SELECT * FROM movies ORDER BY Id LIMIT ? OFFSET ?`

    rows, err := r.db.Query(queryMoviesTable, size, offset)
    if err != nil {
        return nil, fmt.Errorf("%w: select movies: %w", customerrors.ErrDB, err)
    }
    defer rows.Close()

    movies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            return nil, fmt.Errorf("%w: scan movie row: %w", customerrors.ErrDB, err)
        }
        movies = append(movies, m)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: iterate movie rows: %w", customerrors.ErrDB, err)
    }

    if err := r.populateRelations(movies); err != nil {
        return nil, err
    }

    return movies, nil

}

func (r *SQLiteMovieRepository) FindById(id int) (*entity.Movie, error) {
    queryMoviesTable := `SELECT * FROM movies WHERE id = ?`
    row := r.db.QueryRow(queryMoviesTable, id)
    movie := &entity.Movie{}

    err := row.Scan(&movie.Id,&movie.Title,&movie.ReleaseYear, &movie.Duration)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("%w: not found (movieId=%d): %w", customerrors.ErrNotFound, id, err)
        }
        return nil, fmt.Errorf("%w: select movie by id: %w", customerrors.ErrDB, err)
    }

    if err := r.populateRelations([]*entity.Movie{movie}); err != nil {
        return nil, err
    }

    return movie, nil
}

func (r *SQLiteMovieRepository) FindByGenre(genreId int) ([]*entity.Movie, error) {

    query := `
        SELECT 
            m.id, 
            m.title, 
            m.release_year,
            m.duration
        FROM movies m
        JOIN movie_genres mg
        ON mg.movie_id = m.id
        WHERE mg.genre_id = ?
        ORDER BY m.id;
    `

    rows, err := r.db.Query(query, genreId)
    if err != nil {
        return nil, fmt.Errorf("%w: select movies by genre: %w", customerrors.ErrDB, err)
    }
    defer rows.Close()

    foundMovies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return nil, fmt.Errorf("%w: not found (genreId=%d): %w", customerrors.ErrNotFound, genreId, err)
            }
            return nil, fmt.Errorf("%w: scan movie row: %w", customerrors.ErrDB, err)
        }
        foundMovies = append(foundMovies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: iterate movie rows: %w", customerrors.ErrDB, err)
    } 

    if err := r.populateRelations(foundMovies); err != nil {
        return nil, err
    }

    return foundMovies, nil
}

func (r *SQLiteMovieRepository) FindByYear(year int) ([]*entity.Movie, error) {

    query := `SELECT * FROM movies WHERE release_year = ?`

    rows, err := r.db.Query(query, year)
    if err != nil {
        return nil, fmt.Errorf("%w: select movies by year: %w", customerrors.ErrDB, err)
    }
    defer rows.Close()

    movies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}

        if err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration); err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return nil, fmt.Errorf("%w: not found (year=%d): %w", customerrors.ErrNotFound, year, err)
            }
            return nil, fmt.Errorf("%w: scan movie row: %w", customerrors.ErrDB, err)
        }
        movies = append(movies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: iterate movie rows: %w", customerrors.ErrDB, err)
    }

    if err := r.populateRelations(movies); err != nil {
        return nil, err
    }

    return movies, nil
}

func (r *SQLiteMovieRepository) FindByActor(actorId int) ([]*entity.Movie, error) {

    query := `
        SELECT 
            m.id, 
            m.title, 
            m.release_year,
            m.duration
        FROM movies m
        JOIN movie_actors ma
        ON ma.movie_id = m.id
        WHERE ma.actor_id = ?
        ORDER BY m.id;
    `

    rows, err := r.db.Query(query, actorId)
    if err != nil {
        return nil, fmt.Errorf("%w: select movies by actor: %w", customerrors.ErrDB, err)
    }
    defer rows.Close()

    foundMovies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return nil, fmt.Errorf("%w: not found (actorId=%d): %w", customerrors.ErrNotFound, actorId, err)
            }
            return nil, fmt.Errorf("%w: scan movie row: %w", customerrors.ErrDB, err)
        }
        foundMovies = append(foundMovies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: iterate movie rows: %w", customerrors.ErrDB, err)
    } 

    if err := r.populateRelations(foundMovies); err != nil {
        return nil, err
    }

    return foundMovies, nil
}

func (r *SQLiteMovieRepository) FindActors(id int) ([]entity.Actor, error) {
    movie, err := r.FindById(id)

    if err != nil {
        return nil, err
    }

    actors := []entity.Actor{}    

    for _, actor := range movie.Actors {
        actors = append(actors, actor)
    }

    return actors, nil
}

func (r *SQLiteMovieRepository) Create(movie *entity.Movie) (int64, error) {
    tx, err := r.db.Begin()
    if err != nil {
        return 0, fmt.Errorf("%w: transaction movie create start: %w", customerrors.ErrDB, err)
    }

    defer tx.Rollback()

    queryForMovies := `INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?);`
    result, err := tx.Exec(queryForMovies, movie.Title, movie.ReleaseYear, movie.Duration)

    if err != nil {
        return 0, fmt.Errorf("%w: insert movie: %w", customerrors.ErrDB, err)
    }

    movieId, err := result.LastInsertId()

    if err != nil {
        return 0, fmt.Errorf("%w: get inserted movie id: %w", customerrors.ErrDB, err)
    }

    queryForMovieActors := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?);`

    for _, actor := range movie.Actors {
        _, err := tx.Exec(queryForMovieActors, movieId, actor.Id)
        if err != nil {
            var sqliteErr sqlite3.Error
            if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey  {
                return 0, fmt.Errorf("%w: non existing (actorId=%d): %w", customerrors.ErrInvalidReference, actor.Id, err)
            }
            return 0, fmt.Errorf("%w: insert movie_actors (movie=%d, actor=%d): %w", customerrors.ErrDB, movieId, actor.Id, err)
        }        
    }
    
    queryForMovieGenres := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?);`

    for _, genre := range movie.Genres {
        _, err := tx.Exec(queryForMovieGenres, movieId, genre.Id)
        if err != nil {
            var sqliteErr sqlite3.Error
            if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey  {
                return 0, fmt.Errorf("%w: non existing (genreId=%d): %w", customerrors.ErrInvalidReference, genre.Id, err)
            }
            return 0, fmt.Errorf("%w: insert movie_genres (movie=%d, genre=%d): %w", customerrors.ErrDB, movieId, genre.Id, err)
        }      
    }    
    
    err = tx.Commit()
    if err != nil {
        return 0, fmt.Errorf("%w: transaction movie create commit: %w", customerrors.ErrDB, err)
    }

    movie.Id = uint(movieId)    
    
    return movieId, nil
}

func (r *SQLiteMovieRepository) Update(id int, newData *entity.Movie) (int64, error) {

    tx, err := r.db.Begin()
    if err != nil {
        return 0, fmt.Errorf("%w: transaction movie start: %w", customerrors.ErrDB, err)
    }
    defer tx.Rollback()

    queryForMovies := `UPDATE movies 
                       SET title = ?, release_year = ?, duration = ? 
                       WHERE id = ?;`

    result, err := tx.Exec(queryForMovies, newData.Title, newData.ReleaseYear, newData.Duration, id)  
    
    if err != nil {
        return 0, fmt.Errorf("%w: update movie (movie=%d): %w", customerrors.ErrDB, id, err)
    }                    

    rows, err := result.RowsAffected()
    if err != nil {
        return 0, fmt.Errorf("%w: iget affected rows for movie update: %w", customerrors.ErrDB, rows, err)
    }

    if rows == 0 {
        return 0, fmt.Errorf("%w: movie does not exist (movie=%d) %w", customerrors.ErrNotFound, err)
    }    

    _, err = tx.Exec(`DELETE FROM movie_actors WHERE movie_id = ?`, id)
    if err != nil {
        return 0, fmt.Errorf("%w: delete movie_actor failed (movie_id=%d) %w", customerrors.ErrDB, id, err)
    }   

    _, err = tx.Exec(`DELETE FROM movie_genres WHERE movie_id = ?`, id)
    if err != nil {
        return 0, fmt.Errorf("%w: delete movie_genre failed (movie_id=%d) %w", customerrors.ErrDB, id, err)
    }   


    queryForMovieActors := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?);`    

	for _, actor := range newData.Actors {
		_, err := tx.Exec(queryForMovieActors, id, actor.Id)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) &&
				sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return 0, fmt.Errorf(
					"%w: non existing (actorId=%d): %w",
					customerrors.ErrInvalidReference,
					actor.Id,
					err,
				)
			}

			return 0, fmt.Errorf(
				"%w: insert movie_actors (movie=%d, actor=%d): %w",
				customerrors.ErrDB,
				id,
				actor.Id,
				err,
			)
		}
	}
    
    queryForMovieGenres := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?);`     
      
	for _, genre := range newData.Genres {
		_, err := tx.Exec(queryForMovieGenres, id, genre.Id)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) &&
				sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return 0, fmt.Errorf(
					"%w: non existing (genreId=%d): %w",
					customerrors.ErrInvalidReference,
					genre.Id,
					err,
				)
			}

			return 0, fmt.Errorf(
				"%w: insert movie_genres (movie=%d, genre=%d): %w",
				customerrors.ErrDB,
				id,
				genre.Id,
				err,
			)
		}
	}

    if err := tx.Commit(); err != nil {
        return 0, fmt.Errorf("%w: transaction movie update commit: %w", customerrors.ErrDB, err)
    }

    return rows, nil
}

func (r *SQLiteMovieRepository) Delete(id int) (int64, error) {
    tx, err := r.db.Begin()
    if err != nil {
        return 0, fmt.Errorf("%w: transaction movie delete start: %w", customerrors.ErrDB, err)
    }
    defer tx.Rollback()

    queryMovieActorsTable := `DELETE FROM movie_actors WHERE movie_id = ?`
    queryMovieGenresTable := `DELETE FROM movie_genres WHERE movie_id = ?`
    queryForMoviesTable := `DELETE FROM movies WHERE id = ?` 

    _, err = tx.Exec(queryMovieActorsTable, id)
    if err != nil {
        return 0, fmt.Errorf("%w: delete movie_actors (movie=%d): %w", customerrors.ErrDB, id, err)
    }

    _, err = tx.Exec(queryMovieGenresTable, id)
    if err != nil {
        return 0, fmt.Errorf("%w: delete movie_genres (movie=%d): %w", customerrors.ErrDB, id, err)
    }   

    result, err := tx.Exec(queryForMoviesTable, id)
    if err != nil {
        return 0, fmt.Errorf("%w: delete movies (movie=%d): %w", customerrors.ErrDB, id, err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return 0, fmt.Errorf("%w: get affected rows after movie delete: %w", customerrors.ErrDB, err)
    }

    if rows == 0 {
        return 0, fmt.Errorf("%w: movie does not exist (movie=%d)", customerrors.ErrNotFound, id)
    }

    if err := tx.Commit(); err != nil {
        return 0, fmt.Errorf("%w: transaction movie delete commit:  %w", customerrors.ErrDB, err)
    }

    return rows, nil
}

//extra
func (r *SQLiteMovieRepository) FindByExactTitle(title string) (*entity.Movie, error) {
    queryMoviesTable := `SELECT * FROM movies WHERE Title = ?`
    
    row := r.db.QueryRow(queryMoviesTable, title)
    movie := &entity.Movie{}

    err := row.Scan(&movie.Id,&movie.Title,&movie.ReleaseYear, &movie.Duration)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("%w: select movie by id: %w", customerrors.ErrDB, err)
    }

    if err := r.populateRelations([]*entity.Movie{movie}); err != nil {
        return nil, err
    }

    return movie, nil

}

func (r *SQLiteMovieRepository) FindByTitleContains(title string) ([]*entity.Movie, error) {
    queryMoviesTable := `SELECT * FROM movies WHERE LOWER(title) LIKE ? ORDER BY Id`
    
    rows, err := r.db.Query(queryMoviesTable, "%"+title+"%")
    if err != nil {
        return nil, fmt.Errorf("%w: select movies: %w", customerrors.ErrDB, err)
    }
    defer rows.Close()

    movies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            return nil, fmt.Errorf("%w: scan movie row: %w", customerrors.ErrDB, err)
        }
        movies = append(movies, m)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: iterate movie rows: %w", customerrors.ErrDB, err)
    }

    if err := r.populateRelations(movies); err != nil {
        return nil, err
    }

    return movies, nil

}

//Helpers
func (r *SQLiteMovieRepository) GetActorsByMovieIds(movieIds []uint) (map[uint][]entity.Actor, error) {
    
    actorsByMovie := make(map[uint][]entity.Actor)

    if len(movieIds) == 0 {
        return actorsByMovie, nil
    }

    placeholders := make([]string, len(movieIds))
    args := make([]any, len(movieIds))

    for i, id := range movieIds {
        placeholders[i] = "?"
        args[i] = id
    }

    queryMoviesAndActors := fmt.Sprintf(`
        SELECT 
            ma.movie_id, 
            a.id, 
            a.name, 
            a.birthdate
        FROM movie_actors ma
        JOIN actors a ON a.id = ma.actor_id
        WHERE ma.movie_id IN (%s)
        ORDER BY ma.movie_id;
    `, strings.Join(placeholders, ","))

    rows, err := r.db.Query(queryMoviesAndActors, args...)
    if err != nil {
        return nil, fmt.Errorf("%w: select actors by movie ids=%v: %w", customerrors.ErrDB, movieIds, err)
    }
    defer rows.Close()

    for rows.Next() {
        var (
            movieId uint
            actor entity.Actor
            birthdate string
        )

        if err := rows.Scan(&movieId, &actor.Id, &actor.Name, &birthdate); err !=nil {
			return nil, fmt.Errorf("%w: scan actor for movie=%d: %w", customerrors.ErrDB, movieId, err)
        }

        actor.BirthDate, err = time.Parse("2006-01-02", birthdate)

        if err != nil {
			return nil, fmt.Errorf("%w: invalid actor birthdate (actor=%d, value=%s): %w",
                customerrors.ErrInvalidReference,
				actor.Id,
				birthdate,
				err,
			)
        }

        actorsByMovie[movieId] = append(actorsByMovie[movieId], actor)
    }    

    if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"%w: iterate actors by movie ids=%v: %w",
			customerrors.ErrDB,
			movieIds,
			err,
		)
    }    

    return actorsByMovie, nil
    
}

func (r *SQLiteMovieRepository) GetGenresByMovieIds(movieIds []uint) (map[uint][]entity.Genre, error) {
   
    genresByMovie := make(map[uint][]entity.Genre)

    if len(movieIds) == 0 {
        return genresByMovie, nil
    }

    placeholders := make([]string, len(movieIds))
    args := make([]any, len(movieIds))

    for i, id := range movieIds {
        placeholders[i] = "?"
        args[i] = id
    }

    queryMovieGenres := fmt.Sprintf(`
        SELECT 
            mg.movie_id,
            g.id,
            g.name
        FROM movie_genres mg
        JOIN genres g ON g.id = mg.genre_id
        WHERE mg.movie_id IN (%s)
        ORDER BY mg.movie_id
    `, strings.Join(placeholders, ","))

    rows, err := r.db.Query(queryMovieGenres, args...)
    if err != nil {
		return nil, fmt.Errorf("%w: select genres by movie ids=%v: %w",	customerrors.ErrDB,	movieIds, err)
    }
    defer rows.Close()    

    for rows.Next() {
        var (
            movieId uint
            genre entity.Genre
        )

        if err := rows.Scan(&movieId, &genre.Id, &genre.Name); err != nil {
			return nil, fmt.Errorf("%w: scan genre for movie=%d: %w", customerrors.ErrDB, movieId,err)
        }

        genresByMovie[movieId] = append(genresByMovie[movieId], genre)
    }

    if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: iterate genres by movie ids=%v: %w", customerrors.ErrDB,	movieIds, err)
    }    

    return genresByMovie, nil
}

func (r *SQLiteMovieRepository) populateRelations(movies []*entity.Movie) error {
    if len(movies) == 0 {
        return nil
    }

    movieIds := []uint{}
    for _, movie := range movies {
        movieIds = append(movieIds, movie.Id)
    }

    actorsByMovieIds, err := r.GetActorsByMovieIds(movieIds)
    if err != nil {
        return err
    }

    genresByMovieIds, err := r.GetGenresByMovieIds(movieIds)
    if err != nil {
        return err
    }    

    for _, movie := range movies {
        movie.Actors = actorsByMovieIds[movie.Id]
        movie.Genres = genresByMovieIds[movie.Id]
    }    

    return nil
}

func containsActor(actors []entity.Actor, actorId int) bool {
    for _, actor := range actors {
        if actor.Id == uint(actorId) {
            return true
        }
    }

    return false
}

func containsGenre(genres []entity.Genre, genreId int) bool {
    for _, genre := range genres {
        if genre.Id == uint(genreId) {
            return true
        }
    }

    return false
}
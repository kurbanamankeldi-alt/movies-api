package repository

import (
	"database/sql"
    "github.com/kurbanamankeldi-alt/movies-api/entity"
    "time"
    "strings"
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
    FindById(id int) (*entity.Movie, error)
    FindByGenre(genreId int) ([]*entity.Movie, error)	   
    FindByYear(year int) ([]*entity.Movie, error)
    FindByActor(actorId int) ([]*entity.Movie, error) 
    FindActors(id int) ([]entity.Actor, error)
    Create(movie *entity.Movie) (int64, error)
    Update(id int, newData *entity.Movie) (int64, error)
    Delete(id int) (int64, error)
    //extra method
    FilterBy(movieId, actorId, genreId, year int) ([]*entity.Movie, error)
}

func (r *SQLiteMovieRepository) FindAll() ([]*entity.Movie, error) {
    queryMoviesTable := `SELECT * FROM movies ORDER BY Id`

    rows, err := r.db.Query(queryMoviesTable)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    movies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            return nil, err
        }
        movies = append(movies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, err
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

    err := row.Scan(&movie.Id,&movie.Title,&movie.ReleaseYear,&movie.Duration)
    if err != nil {
        return nil, err
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
        return nil, err
    }
    defer rows.Close()

    foundMovies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            return nil, err
        }
        foundMovies = append(foundMovies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, err
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
        return nil, err
    }
    defer rows.Close()

    movies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}

        if err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration); err != nil {
            return nil, err
        }

        movies = append(movies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, err
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
        return nil, err
    }
    defer rows.Close()

    foundMovies := []*entity.Movie{}

    for rows.Next() {
        m := &entity.Movie{}
        err := rows.Scan(&m.Id, &m.Title, &m.ReleaseYear, &m.Duration)
        if err != nil {
            return nil, err
        }
        foundMovies = append(foundMovies, m)
    }

    if err = rows.Err(); err != nil {
        return nil, err
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
        return 0, err
    }

    defer tx.Rollback()

    queryForMovies := `INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?);`
    result, err := tx.Exec(queryForMovies, movie.Title, movie.ReleaseYear, movie.Duration)

    if err != nil {
        return 0, err
    }

    movieId, err := result.LastInsertId()

    if err != nil {
        return 0, err
    }

    queryForMovieActors := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?);`

    for _, actor := range movie.Actors {
        _, err := tx.Exec(queryForMovieActors, movieId, actor.Id)
        if err != nil {
            return 0, fmt.Errorf("failed inserting actor id %v, %w", actor.Id, err)
        }        
    }
    
    queryForMovieGenres := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?);`

    for _, genre := range movie.Genres {
        _, err := tx.Exec(queryForMovieGenres, movieId, genre.Id)
        if err != nil {
            return 0, fmt.Errorf("failed inserting genre id %v,%w", genre.Id, err)
        }        
    }    
    
    err = tx.Commit()
    if err != nil {
        return 0, err
    }

    movie.Id = uint(movieId)    
    
    return movieId, nil
}

func (r *SQLiteMovieRepository) Update(id int, newData *entity.Movie) (int64, error) {

    tx, err := r.db.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    queryForMovies := `UPDATE movies 
                       SET title = ?, release_year = ?, duration = ? 
                       WHERE id = ?;`

    result, err := tx.Exec(queryForMovies, newData.Title, newData.ReleaseYear, newData.Duration, id)  
    
    if err != nil {
        return 0, err
    }                    

    rows, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

    if rows == 0 {
        return 0, fmt.Errorf("movie with id %d not found", id)
    }    

    _, err = tx.Exec(`DELETE FROM movie_actors WHERE movie_id = ?`, id)
    if err != nil {
        return 0, err
    }   

    _, err = tx.Exec(`DELETE FROM movie_genres WHERE movie_id = ?`, id)
    if err != nil {
        return 0, err
    }   


    queryForMovieActors := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?);`    

    for _, actor := range newData.Actors {
        _, err := tx.Exec(queryForMovieActors, id, actor.Id)
        if err != nil {
            return 0, err
        }
    }
    
    queryForMovieGenres := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?);`     
      
    for _, genre := range newData.Genres {
        _, err := tx.Exec(queryForMovieGenres, id, genre.Id)
        if err != nil {
            return 0, err
        }
    }

    if err := tx.Commit(); err != nil {
        return 0, err
    }

    return rows, nil
}

func (r *SQLiteMovieRepository) Delete(id int) (int64, error) {
    tx, err := r.db.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    queryMovieActorsTable := `DELETE FROM movie_actors WHERE movie_id = ?`
    queryMovieGenresTable := `DELETE FROM movie_genres WHERE movie_id = ?`
    queryForMoviesTable := `DELETE FROM movies WHERE id = ?` 

    _, err = tx.Exec(queryMovieActorsTable, id)
    if err != nil {
        return 0, err
    }

    _, err = tx.Exec(queryMovieGenresTable, id)
    if err != nil {
        return 0, err
    }   

    result, err := tx.Exec(queryForMoviesTable, id)
    if err != nil {
        return 0, err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

    if rows == 0 {
        return 0, fmt.Errorf("movie with id %d not found", id)
    }

    if err := tx.Commit(); err != nil {
        return 0, err
    }

    return rows, nil
}

//extra
func (r *SQLiteMovieRepository) FilterBy(movieId, actorId, genreId, year int) ([]*entity.Movie, error) {

    movies, err := r.FindAll()

    if err != nil {
        return nil, err
    }

    if movieId == 0 && actorId == 0 && genreId == 0 && year == 0 {
        return movies, nil
    }

    filtered := []*entity.Movie{}

    for _, movie := range movies {      
        if movieId != 0 && movie.Id != uint(movieId) {
            continue
        }  

        if actorId != 0 && !containsActor(movie.Actors, actorId) {
            continue
        }        
        if genreId != 0 && !containsGenre(movie.Genres, genreId) {
            continue
        }      
        if year != 0 && movie.ReleaseYear != year {
            continue
        }                  

        filtered = append(filtered, movie)
    }

    return filtered, nil
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
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var (
            movieId uint
            actor entity.Actor
            birthdate string
        )

        if err := rows.Scan(&movieId, &actor.Id, &actor.Name, &birthdate); err !=nil {
            return nil, err
        }

        actor.BirthDate, err = time.Parse("2006-01-02", birthdate)

        if err != nil {
            return nil, err
        }

        actorsByMovie[movieId] = append(actorsByMovie[movieId], actor)
    }    

    if err := rows.Err(); err != nil {
        return nil, err
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
        return nil, err
    }
    defer rows.Close()    

    for rows.Next() {
        var (
            movieId uint
            genre entity.Genre
        )

        if err := rows.Scan(&movieId, &genre.Id, &genre.Name); err != nil {
            return nil, err
        }

        genresByMovie[movieId] = append(genresByMovie[movieId], genre)
    }

    if err := rows.Err(); err != nil {
        return nil, err
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
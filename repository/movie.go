package repository

import (
	"database/sql"
    "github.com/kurbanamankeldi-alt/movies-api/entity"
    "time"
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
    FindByActor(actorId int) ([]*entity.Movie, error)
 	//FindByGenre(genre string) ([]*entity.Movie, error)	    
    Create(movie *entity.Movie) (int64, error)
    //Update(movie *entity.Movie) error
    //Delete(id string) error
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

    for _, movie := range movies {
        actors, err := r.GetActorsByMovieId(movie.Id)
        if err != nil {
            return nil, err
        }
        genres, err := r.GetGenresByMovieId(movie.Id)
        if err != nil {
            return nil, err
        }        
        movie.Actors = append(movie.Actors, actors...)
        movie.Genres = append(movie.Genres, genres...)
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

    actors, err := r.GetActorsByMovieId(uint(id))

    if err != nil {
        return nil, err
    }

    genres, err := r.GetGenresByMovieId(uint(id))

    if err != nil {
        return nil, err
    }

    movie.Actors = actors
    movie.Genres = genres

    return movie, nil
}

func (r *SQLiteMovieRepository) FindByActor(actorId int) ([]*entity.Movie, error) {

    allMovies, err := r.FindAll()

    if err != nil {
        return nil, err
    }

    filtered := []*entity.Movie{}

    for _, movie := range allMovies {
        actors := movie.Actors
        for _, actor := range actors {
            if actor.Id == uint(actorId) {
                filtered = append(filtered, movie)
            }
        }
    }

    return filtered, nil
}

func (r *SQLiteMovieRepository) Create(movie *entity.Movie) (int64, error) {
    sql := `INSERT INTO movies (title, release_year, duration) 
        VALUES (?, ?, ?);`
    result, err := r.db.Exec(sql, movie.Title, movie.ReleaseYear, movie.Duration)

    if err != nil {
        return 0, err
    }
    
    return result.LastInsertId()
}

//Helpers
func (r *SQLiteMovieRepository) GetActorsByMovieId(id uint) ([]entity.Actor, error) {
    
    queryMovieActors := `SELECT actor_id FROM movie_actors WHERE movie_id = ?`
    actorRows, err := r.db.Query(queryMovieActors, id)
    if err != nil {
        return nil, err
    }
    defer actorRows.Close()

    var ids []int

    for actorRows.Next() {
        var actorId int
        err := actorRows.Scan(&actorId)
        if err != nil {
            return nil, err
        }
        ids = append(ids, actorId)
    }

    var actors []entity.Actor

    queryActorsTable := `SELECT * FROM actors WHERE id = ?`

    for _, id := range ids {    
        row := r.db.QueryRow(queryActorsTable, id)
        var birthDateStr string
        actor := entity.Actor{}
        
        err := row.Scan(&actor.Id,&actor.Name,&birthDateStr)
        if err != nil {
            return nil, err
        }

        t, err := time.Parse("2006-01-02", birthDateStr)
        if err != nil {
            return nil, err
        }

        actor.BirthDate = t
        actors = append(actors, actor)
    }

    return actors, nil
}

func (r *SQLiteMovieRepository) GetGenresByMovieId(id uint) ([]entity.Genre, error) {
    
    queryMovieGenres := `SELECT genre_id FROM movie_genres WHERE movie_id = ?`
    genreRows, err := r.db.Query(queryMovieGenres, id)
    if err != nil {
        return nil, err
    }
    defer genreRows.Close()

    var ids []int

    for genreRows.Next() {
        var genreId int
        err := genreRows.Scan(&genreId)
        if err != nil {
            return nil, err
        }
        ids = append(ids, genreId)
    }

    var genres []entity.Genre
    
    queryGenresTable := `SELECT * FROM genres WHERE id = ?`

    for _, id := range ids {    
        row := r.db.QueryRow(queryGenresTable, id)
        genre := entity.Genre{}
        
        err := row.Scan(&genre.Id,&genre.Name)
        if err != nil {
            return nil, err
        }

        genres = append(genres, genre)
    }

    return genres, nil
}
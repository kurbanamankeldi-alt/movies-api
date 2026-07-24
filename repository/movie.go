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
    Create(movie *entity.Movie) (int64, error)
    FindById(id int) (*entity.Movie, error)
    //FindAll() ([]*entity.Movie, error)
    //FindByActor(actor string) ([]*entity.Movie, error)
 	//FindByGenre(genre string) ([]*entity.Movie, error)	
    //Update(movie *entity.Movie) error
    //Delete(id string) error
}

func (r *SQLiteMovieRepository) FindById(id int) (*entity.Movie, error) {
    queryMoviesTable := `SELECT * FROM movies WHERE id = ?`
    row := r.db.QueryRow(queryMoviesTable, id)
    movie := &entity.Movie{}

    err := row.Scan(&movie.Id,&movie.Title,&movie.ReleaseYear,&movie.Duration)
    if err != nil {
        return nil, err
    }

    //get ids for actors
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

    actors, err := r.GetActorsForMovieId(ids)

    if err != nil {
        return nil, err
    }

    //get ids for genres
    queryMovieGenres := `SELECT genre_id FROM movie_genres WHERE movie_id = ?`
    genreRows, err := r.db.Query(queryMovieGenres, id)
    if err != nil {
        return nil, err
    }
    defer genreRows.Close()

    var movieIds []int

    for genreRows.Next() {
        var movieId int
        err := genreRows.Scan(&movieId)
        if err != nil {
            return nil, err
        }
        movieIds = append(movieIds, movieId)
    }

    genres, err := r.GetGenresForMovieId(movieIds)

    if err != nil {
        return nil, err
    }

    movie.Actors = actors
    movie.Genres = genres

    return movie, nil
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
func (r *SQLiteMovieRepository) GetActorsForMovieId(ids []int) ([]entity.Actor, error) {
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

func (r *SQLiteMovieRepository) GetGenresForMovieId(ids []int) ([]entity.Genre, error) {
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
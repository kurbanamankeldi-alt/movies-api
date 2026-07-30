package repository

import (
	"database/sql"
    "github.com/kurbanamankeldi-alt/movies-api/entity"
    "time"
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

func (r *SQLiteMovieRepository) FindByGenre(genreId int) ([]*entity.Movie, error) {

    allMovies, err := r.FindAll()

    if err != nil {
        return nil, err
    }

    filtered := []*entity.Movie{}

    for _, movie := range allMovies {
        genres := movie.Genres
        for _, genre := range genres {
            if genre.Id == uint(genreId) {
                filtered = append(filtered, movie)
            }
        }
    }

    return filtered, nil
}

func (r *SQLiteMovieRepository) FindByYear(year int) ([]*entity.Movie, error) {

    allMovies, err := r.FindAll()

    if err != nil {
        return nil, err
    }

    filtered := []*entity.Movie{}

    for _, movie := range allMovies {
        if movie.ReleaseYear == year {
            filtered = append(filtered, movie)
        }
    }

    return filtered, nil
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
    queryForMovies := `INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?);`
    result, err := r.db.Exec(queryForMovies, movie.Title, movie.ReleaseYear, movie.Duration)

    if err != nil {
        return 0, err
    }

    movieId, _ := result.LastInsertId()

    //delete records if actor id or genre is wrong
    queryDeleteMovie := `DELETE FROM movies WHERE id = ?`
    queryDeleteActor := `DELETE FROM movie_actors WHERE movie_id = ?`

    queryForMovieActors := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?);`

    for _, actor := range movie.Actors {
        _, err := r.db.Exec(queryForMovieActors, movieId, actor.Id)
        if err != nil {
            _, _ = r.db.Exec(queryDeleteMovie, movieId)
            return 0, fmt.Errorf("the actor with id %v does not exist", actor.Id)
        }        
    }    

    queryForMovieGenres := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?);`

    for _, genre := range movie.Genres {
        _, err := r.db.Exec(queryForMovieGenres, movieId, genre.Id)
        if err != nil {
             _, _  = r.db.Exec(queryDeleteActor, movieId)         
             _, _  = r.db.Exec(queryDeleteMovie, movieId)
             return 0, fmt.Errorf("the genre with id %v does not exist", genre.Id)
        }        
    }    
    
    return result.LastInsertId()
}

func (r *SQLiteMovieRepository) Update(id int, newData *entity.Movie) (int64, error) {

    queryForMovies := `UPDATE movies 
                       SET title = ?, release_year = ?, duration = ? 
                       WHERE id = ?;`

    queryForMovieActorsDelete := `DELETE FROM movie_actors WHERE movie_id = ?`
    queryForMovieGenresDelete := `DELETE FROM movie_genres WHERE movie_id = ?`
                       
    queryForMovieActors := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?);`
    queryForMovieGenres := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?);`                         

    result, err := r.db.Exec(queryForMovies, newData.Title, newData.ReleaseYear, newData.Duration, id)

    actors, genres := newData.Actors, newData.Genres
    actorIds, genreIds := []uint{}, []uint{}

    for _, actor := range actors {
        actorIds = append(actorIds, actor.Id)
    }

    for _, genre := range genres {
        genreIds = append(genreIds, genre.Id) //10 11 12 1 3 --> 9 10 2 4
    }

    if len(actorIds) > 0 {
        _, _ = r.db.Exec(queryForMovieActorsDelete, id)
        for _, actorId := range actorIds {
            _, err := r.db.Exec(queryForMovieActors, id, actorId)
            if err != nil {
                return 0, err
            }
        }
    }

    if len(genreIds) > 0 {
        _, _ = r.db.Exec(queryForMovieGenresDelete, id)
        for _, genreId := range genreIds {
            _, err := r.db.Exec(queryForMovieGenres, id, genreId)
            if err != nil {
                return 0, err
            }
        }
    }    


    if err != nil {
        return 0, err
    }

    return result.RowsAffected()
}

func (r *SQLiteMovieRepository) Delete(id int) (int64, error) {
    queryForMoviesTable := `DELETE FROM movies WHERE id = ?`
    queryMovieActorsTable := `DELETE FROM movie_actors WHERE movie_id = ?`
    queryMovieGenresTable := `DELETE FROM movie_genres WHERE movie_id = ?`

    _, errFromMovieActors := r.db.Exec(queryMovieActorsTable, id)
    if errFromMovieActors != nil {
        return 0, errFromMovieActors
    }

    _, errFromMovieGenres := r.db.Exec(queryMovieGenresTable, id)
    if errFromMovieGenres != nil {
        return 0, errFromMovieActors
    }

    result, err := r.db.Exec(queryForMoviesTable, id)
    if err != nil {
        fmt.Println("here")
        return 0, err
    }


    return result.RowsAffected()
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

        actors, err1 := r.GetActorsByMovieId(movie.Id)
        genres, err2 := r.GetGenresByMovieId(movie.Id)

        if err1 != nil {
            return nil, err1
        }

        if err2 != nil {
            return nil, err2
        }          

        if actorId != 0 && !containsActor(actors, actorId) {
            continue
        }        
        if genreId != 0 && !containsGenre(genres, genreId) {
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

//10 9 10 | 2 4
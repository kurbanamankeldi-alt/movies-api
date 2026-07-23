package repository

import (
	"database/sql"
    "github.com/kurbanamankeldi-alt/movies-api/entity"
)

type SQLiteMovieRepository struct {
    db *sql.DB
}

func NewSQLiteMovieRepository(db *sql.DB) *SQLiteMovieRepository{
    return &SQLiteMovieRepository{db: db}
}

type MovieRepository interface {
    Create(movie *entity.Movie) (int64, error)
    //FindByID(id string) (*entity.Movie, error)
    //FindAll() ([]*entity.Movie, error)
    //FindByActor(actor string) ([]*entity.Movie, error)
 	//FindByGenre(genre string) ([]*entity.Movie, error)	
    //Update(movie *entity.Movie) error
    //Delete(id string) error
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
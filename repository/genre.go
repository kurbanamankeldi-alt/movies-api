package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
)

type SQLiteGenreRepository struct {
	db *sql.DB
}

func NewSQLiteGenreRepository(db *sql.DB) *SQLiteGenreRepository {
	return &SQLiteGenreRepository{db: db}
}

type GenreRepository interface {
	Create(genre *entity.Genre) (int64, error)
	GetAll(moviesFlag bool) ([]entity.Genre, error)
	GetByID(id int) (entity.Genre, error)
	Update(id int, genre entity.GenrePatchRequest) (entity.Genre, error)
	Delete(id int, force bool) (int64, error)
	DeleteConnection(id int, movies []int) (int64, error)
}

func (g *SQLiteGenreRepository) Create(genre *entity.Genre) (int64, error) {
	tx, err := g.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	query := `INSERT INTO genres (name) VALUES (?);`
	result, err := tx.Exec(query, genre.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if genre.MovieIds != nil {
		if err := CreateGenreConnection(tx, id, genre.MovieIds); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, nil
	}
	return id, nil
}
func (g *SQLiteGenreRepository) GetAll(moviesFlag bool) ([]entity.Genre, error) {
	query := `SELECT id, name FROM genres`
	rows, err := g.db.Query(query)
	if err != nil {
		return []entity.Genre{}, err
	}
	defer rows.Close()
	genres := []entity.Genre{}
	for rows.Next() {
		var id uint
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return []entity.Genre{}, err
		}

		if moviesFlag {
			movies, err := g.GetMovies(int(id))
			if err != nil {
				return []entity.Genre{}, err
			}
			genres = append(genres, entity.Genre{Id: id, Name: name, Movies: movies})
		} else {
			genres = append(genres, entity.Genre{Id: id, Name: name})
		}
	}
	return genres, nil
}
func (g *SQLiteGenreRepository) GetByID(id int) (entity.Genre, error) {
	query := `SELECT name FROM genres WHERE id = ?`
	row := g.db.QueryRow(query, id)
	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return entity.Genre{}, fmt.Errorf("there is no genre with this id: %v", id)
	} else if err != nil {
		return entity.Genre{}, err
	}
	movies, err := g.GetMovies(int(id))
	if err != nil {
		return entity.Genre{}, err
	}
	return entity.Genre{Id: uint(id), Name: name, Movies: movies}, nil
}
func (g *SQLiteGenreRepository) Update(id int, genre entity.GenrePatchRequest) (entity.Genre, error) {
	tx, err := g.db.Begin()
	if err != nil {
		return entity.Genre{}, err
	}
	defer tx.Rollback()
	query := `SELECT name FROM genres WHERE id =?`
	row := tx.QueryRow(query, id)
	var name string
	err = row.Scan(&name)
	if err == sql.ErrNoRows {
		return entity.Genre{}, fmt.Errorf("there is no genre with this id: %v", id)
	} else if err != nil {
		return entity.Genre{}, err
	}
	if genre.Name != nil {
		name = *genre.Name
	}
	newQuery := `UPDATE genres SET name = ? WHERE id = ?`
	result, err := tx.Exec(newQuery, name, id)
	if err != nil {
		return entity.Genre{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return entity.Genre{}, fmt.Errorf("there is no genre with this id: %v", id)
	}
	if genre.MovieIds != nil {
		err := CreateGenreConnection(tx, int64(id), genre.MovieIds)
		if err != nil {
			return entity.Genre{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return entity.Genre{}, err
	}
	return entity.Genre{Id: uint(id), Name: name}, nil
}
func (g *SQLiteGenreRepository) Delete(id int, force bool) (int64, error) {
	query := `SELECT name FROM genres WHERE id = ?`
	row := g.db.QueryRow(query, id)
	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("there is no genre with this id: %v", id)
	} else if err != nil {
		return 0, err
	}
	countFilms := 0
	queryCount := `SELECT COUNT(*) FROM movie_genres WHERE genre_id = ?`
	if err := g.db.QueryRow(queryCount, id).Scan(&countFilms); err != nil {
		return 0, err
	}
	if countFilms > 0 && !force {
		return 0, fmt.Errorf("You can't delete %s genre, because it's connected to %d films", name, id)
	}
	if force {
		queryDeleteConnection := `DELETE FROM movie_genres WHERE genre_id=?`
		_, err = g.db.Exec(queryDeleteConnection, id)
		if err != nil {
			return 0, err
		}
	}
	queryDelete := `DELETE FROM genres WHERE id = ?`
	result, err := g.db.Exec(queryDelete, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func (g *SQLiteGenreRepository) DeleteConnection(id int, movies []int) (int64, error) {
	str := make([]string, len(movies))
	args := []any{id}
	for i, movieID := range movies {
		str[i] = "?"
		args = append(args, movieID)
	}
	placeholder := strings.Join(str, ",")
	queryDeleteConnection := fmt.Sprintf(`DELETE FROM movie_genres WHERE genre_id=? AND movie_id IN (%s)`, placeholder)
	result, err := g.db.Exec(queryDeleteConnection, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// helper
func CreateGenreConnection(tx *sql.Tx, idGenre int64, idMovies []int) error {
	for _, id := range idMovies {
		_, err := tx.Exec(`INSERT OR IGNORE INTO movie_genres(movie_id, genre_id) VALUES (?,?)`, id, idGenre)
		if err != nil {
			return err
		}
	}
	return nil
}
func (g *SQLiteGenreRepository) GetMovies(id int) ([]entity.Movie, error) {
	query := `SELECT movies.id, movies.title, movies.release_year, movies.duration
	FROM movies
	JOIN movie_genres ON movies.id=movie_genres.movie_id
	WHERE movie_genres.genre_id = ?`
	rows, err := g.db.Query(query, id)
	if err != nil {
		return []entity.Movie{}, err
	}
	defer rows.Close()
	movies := []entity.Movie{}
	for rows.Next() {
		var id uint
		var title string
		var year int
		var duration float64
		err := rows.Scan(&id, &title, &year, &duration)
		if err != nil {
			return []entity.Movie{}, err
		}
		movies = append(movies, entity.Movie{Id: id, Title: title, ReleaseYear: year, Duration: duration})
	}
	return movies, nil
}

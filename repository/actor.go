package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
)

type SQLiteActorRepository struct {
	db *sql.DB
}

func NewSQLiteActorRepository(db *sql.DB) *SQLiteActorRepository {
	return &SQLiteActorRepository{db: db}
}

type ActorRepository interface {
	Create(actor *entity.Actor) (int64, error)
	GetAll() ([]entity.Actor, error)
	GetByID(id int) (entity.Actor, error)
	GetByName(name string) ([]entity.Actor, error)
	Update(id int, actor entity.ActorPatchRequest) (entity.Actor, error)
	Delete(id int, force bool) (int64, error)
}

func (a *SQLiteActorRepository) Create(actor *entity.Actor) (int64, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	query := `INSERT INTO actors (name, birthdate) VALUES (?, ?);`

	result, err := tx.Exec(query, actor.Name, actor.BirthDate.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if actor.MovieIds != nil {
		if err := CreateConnection(tx, id, actor.MovieIds); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}
func (a *SQLiteActorRepository) GetAll() ([]entity.Actor, error) {
	query := `SELECT id, name, birthdate FROM actors`
	rows, err := a.db.Query(query)
	if err != nil {
		return []entity.Actor{}, err
	}
	defer rows.Close()
	actors := []entity.Actor{}
	for rows.Next() {
		var id uint
		var name, birthdate string
		if err := rows.Scan(&id, &name, &birthdate); err != nil {
			return []entity.Actor{}, err
		}
		birthTime, err := time.Parse("2006-01-02", birthdate)
		if err != nil {
			return []entity.Actor{}, err
		}
		actors = append(actors, entity.Actor{Id: id, Name: name, BirthDate: birthTime})
	}
	return actors, nil
}
func (a *SQLiteActorRepository) GetByID(id int) (entity.Actor, error) {
	query := `SELECT name, birthdate FROM actors WHERE id = ?`
	row := a.db.QueryRow(query, id)
	var name, birthdate string
	err := row.Scan(&name, &birthdate)
	if err == sql.ErrNoRows {
		return entity.Actor{}, fmt.Errorf("there is no actor with this id: %v", id)
	} else if err != nil {
		return entity.Actor{}, err
	}
	birthTime, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return entity.Actor{}, err
	}
	return entity.Actor{Id: uint(id), Name: name, BirthDate: birthTime}, nil
}
func (a *SQLiteActorRepository) GetByName(name string) ([]entity.Actor, error) {
	query := `SELECT id, name, birthdate FROM actors WHERE name LIKE ?`
	searchPattern := "%" + name + "%"
	rows, err := a.db.Query(query, searchPattern)
	if err != nil {
		return []entity.Actor{}, err
	}
	defer rows.Close()
	actors := []entity.Actor{}
	for rows.Next() {
		var id uint
		var nameActual, birthdate string
		if err := rows.Scan(&id, &nameActual, &birthdate); err != nil {
			return []entity.Actor{}, err
		}
		birthTime, err := time.Parse("2006-01-02", birthdate)
		if err != nil {
			return []entity.Actor{}, err
		}
		actors = append(actors, entity.Actor{Id: uint(id), Name: nameActual, BirthDate: birthTime})
	}
	return actors, nil
}
func (a *SQLiteActorRepository) Update(id int, actor entity.ActorPatchRequest) (entity.Actor, error) {
	query := `SELECT name, birthdate FROM actors WHERE id = ?`
	row := a.db.QueryRow(query, id)
	var name, birthdate string
	err := row.Scan(&name, &birthdate)
	if err == sql.ErrNoRows {
		return entity.Actor{}, fmt.Errorf("there is no actor with this id: %v", id)
	} else if err != nil {
		return entity.Actor{}, err
	}
	if actor.BirthDate != nil {
		birthdate = *actor.BirthDate
	}
	if actor.Name != nil {
		name = *actor.Name
	}
	birthTime, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return entity.Actor{}, err
	}
	newQuery := `UPDATE actors SET name = ?, birthdate = ? WHERE id = ?`
	result, err := a.db.Exec(newQuery, name, birthTime.Format("2006-01-02"), id)
	if err != nil {
		return entity.Actor{}, err
	}
	//if someone delete this entity before this func updates everything
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return entity.Actor{}, fmt.Errorf("there is no actor with this id: %v", id)
	}
	return entity.Actor{Id: uint(id), Name: name, BirthDate: birthTime}, nil
}
func (a *SQLiteActorRepository) Delete(id int, force bool) (int64, error) {
	query := `SELECT name, birthdate FROM actors WHERE id = ?`
	row := a.db.QueryRow(query, id)
	var name, birthdate string
	err := row.Scan(&name, &birthdate)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("there is no actor with this id: %v", id)
	} else if err != nil {
		return 0, err
	}
	countFilms := 0
	queryCount := `SELECT COUNT(*) FROM movie_actors WHERE actor_id = ? `
	if err = a.db.QueryRow(queryCount, id).Scan(&countFilms); err != nil {
		return 0, err
	}
	if countFilms > 0 && !force {
		return 0, fmt.Errorf("You can't delete %s, because he/she plays in %d films", name, countFilms)
	}
	if force {
		queryDeleteConnection := `DELETE FROM movie_actors WHERE actor_id=?`
		_, err = a.db.Exec(queryDeleteConnection, id)
		if err != nil {
			return 0, err
		}
	}
	queryDelete := `DELETE FROM actors WHERE id = ?`
	result, err := a.db.Exec(queryDelete, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// helper
func CreateConnection(tx *sql.Tx, idActor int64, idMovies []int) error {
	for _, id := range idMovies {
		_, err := tx.Exec(`INSERT OR IGNORE INTO movie_actors(movie_id, actor_id) VALUES (?,?)`,
			id, idActor)
		if err != nil {
			return err
		}
	}
	return nil
}

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
	GetByName(name string) (entity.Actor, error)
}

func (a *SQLiteActorRepository) Create(actor *entity.Actor) (int64, error) {
	sql := `INSERT INTO actors (name, birthdate) VALUES (?, ?);`

	result, err := a.db.Exec(sql, actor.Name, actor.BirthDate.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (a *SQLiteActorRepository) GetAll() ([]entity.Actor, error) {
	sql := `SELECT id, name, birthdate FROM actors`
	rows, err := a.db.Query(sql)
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
	sql := `SELECT id, name, birthdate FROM actors`
	rows, err := a.db.Query(sql)
	if err != nil {
		return entity.Actor{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var idActor uint
		var name, birthdate string
		if err := rows.Scan(&idActor, &name, &birthdate); err != nil {
			return entity.Actor{}, err
		}
		if id == int(idActor) {
			birthTime, err := time.Parse("2006-01-02", birthdate)
			if err != nil {
				return entity.Actor{}, err
			}
			return entity.Actor{Id: idActor, Name: name, BirthDate: birthTime}, nil
		}
	}
	return entity.Actor{}, fmt.Errorf("there is no actor with this id: %v", id)
}
func (a *SQLiteActorRepository) GetByName(name string) (entity.Actor, error) {
	sql := `SELECT id, name, birthdate FROM actors`
	rows, err := a.db.Query(sql)
	if err != nil {
		return entity.Actor{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id uint
		var nameActor, birthdate string
		if err := rows.Scan(&id, &nameActor, &birthdate); err != nil {
			return entity.Actor{}, err
		}
		if name == nameActor {
			birthTime, err := time.Parse("2006-01-02", birthdate)
			if err != nil {
				return entity.Actor{}, err
			}
			return entity.Actor{Id: id, Name: nameActor, BirthDate: birthTime}, nil
		}
	}
	return entity.Actor{}, fmt.Errorf("there is no actor with this name: %v", name)
}

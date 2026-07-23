package db

import (
    "database/sql"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
)

func Init() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./my.db")
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    if _, err := CreateTables(db); err != nil {
        db.Close()
        return nil, err
    }

    fmt.Println("Connected to SQLite")
    return db, nil
}

func CreateTables(db *sql.DB) (sql.Result, error) {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		release_year INTEGER NOT NULL,
		duration REAL NOT NULL
	);

	CREATE TABLE IF NOT EXISTS actors (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		birthdate DATE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS genres (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS movie_actors (
		movie_id INTEGER NOT NULL,
		actor_id INTEGER NOT NULL,
		PRIMARY KEY(movie_id, actor_id),
		FOREIGN KEY(movie_id) REFERENCES movies(id),
		FOREIGN KEY(actor_id) REFERENCES actors(id)
	);

	CREATE TABLE IF NOT EXISTS movie_genres (
		movie_id INTEGER NOT NULL,
		genre_id INTEGER NOT NULL,
		PRIMARY KEY(movie_id, genre_id),
		FOREIGN KEY(movie_id) REFERENCES movies(id),
		FOREIGN KEY(genre_id) REFERENCES genres(id)
	);
	`

	return db.Exec(query)
}
package db

import (
	"database/sql"
	"log"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./gradr.db")
	if err != nil {
		log.Fatal(err)
	}

	createTablesSQL := `
CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    subject TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS assignment_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    correct INTEGER NOT NULL,
    total INTEGER NOT NULL,
    percentage REAL GENERATED ALWAYS AS (CAST(correct AS REAL) / total * 100),
    class_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    FOREIGN KEY(class_id) REFERENCES classes(id)
    FOREIGN KEY(type_id) REFERENCES types(id)
);

CREATE TABLE IF NOT EXISTS assignment_weights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    weight INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    FOREIGN KEY(class_id) REFERENCES classes(id)
    FOREIGN KEY(type_id) REFERENCES types(id)
);`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatal(err)
	}

	createAssignmentTypes := `
INSERT OR IGNORE INTO assignment_types (name)
  VALUES ("Test"), ("Quiz"), ("Homework");
  `
	_, err = db.Exec(createAssignmentTypes)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

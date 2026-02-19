package sqlite

import (
	"database/sql"

	"github.com/Deepsayan-Das/student-api-go/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	age INTEGER NOT NULL
	)`)

	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int, error) {
	stat, err := s.Db.Prepare(`INSERT INTO students (name,email,age)VALUES(?,?,?)`) //prevent sql injection
	if err != nil {
		return 0, err
	}
	defer stat.Close()

	res, err := stat.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastid), nil
}

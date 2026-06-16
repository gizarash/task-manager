package store

import (
	"database/sql"
	"fmt"

	"github.com/gizarash/task-manager/internal/model"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB `json:"-"`
}

func New(dsn string) (*Store, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: opening database %s from New: %w", dsn, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("store: pinging database %s from New: %w", dsn, err)
	}

	query := `
    create table if not exists todos(
        id integer primary key,
        title text,
        done boolean
    );
	`
	_, err = db.Exec(query)

	if err != nil {
		return nil, fmt.Errorf("store: unable to create table todos from New: %w", err)
	}

	var store Store
	store.db = db

	return &store, nil
}

func (s *Store) Add(title string) model.Todo {
	query := `
		INSERT INTO todos (title, done) VALUES (?, ?)
		RETURNING id, title, done;
	`
	var todo model.Todo

	err := s.db.QueryRow(query, title, false).Scan(&todo.Id, &todo.Title, &todo.Done)
	if err != nil {
		return model.Todo{}
	}

	return todo
}
func (s *Store) MarkDone(id int) bool {
	query := `
		update todos set done = ? 
		where id = ? and done = ?
	`
	res, err := s.db.Exec(query, true, id, false)
	if err != nil {
		return false
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func (s *Store) Delete(id int) bool {
	query := `
		delete from todos where id = ?
	`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return false
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func (s *Store) List() []model.Todo {
	var todos []model.Todo

	query := `
		select id, title, done from todos
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return todos
	}
	defer rows.Close()

	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Done)
		if err != nil {
			return todos
		}
		todos = append(todos, todo)
	}

	return todos
}
func (s *Store) Save() error {
	return nil
}
func (s *Store) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("store: closing db connection from Save: %w", err)
	}
	return nil
}

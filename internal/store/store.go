package store

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/gizarash/task-manager/internal/model"
)

type Store struct {
	CurrentId int          `json:"current_id"`
	Todos     []model.Todo `json:"todos"`
	filePath  string       `json:"-"`
}

func New(filePath string) (*Store, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("store: opening file %s from Load: %w", filePath, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("store: reading stats of file %s from Load: %w", filePath, err)
	}

	var store Store
	if stat.Size() == 0 {
		store.CurrentId = 1
		store.Todos = []model.Todo{}
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&store)
		if err != nil {
			return nil, fmt.Errorf("store: decoding file %s from Load: %w", filePath, err)
		}
	}
	store.filePath = filePath

	return &store, nil
}

func (s *Store) Add(title string) model.Todo {
	var newTodo = model.Todo{Id: s.CurrentId, Title: title, Done: false}
	s.Todos = append(s.Todos, newTodo)
	s.CurrentId++
	return newTodo
}
func (s *Store) MarkDone(id int) bool {
	isChanged := false
	for i, t := range s.Todos {
		if t.Id == id && !s.Todos[i].Done {
			s.Todos[i].Done = true
			isChanged = true
			break
		}
	}
	return isChanged
}
func (s *Store) Delete(id int) bool {
	isDeleted := false
	for i, t := range s.Todos {
		if t.Id == id {
			s.Todos = slices.Delete(s.Todos, i, i+1)
			isDeleted = true
			break
		}
	}
	return isDeleted
}
func (s *Store) List() []model.Todo {
	return s.Todos
}
func (s *Store) Save() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("store: opening file %s from Save: %w", s.filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s)
	if err != nil {
		return fmt.Errorf("store: encoding json to file %s from Save: %w", s.filePath, err)
	}
	return nil
}

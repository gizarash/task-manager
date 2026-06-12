package store

import (
	"path/filepath"
	"testing"
)

func NewTestStore(t *testing.T) *Store {
	t.Helper()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_store.json")
	s, err := New(tmpFile)
	if err != nil {
		t.Fatalf("unable to initialise store: %v", err)
	}

	return s
}
func TestAdd(t *testing.T) {
	s := NewTestStore(t)
	todo1 := s.Add("test todo 1")
	todos := s.List()

	t.Run("add new 1 record", func(t *testing.T) {
		if len(todos) != 1 {
			t.Errorf("expected 1 element, got: %d", len(todos))
		}
	})
	t.Run("default status is false", func(t *testing.T) {
		if todo1.Done != false {
			t.Errorf("expected Done = false, got: %t", todo1.Done)
		}
	})
	t.Run("increment id", func(t *testing.T) {
		todo2 := s.Add("test todo 2")
		if todo2.Id != todo1.Id+1 {
			t.Errorf("expected Id = 2, got: %d", todo2.Id)
		}
	})
}

func TestMarkDone(t *testing.T) {
	s := NewTestStore(t)
	todo := s.Add("test todo 1")

	t.Run("return true on success", func(t *testing.T) {
		res := s.MarkDone(todo.Id)
		if res != true {
			t.Errorf("expected result = true, got: %t", res)
		}
	})
	t.Run("mark as done", func(t *testing.T) {
		todos := s.List()
		if todos[0].Done != true {
			t.Errorf("expected Done = true, got: %t", todos[0].Done)
		}
	})
	t.Run("not found", func(t *testing.T) {
		res := s.MarkDone(999)
		if res != false {
			t.Errorf("expected result = false, got: %t", res)
		}
	})
	t.Run("already done", func(t *testing.T) {
		res := s.MarkDone(todo.Id)
		if res != false {
			t.Errorf("expected result = false, got: %t", res)
		}
	})
}

func TestDelete(t *testing.T) {
	s := NewTestStore(t)

	todo := s.Add("test todo 1")

	t.Run("return true on success", func(t *testing.T) {
		res := s.Delete(todo.Id)
		if res != true {
			t.Errorf("expected result = true, got: %t", res)
		}
	})
	t.Run("is deleted", func(t *testing.T) {
		if len(s.List()) != 0 {
			t.Errorf("expected 0 records, got: %d", len(s.List()))
		}
	})
	t.Run("not found", func(t *testing.T) {
		res := s.Delete(999)
		if res != false {
			t.Errorf("expected result = false, got: %t", res)
		}
	})
	t.Run("already deleted", func(t *testing.T) {
		res := s.Delete(todo.Id)
		if res != false {
			t.Errorf("expected result = false, got: %t", res)
		}
	})
}

func TestList(t *testing.T) {
	s := NewTestStore(t)
	s.Add("test todo 1")
	s.Add("test todo 2")
	res := s.List()
	if len(res) != 2 {
		t.Errorf("check if function return correct number of todos - expected result = 2, got: %d", len(res))
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/gizarash/task-manager/internal/model"
	"github.com/gizarash/task-manager/internal/store"
)

func NewTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	
	mux := http.NewServeMux()
	store, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("unable to initialise store: %v", err)
	}
	t.Cleanup(func() {
		store.Close()
	})
	th := &TaskHandler{
		store: store,
	}

	mux.HandleFunc("GET /tasks", th.handleListTasks)

	mux.HandleFunc("POST /tasks", th.handleCreateTask)

	mux.HandleFunc("PATCH /tasks/{id}/done", th.handleUpdateTask)

	mux.HandleFunc("DELETE /tasks/{id}", th.handleDeleteTask)

	ts := httptest.NewServer(mux)

	return ts
}

func TestList(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/tasks", ts.URL))
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	t.Run("return 200", func(t *testing.T) {
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected result = 200, got: %d", res.StatusCode)
		}
	})

	todos := []model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	if err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
	t.Run("return empty array", func(t *testing.T) {
		if len(todos) !=0 {
			t.Errorf("expected empty array, got: %+v", todos)
		}
	})
}

func TestAdd(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	todo := model.Todo{
		Title: "test",
	}
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("unable to marshal json: %v", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/tasks", ts.URL), "application/json", bytes.NewBuffer(todoBytes))
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()
	t.Run("return 201 on success", func(t *testing.T) {
		if res.StatusCode != http.StatusCreated {
			t.Errorf("expected result = 201, got: %d", res.StatusCode)
		}
	})
	testText := []byte("test text")
	res, err = http.Post(fmt.Sprintf("%s/tasks", ts.URL), "application/json", bytes.NewBuffer(testText))
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	t.Run("return 400 on invalid body", func(t *testing.T) {
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected result = 400, got: %d", res.StatusCode)
		}
	})
}

func TestDone(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	todo := model.Todo{
		Title: "test",
	}
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("unable to marshal json: %v", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/tasks", ts.URL), "application/json", bytes.NewBuffer(todoBytes))
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	var result model.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode json from add response: %v", err)
	}

	re := regexp.MustCompile(`\d+$`)
	numberStr := re.FindString(result.Message)
	id := 0
	
	if numberStr != "" {
		id, _ = strconv.Atoi(numberStr)
	} else {
		t.Fatalf("unable to parse id from response: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/tasks/%d/done", ts.URL, id), http.NoBody)
	if err != nil {
		t.Fatalf("unable to create request: %v", err)
	}
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	t.Run("return 200 on success", func(t *testing.T) {
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected result = 200, got: %d", res.StatusCode)
		}
	})

	req, err = http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/tasks/2/done", ts.URL), http.NoBody)
	if err != nil {
		t.Fatalf("unable to create request: %v", err)
	}

	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	t.Run("return 404 on invalid id", func(t *testing.T) {
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("expected result = 404, got: %d", res.StatusCode)
		}
	})
}

func TestDelete(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	todo := model.Todo{
		Title: "test",
	}
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("unable to marshal json: %v", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/tasks", ts.URL), "application/json", bytes.NewBuffer(todoBytes))
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	var result model.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode json from add response: %v", err)
	}

	re := regexp.MustCompile(`\d+$`)
	numberStr := re.FindString(result.Message)
	id := 0
	
	if numberStr != "" {
		id, _ = strconv.Atoi(numberStr)
	} else {
		t.Fatalf("unable to parse id from response: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/%d", ts.URL, id), http.NoBody)
	if err != nil {
		t.Fatalf("unable to create request: %v", err)
	}
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	t.Run("return 200 on success", func(t *testing.T) {
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected result = 200, got: %d", res.StatusCode)
		}
	})

	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/2", ts.URL), http.NoBody)
	if err != nil {
		t.Fatalf("unable to create request: %v", err)
	}

	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("unable to make request: %v", err)
	}
	defer res.Body.Close()

	t.Run("return 404 on invalid id", func(t *testing.T) {
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("expected result = 404, got: %d", res.StatusCode)
		}
	})
}
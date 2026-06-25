package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gizarash/task-manager/internal/model"
	"github.com/gizarash/task-manager/internal/store"
)

type TaskHandler struct {
	store *store.Store
}

func (h *TaskHandler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	todos := h.store.List()
	w.Header().Set("Content-Type", "application/json")
	data := todos
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todo := model.Todo{}
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
		errorData := model.Response{Message: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorData)
		return
	}

	result := h.store.Add(todo.Title)
	err = h.store.Save()
	if err != nil {
		slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
		errorData := model.Response{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorData)
	} else {
		data := model.Response{Message: fmt.Sprintf("record '%s' successfully added with id = %d", result.Title, result.Id)}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	}
}

func (h *TaskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
		errorData := model.Response{Message: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorData)
		return
	}

	result := h.store.MarkDone(id)

	if result {
		err = h.store.Save()
		if err != nil {
			slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
			errorData := model.Response{Message: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorData)
		} else {
			data := model.Response{Message: fmt.Sprintf("record with id = %d successfully marked as done", id)}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data)
		}
	} else {
		slog.Warn(fmt.Sprintf("%s %s - record not found or have already marked as done", r.Method, r.URL))
		errorData := model.Response{Message: "record not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorData)
	}
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
		errorData := model.Response{Message: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorData)
		return
	}

	result := h.store.Delete(id)

	if result {
		err = h.store.Save()
		if err != nil {
			slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
			errorData := model.Response{Message: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorData)
		} else {
			data := model.Response{Message: fmt.Sprintf("record with id = %d successfully deleted", id)}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data)
		}
	} else {
		slog.Warn(fmt.Sprintf("%s %s - record not found", r.Method, r.URL))
		errorData := model.Response{Message: "record not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorData)
	}
}
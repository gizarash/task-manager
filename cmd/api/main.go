package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gizarash/task-manager/internal/model"
	"github.com/gizarash/task-manager/internal/store"
)

func main() {
	mux := http.NewServeMux()
	store, err := store.New("store.json")
	if err != nil {
		slog.Error("unable to initialise store", "error", err)
		os.Exit(1)
	}

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		todos := store.List()
		slog.Info(fmt.Sprintf("%s %s - returned %d records", r.Method, r.URL, len(todos)))
		w.Header().Set("Content-Type", "application/json")
		data := todos
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
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

		result := store.Add(todo.Title)
		err = store.Save()
		if err != nil {
			slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
			errorData := model.Response{Message: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorData)
		} else {
			slog.Info(fmt.Sprintf("%s %s - record '%s' successfully added with id = %d", r.Method, r.URL, result.Title, result.Id))
			data := model.Response{Message: fmt.Sprintf("record '%s' successfully added with id = %d", result.Title, result.Id)}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(data)
		}
	})

	mux.HandleFunc("PATCH /tasks/{id}/done", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
			errorData := model.Response{Message: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorData)
			return
		}

		result := store.MarkDone(id)

		if result {
			err = store.Save()
			if err != nil {
				slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
				errorData := model.Response{Message: err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorData)
			} else {
				slog.Info(fmt.Sprintf("%s %s - record with id = %d successfully marked as done", r.Method, r.URL, id))
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
	})

	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
			errorData := model.Response{Message: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorData)
			return
		}

		result := store.Delete(id)

		if result {
			err = store.Save()
			if err != nil {
				slog.Error(fmt.Sprintf("%s %s - got error: '%s'", r.Method, r.URL, err.Error()))
				errorData := model.Response{Message: err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorData)
			} else {
				slog.Info(fmt.Sprintf("%s %s - record with id = %d successfully deleted", r.Method, r.URL, id))
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
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting server at port 8080...")
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				todos := store.List()
				done := 0
				pending := 0

				for _, todo := range todos {
					if todo.Done {
						done++
					} else {
						pending++
					}
				}

				slog.Info(fmt.Sprintf("stats: total=%d, done=%d, pending=%d", len(todos), done, pending))
			}
		}
	}(ctx)

	<-ctx.Done()
	slog.Info("server is shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("server forced to shutdown: %s", err))
	}

	slog.Info("server stopped")
}

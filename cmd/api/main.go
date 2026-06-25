package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gizarash/task-manager/internal/store"
)

func main() {
	mux := http.NewServeMux()
	store, err := store.New("store.db")
	if err != nil {
		slog.Error("unable to initialise store", "error", err)
		os.Exit(1)
	}
	th := &TaskHandler{
		store: store,
	}

	mux.HandleFunc("GET /tasks", th.handleListTasks)

	mux.HandleFunc("POST /tasks", th.handleCreateTask)

	mux.HandleFunc("PATCH /tasks/{id}/done", th.handleUpdateTask)

	mux.HandleFunc("DELETE /tasks/{id}", th.handleDeleteTask)

	handler := LoggingMiddleware(mux)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
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

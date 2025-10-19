package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"spreadlove/db"
)

//go:embed schema.sql
var schemaSQL string

type App struct {
	db      *sql.DB
	queries *db.Queries
}

func main() {
	app := &App{}

	if err := app.setupDB(); err != nil {
		log.Fatal("Failed to setup database:", err)
	}

	defer app.db.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bye"))
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Block
	<-quit
	log.Println("Shutdown Server ...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}

func (app *App) setupDB() error {
	database, err := sql.Open("sqlite3", "./love.db")
	if err != nil {
		return err
	}

	app.db = database
	app.queries = db.New(database)

	// Setup DB from `schema.sql`
	_, err = database.Exec(schemaSQL)
	return err
}

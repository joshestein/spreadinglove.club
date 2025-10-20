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
	"github.com/go-chi/render"

	"spreadlove/db"
)

//go:embed schema.sql
var schemaSQL string

type App struct {
	db      *sql.DB
	queries *db.Queries
}

type MessageResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
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

	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	r.Get("/submit", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/submit.html")
	})

	r.Post("/submit", app.handleSubmitMessage)

	r.Route("/api", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/message", app.handleGetRandomMessage)
	})

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	server := &http.Server{
		Addr:    ":" + port,
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

func (app *App) handleGetRandomMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	msg, err := app.queries.GetRandomMessage(ctx)

	if err != nil {
		log.Printf("Error fetching random message: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		ID:      msg.ID,
		Content: msg.Content,
	}

	if msg.CreatedAt.Valid {
		response.CreatedAt = msg.CreatedAt.Time
	}

	render.JSON(w, r, response)
}

func (app *App) handleSubmitMessage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	_, err := app.queries.CreatePendingMessage(ctx, content)

	if err != nil {
		log.Printf("Error submitting message message: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

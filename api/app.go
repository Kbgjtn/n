package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Kbgjtn/notethingness-api.git/db"
)

type Config struct {
	port    string
	host    string
	url     string
	connStr string
}

type Server struct {
	router http.Handler
	db     *sql.DB
	config *Config
}

func NewServer() *Server {
	config := &Config{
		port:    ":" + os.Getenv("PORT"),
		host:    os.Getenv("HOST"),
		url:     fmt.Sprintf("http://%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		connStr: os.Getenv("DB_URL"),
	}

	db := db.NewDatabase("file://db/migration")
	slog.Info("[ ☘️ Connect to DB POSTGRES ]")
	dbConnection, err := db.Connect(config.connStr)
	if err != nil {
		panic(err)
	}

	server := &Server{
		db:     dbConnection,
		config: config,
	}

	slog.Info("[ ☘️ Run migration rollback ]")
	if err := db.RollbackMigration(config.connStr); err != nil {
		panic(err)
	}

	slog.Info("[ ☘️ Run migration ]")
	if err := db.RunMigration(config.connStr); err != nil {
		panic(err)
	}

	server.Routes()
	return server
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.config.port,
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("[ Server started on port: " + s.config.port + " ]")
	defer func() {
		if err := s.db.Close(); err != nil {
			slog.Error("[ Failed to close DB connection ]" + "\nError: " + err.Error())
		}
	}()

	// Using a buffered channel to avoid goroutine leaks
	channel := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			channel <- error(err)
		}
		close(channel)
	}()

	select {
	case err := <-channel:
		slog.Error("[ Failed to start server ]" + "\nError: " + err.Error())
		return err
	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			return error(err)
		}
		return nil
	}
}

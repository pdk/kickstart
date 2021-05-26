package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Server provides webservices.
type Server struct {
	DB *sql.DB
}

// New returns a new Server.
func New(db *sql.DB) Server {
	return Server{
		DB: db,
	}
}

// ListenAndServe does a listen. Returns a function that can be called to
// shutdown the server.
func (s Server) ListenAndServe() func() {

	router := s.makeRouter()

	httpServer := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// make the shutdown function for clean shutdowns.
	shutdown := func() {
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(timeout); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}

	// launch the signal catcher
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Printf("ListenAndServe received interrupt signal. Shutting down...")
		shutdown()
		log.Printf("ListenAndServe shutdown complete")

		os.Exit(1) // indicate server interrupted
	}()

	// launch the http listener
	go func() {
		log.Printf("listening at %s", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// Completed serving
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return shutdown
}

package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"app/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()

	router.Handle("/*", public())
	routes.Register(router)

	port := os.Getenv("LISTEN_PORT")
	slog.Info("HTTP server started", "port", port)
	http.ListenAndServe(port, router)
}

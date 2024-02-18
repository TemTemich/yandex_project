package app

import (
	"context"
	"log"
	"net/http"
	"orchestrator/internal/core/adapters/calculator"
	"orchestrator/internal/core/adapters/storage"
	handlers "orchestrator/internal/core/handlers/http"
	"orchestrator/internal/core/services"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Run() error {

	storage := storage.NewStorage()

	calcHost := os.Getenv("DAEMON_URL")
	client := &http.Client{}
	calc := calculator.NewCalculator(
		calcHost,
		client,
	)
	service := services.NewService(
		storage,
		calc,
	)
	go service.Run()

	h := handlers.NewHandlers(service)
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Route("/operation", func(r chi.Router) {
		r.Get("/all", h.GetOperations)
		r.Post("/update", h.UpdateOperation)
	})

	router.Route("/expressions", func(r chi.Router) {
		r.Post("/", h.AddExpression)
		r.Get("/", h.GetExpressions)
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			r = r.WithContext(context.WithValue(r.Context(), "id", id))
			h.GetExpressionByID(w, r)
		})
	})

	log.Println("server is running on port 8080")
	return http.ListenAndServe(":8080", router)
}

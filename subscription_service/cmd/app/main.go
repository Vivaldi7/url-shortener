package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/config"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/db"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/handlers"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/repository"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/service"
	"github.com/vivaldi7/golang_code/subscription_service/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)

	// Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		logger.Log.Error("Failed to connect to database", "error", err)
	}

	defer database.Close()

	// Run migrations
	//	runMigration(database)

	// Initialize repository, service, handler
	repo := repository.NewSubscriptionRepositoty(database)
	subscriptionService := service.NewSubscriptionService(repo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTOINS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	// Routes
	r.Route("api/v1/subscriptions", func(r chi.Router) {
		r.Post("/", subscriptionHandler.Create)
		r.Get("/", subscriptionHandler.List)
		r.Post("/", subscriptionHandler.GetTotalCost)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", subscriptionHandler.Get)
			r.Get("/", subscriptionHandler.Update)
			r.Get("/", subscriptionHandler.Delete)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	serverAddr := ":" + cfg.ServerPort
	logger.Log.Info("Starting server", "address", serverAddr)

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		logger.Log.Error("Failed to start server", "error", err)
		panic(err)
	}

}

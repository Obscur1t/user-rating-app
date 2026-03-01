package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rating/internal/db"
	"rating/internal/handler"
	"rating/internal/logger"
	"rating/internal/middleware"
	"rating/internal/repo/postgres"
	"rating/internal/service"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	dbUrl := os.Getenv("DB_URL")
	addr := os.Getenv("SERVER_ADDR")
	logLevel := os.Getenv("LOG_LEVEL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	if addr == "" {
		addr = ":8080"
	}

	if logLevel == "" {
		log.Fatal("LOG_LEVEL environment variable is not set")
	}

	pool, err := db.NewDb(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	logger := logger.SetupLogger(logLevel)

	userRepo := postgres.NewUserRepo(pool)
	userService := service.NewUserService(userRepo)
	userHandlers := handler.NewUserHandler(userService, logger)

	mux := http.NewServeMux()
	chainedHandler := middleware.Chain(
		mux,
		middleware.RecoveryMiddleware(logger),
		middleware.LoggerMiddleware(logger),
	)

	mux.HandleFunc("POST /users", userHandlers.CreateUserHandler)
	mux.HandleFunc("GET /users", userHandlers.GetUsers)
	mux.HandleFunc("GET /users/{nickname}", userHandlers.GetUser)
	mux.HandleFunc("PATCH /users/{nickname}", userHandlers.ChangeData)
	mux.HandleFunc("DELETE /users/{nickname}", userHandlers.Delete)

	server := &http.Server{
		Addr:         addr,
		Handler:      chainedHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	defer pool.Close()
}

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
	"rating/internal/repo/postgres"
	"rating/internal/service"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	dbUrl := os.Getenv("DB_URL")
	addr := os.Getenv("SERVER_ADDR")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	if addr == "" {
		addr = ":8080"
	}

	pool, err := db.NewDb(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	userRepo := postgres.NewUserRepo(pool)
	userService := service.NewUserService(userRepo)
	userHandlers := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", userHandlers.CreateUserHandler)
	mux.HandleFunc("GET /users", userHandlers.GetUsers)
	mux.HandleFunc("GET /users/{nickname}", userHandlers.GetUser)
	mux.HandleFunc("PATCH /users/{nickname}", userHandlers.ChangeData)
	mux.HandleFunc("DELETE /users/{nickname}", userHandlers.Delete)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
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

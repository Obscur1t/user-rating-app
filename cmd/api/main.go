package main

import (
	"context"
	"log"
	"net/http"
	"rating/internal/db"
	"rating/internal/handler"
	usersrepo "rating/internal/repo/users_repo"
	"rating/internal/service"
	"time"
)

//Table:
// id BIGSERIAL PRIMARY KEY
// name TEXT NOT NULL
// nickname TEXT NOT NULL UNIQUE
// likes INT NOT NULL DEFAULT 0
// viewers INT NOT NULL DEFAULT 0
// rating NUMERIC GENERATED ALWAYS AS (
//     CASE WHEN viewers > 0
//          THEN likes::NUMERIC / viewers
//          ELSE 0
//     END
// ) STORED

func main() {
	pool, err := db.NewDb(context.Background())
	if err != nil {
		log.Fatalf("failed to create pool %v", err)
	}

	userRepo := usersrepo.NewUserRepo(pool)
	userService := service.NewUserService(userRepo)
	userHandlers := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", userHandlers.CreateUserHandler)
	mux.HandleFunc("GET users", userHandlers.GetUsers)
	mux.HandleFunc("GET users/{nickname}", userHandlers.GetUser)
	mux.HandleFunc("PATCH users/{nickname}", userHandlers.ChangeData)
	mux.HandleFunc("DELETE users/{nickname}", userHandlers.Delete)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nirupam52/expenseTrack/internal/config"
	"github.com/nirupam52/expenseTrack/internal/db"
	"github.com/nirupam52/expenseTrack/internal/handlers"
	"github.com/nirupam52/expenseTrack/internal/repository"
	"github.com/nirupam52/expenseTrack/internal/response"
)

func main() {
	appConfig := config.LoadConfig()

	database, err := db.OpenDB(appConfig.DBPath)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer database.Close()

	if err := db.InitDB(context.Background(), database); err != nil {
		log.Fatalf("could not initialise database: %v", err)
	}

	log.Printf("database ready at %s", appConfig.DBPath)

	userRepo := repository.NewUserRepository(database)
	expenseRepo := repository.NewExpenseRepository(database)
	authRepo := repository.NewAuthRepository(database)

	userHandler := handlers.NewUserHandler(userRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo)
	authHandler := handlers.NewAuthHandler(authRepo)

	protect := handlers.NewAuthMiddleware(authRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		if err := response.WriteSuccess(w, http.StatusOK, "pong"); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})

	authHandler.RegisterRoutes(mux, protect)
	userHandler.RegisterRoutes(mux, protect)
	expenseHandler.RegisterRoutes(mux, protect)

	addr := fmt.Sprintf(":%s", appConfig.Port)
	log.Printf("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

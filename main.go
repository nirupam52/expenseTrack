package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nirupam52/expenseTrack/internal/config"
	"github.com/nirupam52/expenseTrack/internal/db"
)

func main() {
	appConfig := config.LoadConfig()

	database, err := db.OpenDB(appConfig.DbConfig.DBPath)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer database.Close()

	// Run schema migrations and enable foreign keys.
	if err := db.InitDB(context.Background(), database); err != nil {
		log.Fatalf("could not initialise database: %v", err)
	}

	log.Printf("database ready at %s", appConfig.DbConfig.DBPath)

	mux := http.NewServeMux()

	// Health-check — useful to confirm the server is running.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
		log.Print("request received new")
	})

	addr := fmt.Sprintf(":%s", appConfig.Port)
	log.Printf("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

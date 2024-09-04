package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Todo-list/internal/app/api"
	"Todo-list/internal/app/config"
	httpHandler "Todo-list/internal/app/http"
	"Todo-list/internal/app/repository/database"
	"Todo-list/internal/app/repository/postgresql"
	"github.com/gorilla/mux"
)

func main() {

	// Create config
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new connection pool to database
	pool, err := database.NewPool(cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

	// Create a new Database with connection pool
	db := database.NewDatabase(pool)

	// Create a new repo with Database
	repo := postgresql.New(*db)

	// Create a new service
	taskService := api.New(repo)

	// Create Http handler
	handler := httpHandler.New(taskService)

	// Init Router
	r := mux.NewRouter()

	handler.RegisterRoutes(r)

	// Start HTTP server
	if err = http.ListenAndServe(cfg.HttpPort, r); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

}

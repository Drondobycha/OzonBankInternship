package main

import (
	"fmt"
	"log"
	"my_project/inmemory"
	postgres "my_project/postgreSQL"
	"net/http"
	"os"
)

func main() {
	storageMode := os.Getenv("STORAGE_MODE")

	var store Storage
	switch storageMode {
	case "inmemory":
		store = inmemory.NewUrlStore()
	case "postgres":
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			log.Fatal("DATABASE_URL is not set")
		}
		var err error
		store, err = postgres.NewURLStore(connStr)
		if err != nil {
			log.Fatalf("Failed to create PostgreSQL store: %v", err)
		}
	default:
		log.Fatalf("Unknown storage mode: %s", storageMode)
	}

	// Используем store в вашем сервисе
	fmt.Println("Storage mode:", storageMode)
	// Запуск сервера
	http.HandleFunc("/shorten", HandleShorten(store))

	http.HandleFunc("/redirect/", HandleRedirect(store))

	fmt.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

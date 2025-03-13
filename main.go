package main

import (
	"flag"
	"fmt"
	"log"
	"my_project/inmemory"
	postgres "my_project/postgreSQL"
	"net/http"
)

func main() {
	storageType := flag.String("storage", "memory", "Тип хранилища (memory или postgres)")
	postgresConnStr := flag.String("postgres-conn", "postgres://OzonBankTest:12345678@localhost5432/url_shortener?sslmode=disable", "Строка подключения к PostgreSQL")
	flag.Parse()

	var store Storage
	var err error

	switch *storageType {
	case "memory":
		store = inmemory.NewUrlStore() // Использование пакета memory
	case "postgres":
		store, err = postgres.NewURLStore(*postgresConnStr)
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
		}
	default:
		log.Fatalf("Unknown storage type: %s", *storageType)
	}
	// Запуск сервера
	http.HandleFunc("/shorten", HandleShorten(store))

	http.HandleFunc("/redirect/", HandleRedirect(store))

	fmt.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

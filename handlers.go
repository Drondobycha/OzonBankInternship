package main

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

func HandleShorten(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверка на пустой URL
		if req.URL == "" {
			http.Error(w, "URL cannot be empty", http.StatusBadRequest)
			return
		}

		shortURL, created, err := store.Save(req.URL)
		if err != nil {
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		status := http.StatusCreated
		if !created {
			status = http.StatusOK
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(Response{ShortURL: shortURL})
	}
}

func HandleRedirect(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем короткий URL из пути запроса
		shortURL := r.URL.Path[len("/redirect/"):]

		// Получаем оригинальный URL по короткому URL
		originalURL, exists, err := store.Get(shortURL)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Перенаправляем на оригинальный URL
		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

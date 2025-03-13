package postgres

import (
	"context"
	"testing"
)

func TestPostgresSaveAndGet(t *testing.T) {
	connStr := "postgres://test:test@localhost:5432/test?sslmode=disable"
	store, err := NewURLStore(connStr)
	if err != nil {
		t.Fatalf("Failed to create URLStore: %v", err)
	}

	// Очистка таблицы перед тестом
	_, err = store.pool.Exec(context.Background(), "DELETE FROM urls")
	if err != nil {
		t.Fatalf("Failed to clean up table: %v", err)
	}

	// Тест сохранения URL
	shortURL, created, err := store.Save("https://www.example.com")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if !created {
		t.Error("Expected URL to be created")
	}

	// Проверка, что URL можно получить
	originalURL, exists, err := store.Get(shortURL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected URL to exist")
	}
	if originalURL != "https://www.example.com" {
		t.Errorf("Expected original URL to be 'https://www.example.com', got %v", originalURL)
	}
}

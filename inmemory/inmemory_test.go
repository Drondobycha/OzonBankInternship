package inmemory

import (
	"testing"
)

func TestInMemorySaveAndGet(t *testing.T) {
	store := NewUrlStore()

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

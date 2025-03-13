package main

import (
	"testing"
)

func TestMockStorage(t *testing.T) {
	mockStore := &MockStorage{
		SaveFunc: func(url string) (string, bool, error) {
			return "abc123", true, nil
		},
		GetFunc: func(shortURL string) (string, bool, error) {
			return "https://www.example.com", true, nil
		},
	}

	shortURL, created, err := mockStore.Save("https://www.example.com")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if !created {
		t.Error("Expected URL to be created")
	}
	if shortURL != "abc123" {
		t.Errorf("Expected short URL to be 'abc123', got %v", shortURL)
	}

	originalURL, exists, err := mockStore.Get("abc123")
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

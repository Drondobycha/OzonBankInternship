package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	SaveFunc func(url string) (string, bool, error)
	GetFunc  func(shortURL string) (string, bool, error)
}

func (m *MockStorage) Save(url string) (string, bool, error) {
	return m.SaveFunc(url)
}

func (m *MockStorage) Get(shortURL string) (string, bool, error) {
	return m.GetFunc(shortURL)
}

func TestHandleShorten(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    Request
		mockSaveFunc   func(url string) (string, bool, error)
		expectedStatus int
		expectedBody   Response
	}{
		{
			name: "successful shorten",
			requestBody: Request{
				URL: "https://www.example.com",
			},
			mockSaveFunc: func(url string) (string, bool, error) {
				return "abc123", true, nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody: Response{
				ShortURL: "abc123",
			},
		},
		{
			name: "invalid request body",
			requestBody: Request{
				URL: "",
			},
			mockSaveFunc: func(url string) (string, bool, error) {
				return "", false, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   Response{},
		},
		{
			name: "database error",
			requestBody: Request{
				URL: "https://www.example.com",
			},
			mockSaveFunc: func(url string) (string, bool, error) {
				return "", false, fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStorage{
				SaveFunc: tt.mockSaveFunc,
			}

			reqBodyBytes, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBodyBytes))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := HandleShorten(mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated || tt.expectedStatus == http.StatusOK {
				var resp Response
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.ShortURL != tt.expectedBody.ShortURL {
					t.Errorf("Handler returned unexpected body: got %v want %v", resp, tt.expectedBody)
				}
			}
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	tests := []struct {
		name           string
		shortURL       string
		mockGetFunc    func(shortURL string) (string, bool, error)
		expectedStatus int
		expectedURL    string
	}{
		{
			name:     "successful redirect",
			shortURL: "abc123",
			mockGetFunc: func(shortURL string) (string, bool, error) {
				return "https://www.example.com", true, nil
			},
			expectedStatus: http.StatusFound,
			expectedURL:    "https://www.example.com",
		},
		{
			name:     "short URL not found",
			shortURL: "unknown",
			mockGetFunc: func(shortURL string) (string, bool, error) {
				return "", false, nil
			},
			expectedStatus: http.StatusNotFound,
			expectedURL:    "",
		},
		{
			name:     "database error",
			shortURL: "abc123",
			mockGetFunc: func(shortURL string) (string, bool, error) {
				return "", false, fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStorage{
				GetFunc: tt.mockGetFunc,
			}

			req, err := http.NewRequest("GET", "/redirect/"+tt.shortURL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := HandleRedirect(mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusFound {
				if location := rr.Header().Get("Location"); location != tt.expectedURL {
					t.Errorf("Handler returned unexpected Location header: got %v want %v", location, tt.expectedURL)
				}
			}
		})
	}
}

package inmemory

import (
	"math/rand"
	"sync"
	"time"
)

type URLStore struct {
	mu          sync.RWMutex
	shortToLong map[string]string
	LongToshort map[string]string
}

func NewUrlStore() *URLStore {
	return &URLStore{
		shortToLong: make(map[string]string),
		LongToshort: make(map[string]string),
	}
}

func (s *URLStore) Save(originalURL string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if shortURL, exist := s.LongToshort[originalURL]; exist {
		return shortURL, false, nil
	}

	shortURL := generateShortURL()

	s.shortToLong[shortURL] = originalURL
	s.LongToshort[originalURL] = shortURL

	return shortURL, true, nil
}

func (s *URLStore) Get(shortUrl string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	originalURL, exist := s.shortToLong[shortUrl]
	return originalURL, exist, nil
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShortURL() string {
	const (
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
		keyLen  = 10
	)

	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

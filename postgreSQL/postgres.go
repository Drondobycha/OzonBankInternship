package postgres

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLStore struct {
	pool *pgxpool.Pool
}

func NewURLStore(connStr string) (*URLStore, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			short_url TEXT PRIMARY KEY,
			original_url TEXT UNIQUE
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return &URLStore{pool: pool}, nil
}

func (s *URLStore) Save(originalURL string) (string, bool, error) {
	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var shortURL string
	err = tx.QueryRow(ctx, `
		SELECT short_url FROM urls WHERE original_url = $1;
	`, originalURL).Scan(&shortURL)

	if err == nil {
		return shortURL, false, nil
	} else if err.Error() != "no rows in result set" {
		log.Printf("Database query error: %v", err)
		return "", false, fmt.Errorf("failed to query database: %w", err)
	}

	for {
		shortURL = generateShortURL()
		var exists bool
		err := tx.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = $1);
		`, shortURL).Scan(&exists)
		if err != nil {
			log.Printf("Database uniqueness check error: %v", err)
			return "", false, fmt.Errorf("failed to check short URL uniqueness: %w", err)
		}
		if !exists {
			break
		}
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO urls (short_url, original_url) VALUES ($1, $2);
	`, shortURL, originalURL)
	if err != nil {
		log.Printf("Database insert error: %v", err)
		return "", false, fmt.Errorf("failed to insert into database: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Transaction commit error: %v", err)
		return "", false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return shortURL, true, nil
}

func (s *URLStore) Get(shortURL string) (string, bool, error) {
	var originalURL string
	err := s.pool.QueryRow(context.Background(), `
		SELECT original_url FROM urls WHERE short_url = $1;
	`, shortURL).Scan(&originalURL)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", false, nil
		}
		log.Printf("Database query error: %v", err)
		return "", false, fmt.Errorf("failed to query database: %w", err)
	}

	return originalURL, true, nil
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	const keyLen = 10
	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

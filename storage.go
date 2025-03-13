package main

type Storage interface {
	Save(originalURL string) (shortURL string, created bool, err error)
	Get(shortURL string) (originalURL string, exists bool, err error)
}

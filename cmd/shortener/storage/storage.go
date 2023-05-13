package storage

import "github.com/google/uuid"

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) string
}

type LocalStorage struct {
	shortURL map[string]string
}

func (s *LocalStorage) PutURL(body string) string {
	id := uuid.New().String()
	s.shortURL[id] = body
	return id
}

func (s *LocalStorage) GetURL(id string) string {
	return s.shortURL[id]
}

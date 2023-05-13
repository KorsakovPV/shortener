package storage

import "github.com/google/uuid"

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) string
}

type localStorageStruct struct {
	shortURL map[string]string
}

func (s *localStorageStruct) PutURL(body string) string {
	id := uuid.New().String()
	s.shortURL[id] = body
	return id
}

func (s *localStorageStruct) GetURL(id string) string {
	return s.shortURL[id]
}

var LocalStorage AbstractStorage = &localStorageStruct{
	shortURL: map[string]string{},
}

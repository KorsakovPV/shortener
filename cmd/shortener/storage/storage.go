package storage

import "github.com/KorsakovPV/shortener/cmd/shortener/storage/localStorage"

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) string
}

var LocalStorage AbstractStorage = &localStorage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

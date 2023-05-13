package storage

import "github.com/KorsakovPV/shortener/cmd/shortener/storage/localstorage"

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) (string, error)
}

var LocalStorage AbstractStorage = &localstorage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

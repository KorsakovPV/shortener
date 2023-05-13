package storage

import "github.com/KorsakovPV/shortener/cmd/shortener/storage/local_storage"

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) string
}

var LocalStorage AbstractStorage = &local_storage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

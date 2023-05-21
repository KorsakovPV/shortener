package storage

import "github.com/KorsakovPV/shortener/cmd/shortener/storage/localstorage"

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) (string, error)
}

var localStorage AbstractStorage = &localstorage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

func GetStorage() AbstractStorage {
	return localStorage
}

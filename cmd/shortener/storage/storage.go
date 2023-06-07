package storage

import "github.com/KorsakovPV/shortener/cmd/shortener/storage/localstorage"

type AbstractStorage interface {
	PutURL(string) (string, error)
	GetURL(string) (string, error)
	LoadBackupURL() error
}

var localStorage AbstractStorage = &localstorage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

func GetStorage() AbstractStorage {
	return localStorage
}

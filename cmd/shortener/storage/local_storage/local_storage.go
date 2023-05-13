package local_storage

import "github.com/google/uuid"

type LocalStorageStruct struct {
	ShortURL map[string]string
}

func (s *LocalStorageStruct) PutURL(body string) string {
	id := uuid.New().String()
	s.ShortURL[id] = body
	return id
}

func (s *LocalStorageStruct) GetURL(id string) string {
	return s.ShortURL[id]
}

//var LocalStorage AbstractStorage = &localStorageStruct{
//	shortURL: map[string]string{},
//}

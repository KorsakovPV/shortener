package localstorage

import (
	"fmt"
	"github.com/google/uuid"
)

type LocalStorageStruct struct {
	ShortURL map[string]string
}

func (s *LocalStorageStruct) PutURL(body string) string {
	id := uuid.New().String()
	s.ShortURL[id] = body
	return id
}

func (s *LocalStorageStruct) GetURL(id string) (string, error) {
	url, ok := s.ShortURL[id]
	if !ok {
		return url, fmt.Errorf("id %s not found", id)
	} else {
		return url, nil
	}
}

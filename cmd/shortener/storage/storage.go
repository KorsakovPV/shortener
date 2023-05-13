package storage

type AbstractStorage interface {
	PutURL(string) string
	GetURL(string) string
}

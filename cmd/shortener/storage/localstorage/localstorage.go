package localstorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/internal/models"
)

type ShortURL struct {
	UUID        string `json:"uuid"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event ShortURL) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

// Consumer
type Consumer struct {
	file   *os.File
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadShortURL() ([]*ShortURL, error) {
	events := []*ShortURL{}
	for {
		data, _, err := c.reader.ReadLine()
		if err == io.EOF {
			return events, nil
		}
		if err != nil {
			return nil, err
		}
		url := &ShortURL{}
		err = json.Unmarshal(data, url)
		if err != nil {
			return nil, err
		}
		events = append(events, url)
	}
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

type LocalStorageStruct struct {
	ShortURL map[string]string
}

func (s *LocalStorageStruct) PutURL(id string, body string) (string, error) {

	cfg := config.GetConfig()

	if cfg.FlagFileStoragePath != "" {
		Produc, err := NewProducer(cfg.FlagFileStoragePath)
		if err != nil {
			return "", err
		}
		defer Produc.Close()

		url := ShortURL{UUID: id, OriginalURL: body}

		_, err = json.Marshal(&url)
		if err != nil {
			return "", err
		}

		if err := Produc.WriteEvent(url); err != nil {
			return "", err
		}
	}

	s.ShortURL[id] = body
	return id, nil
}

func (s *LocalStorageStruct) PutURLBatch(body []models.RequestBatch) ([]models.ResponseButch, error) {
	cfg := config.GetConfig()

	bodyResponseButch := make([]models.ResponseButch, len(body))
	urls := make([]ShortURL, len(body))

	for i := 0; i < len(body); i++ {
		id := body[i].UUID
		url := body[i].URL

		_, ok := s.ShortURL[id]
		if ok {
			return nil, fmt.Errorf("id %s is exist", id)
		}
		bodyResponseButch[i].UUID = id
		bodyResponseButch[i].URL = fmt.Sprintf("%s/%s", cfg.FlagBaseURLAddr, id)
		urls[i].UUID = id
		urls[i].OriginalURL = url
	}

	if cfg.FlagFileStoragePath != "" {
		Produc, err := NewProducer(cfg.FlagFileStoragePath)
		if err != nil {
			return nil, err
		}
		defer Produc.Close()

		_, err = json.Marshal(&urls)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(body); i++ {

			if err := Produc.WriteEvent(urls[i]); err != nil {
				return nil, err
			}
			s.ShortURL[urls[i].UUID] = urls[i].OriginalURL
		}
	}

	return bodyResponseButch, nil
}

func (s *LocalStorageStruct) GetURL(id string) (string, error) {
	url, ok := s.ShortURL[id]
	if !ok {
		return url, fmt.Errorf("id %s not found", id)
	} else {
		return url, nil
	}
}

func (s *LocalStorageStruct) InitStorage() error {
	cfg := config.GetConfig()

	Cons, err := NewConsumer(cfg.FlagFileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer Cons.Close()

	urls, err := Cons.ReadShortURL()
	if err != nil {
		return err
	}

	for _, url := range urls {
		s.ShortURL[url.UUID] = url.OriginalURL
	}
	return nil
}

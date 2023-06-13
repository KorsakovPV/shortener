package storage

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage/dbstorage"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage/localstorage"
	"github.com/KorsakovPV/shortener/internal/models"
	"go.uber.org/zap"
)

type AbstractStorage interface {
	PutURL(string, string) (string, error)
	GetURL(string) (string, error)
	PutURLBatch([]models.RequestBatch) ([]models.ResponseButch, error)
	InitStorage() error
}

type Struct struct{}

func (s Struct) PutURL(id string, body string) (string, error) {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		id, err := dbStorage.PutURL(id, body)
		return id, err
	}

	_, err := localStorage.PutURL(id, body)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Struct) GetURL(id string) (string, error) {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		url, err := dbStorage.GetURL(id)
		if err != nil {
			return "", err
		}
		return url, err
	}

	url, err := localStorage.GetURL(id)
	if err != nil {
		return "", err
	}
	return url, err
}

func (s Struct) PutURLBatch(body []models.RequestBatch) ([]models.ResponseButch, error) {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		bodyResponseButch, err := dbStorage.PutURLBatch(body)
		if err != nil {
			return nil, err
		}
		return bodyResponseButch, nil
	}

	bodyResponseButch, err := localStorage.PutURLBatch(body)
	if err != nil {
		return nil, err
	}

	return bodyResponseButch, nil
}

func (s Struct) InitStorage() error {
	//TODO implement me
	panic("implement me")
}

var s AbstractStorage = Struct{}

var localStorage AbstractStorage = &localstorage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

var dbStorage AbstractStorage = &dbstorage.DBStorageStruct{}

func InitStorage() error {
	sugar := logging.GetSugarLogger()

	cfg := config.GetConfig()

	// Если в конфиге есть url для базы, то работаем с базой
	if cfg.FlagDataBaseDSN != "" {
		err := InitDBStorage(cfg, sugar)
		if err != nil {
			return err
		}
		return nil
	}

	// Если в конфиге нет url для базы, то работаем с файлом
	err := InitLocalStorage(cfg, sugar)
	if err != nil {
		return err
	}
	return nil
}

func InitDBStorage(cfg *config.Сonfiguration, sugar zap.SugaredLogger) error {
	err := dbStorage.InitStorage()
	if err != nil {
		sugar.Errorf("ERROR Init DB Storage. %s", err)
		return err
	}
	sugar.Infof("Use db storage %s", cfg.FlagDataBaseDSN)
	return nil
}

func InitLocalStorage(cfg *config.Сonfiguration, sugar zap.SugaredLogger) error {

	// Если в конфиге есть имя файла, то загружаем его
	if cfg.FlagFileStoragePath != "" {
		err := localStorage.InitStorage()
		if err != nil {
			sugar.Errorf("ERROR Init Local Storage. %s", err)
			return err
		}
		sugar.Infof("Use local storage %s", cfg.FlagFileStoragePath)
		return nil
	}

	// Если в конфиге нет имя файла, то работаем только с памятью
	sugar.Infoln("Use memory storage")
	return nil
}

func GetStorage() AbstractStorage {
	return s
}

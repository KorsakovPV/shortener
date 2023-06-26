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
	PutURL(string, string, interface{}) (string, error)
	GetURL(string) (string, error)
	PutURLBatch([]models.RequestBatch, interface{}) ([]models.ResponseButch, error)
	GetURLBatch(interface{}) ([]models.ResponseButchForUser, error)
	DeleteURLBatch([]string, interface{}) error
	InitStorage() error
}

type Struct struct{}

func (s Struct) PutURL(id string, body string, userID interface{}) (string, error) {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		id, err := dbStorage.PutURL(id, body, userID)
		return id, err
	}

	_, err := localStorage.PutURL(id, body, userID)
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

func (s Struct) GetURLBatch(userID interface{}) ([]models.ResponseButchForUser, error) {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		bodyResponseButch, err := dbStorage.GetURLBatch(userID)
		if err != nil {
			return nil, err
		}
		return bodyResponseButch, nil
	}

	bodyResponseButch, err := localStorage.GetURLBatch(userID)
	if err != nil {
		return nil, err
	}

	return bodyResponseButch, nil
}

func (s Struct) DeleteURLBatch(req []string, userID interface{}) error {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		err := dbStorage.DeleteURLBatch(req, userID)
		return err
	}

	err := localStorage.DeleteURLBatch(req, userID)
	return err
}

func (s Struct) PutURLBatch(body []models.RequestBatch, userID interface{}) ([]models.ResponseButch, error) {
	cfg := config.GetConfig()

	if cfg.FlagDataBaseDSN != "" {
		bodyResponseButch, err := dbStorage.PutURLBatch(body, userID)
		if err != nil {
			return nil, err
		}
		return bodyResponseButch, nil
	}

	bodyResponseButch, err := localStorage.PutURLBatch(body, userID)
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

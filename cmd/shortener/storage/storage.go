package storage

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage/dbstorage"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage/localstorage"
	"go.uber.org/zap"
)

type AbstractStorage interface {
	PutURL(string) (string, error)
	GetURL(string) (string, error)
	LoadBackupURL() error
}

var storage AbstractStorage

var localStorage AbstractStorage = &localstorage.LocalStorageStruct{
	ShortURL: map[string]string{},
}

var dbStorage AbstractStorage = &dbstorage.DBStorageStruct{
	//ShortURL: map[string]string{},
}

func InitStorage() error {
	sugar := logging.GetSugarLogger()

	cfg := config.GetConfig()

	// Если в конфиге есть url для базы, то работаем с базой
	//if cfg.FlagDataBaseDSN != "" {
	//	return InitDbStorage(sugar, cfg)
	//}

	// Если в конфиге нет url для базы, то работаем с файлом
	return InitLocalStorage(cfg, sugar)

}

//func InitDbStorage(sugar zap.SugaredLogger, cfg *config.Сonfiguration) error {
//	storage = dbStorage
//	// TODO добавить создание таблиц.
//	sugar.Infof("Use db storage %s", cfg.FlagDataBaseDSN)
//	return nil
//}

func InitLocalStorage(cfg *config.Сonfiguration, sugar zap.SugaredLogger) error {
	storage = localStorage

	// Если в конфиге есть имя файла, то загружаем его
	if cfg.FlagFileStoragePath != "" {
		err := storage.LoadBackupURL()
		if err != nil {
			sugar.Errorf("ERROR LoadBackupURL. %s", err)
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
	return storage
}

package logging

import "go.uber.org/zap"

var sugar zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = logger.Sync()
	}()

	sugar = *logger.Sugar()
}

func GetSugarLogger() zap.SugaredLogger {
	return sugar
}

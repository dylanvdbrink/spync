package storage

import (
	"errors"
	"go.uber.org/zap"
	"os"
)

func GetOrCreateFile(filepath string) (*os.File, bool, error) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	if _, existsErr := os.Stat(filepath); errors.Is(existsErr, os.ErrNotExist) {
		newFile, createErr := os.Create(filepath)
		if createErr != nil {
			sugar.Error("file creation error: " + createErr.Error())
			return nil, false, createErr
		}

		return newFile, true, nil
	} else {
		file, openErr := os.OpenFile(filepath, os.O_RDWR, 0755)
		if openErr != nil {
			sugar.Error("error opening file: " + openErr.Error())
			return nil, false, openErr
		}
		return file, false, nil
	}
}

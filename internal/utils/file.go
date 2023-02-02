package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var ErrFilePath = errors.New("the provided filepath is invalid")

func GetRealFilePath(filePath string) (string, error) {
	var err error
	filePath, err = homedir.Expand(filePath)
	if err != nil {
		return "", ErrorWithCause(ErrFilePath, err.Error())
	}
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return "", ErrorWithCause(ErrFilePath, err.Error())
	}
	return filePath, nil
}

func ReadFile(filePath string) ([]byte, error) {
	var err error
	filePath, err = GetRealFilePath(filePath)
	if err != nil {
		return nil, ErrorWithCause(ErrFilePath, err.Error())
	}
	var content []byte
	content, err = os.ReadFile(filePath)
	if err != nil {
		return nil, ErrorWithCause(ErrFilePath, err.Error())
	}

	return content, nil
}

package strategy

import (
	"os"
	"path"
	"path/filepath"
	"stregy/internal/domain/strategy"
	"stregy/pkg/utils"
)

type storage struct{}

func NewStorage() strategy.Storage {
	return storage{}
}

func (storage) SaveStrategy(name string, implementation *string) (string, error) {
	wd, _ := os.Getwd()
	dir := path.Join(wd, "local", "strategies", name)
	os.MkdirAll(dir, os.ModePerm)
	utils.RemoveContents(dir)

	archivePath := filepath.Join(dir, "strategy.zip")
	f, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	f.Write([]byte(*implementation))
	f.Close()

	utils.Unzip(archivePath, dir)
	os.Remove(archivePath)

	return dir, nil
}

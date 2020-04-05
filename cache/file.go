package cache

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// FileCache file struct
type FileCache struct {
	logger    *logrus.Entry
	cacheFile string
}

// NewFileCache create a new file cache
func NewFileCache(filePath string) *FileCache {
	return &FileCache{
		logger: logrus.
			WithField("file", filePath).
			WithField("component", "FileCache"),
		cacheFile: filePath,
	}
}

func (f *FileCache) Read() (ip string, err error) {
	file, err := os.Open(f.cacheFile)
	if err != nil {
		f.logger.Error("Error opening file")
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		f.logger.Error("Error reading from file")
		return "", err
	}
	err = file.Close()
	if err != nil {
		f.logger.Error("Error closing file")
		return "", err
	}
	return string(fileBytes), err
}

func (f *FileCache) Write(newIP string) (err error) {
	file, err := os.Create(f.cacheFile)
	if err != nil {
		f.logger.Error("Error creating cache file!")
		return err
	}
	_, err = file.WriteString(newIP)
	if err != nil {
		f.logger.Error("Error writing to cache file!")
		return err
	}
	log.Printf("Stored IP (%s) in cache.", newIP)
	err = file.Close()
	if err != nil {
		f.logger.Error("Error closing cache file!")
	}
	return
}

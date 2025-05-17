package cache

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// FileCache file struct
type FileCache struct {
	logger    *logrus.Entry
	cacheFile string
}

type CacheContents struct {
	IpAddress string
	CacheTime time.Time
}

// NewFileCache create a new file cache
func NewFileCache(filePath string, logger *logrus.Entry) *FileCache {
	return &FileCache{
		logger: logger.
			WithField("file", filePath).
			WithField("component", "FileCache"),
		cacheFile: filePath,
	}
}

func (f *FileCache) Read() (cacheContents *CacheContents, err error) {
	info, err := os.Stat(f.cacheFile)
	if err != nil || info.IsDir() {
		f.logger.Info(f.cacheFile + " is not a file")
		return &CacheContents{
			IpAddress: "",
		}, errors.New("Cache is not a valid file")
	}
	file, err := os.Open(f.cacheFile)
	if err != nil {
		f.logger.Error("Error opening file")
		return nil, err
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		f.logger.Error("Error reading from file")
		return nil, err
	}
	err = file.Close()
	if err != nil {
		f.logger.Error("Error closing file")
		return nil, err
	}
	var contents CacheContents
	err = json.Unmarshal(fileBytes, &contents)
	return &contents, err
}

func (f *FileCache) Write(newIP string) (err error) {
	file, err := os.Create(f.cacheFile)
	if err != nil {
		f.logger.Error("Error creating cache file!")
		return err
	}
	contents := CacheContents{
		IpAddress: newIP,
		CacheTime: time.Now(),
	}
	jsonBytes, err := json.MarshalIndent(contents, "", " ")
	if err != nil {
		f.logger.Error("Error marshalling JSON!")
		return err
	}
	_, err = file.Write(jsonBytes)
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

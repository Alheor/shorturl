// Package repository
// file implementation
package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/Alheor/shorturl/internal/config"
	"os"
	"sync"
)

type shortURL struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ShortNameFile struct
type ShortNameFile struct {
	URLMap map[string]string
	sync.RWMutex
	file *os.File
}

func (sn *ShortNameFile) Init() error {
	sn.URLMap = make(map[string]string)

	err := load(sn, config.Options.FileStoragePath)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(config.Options.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	sn.file = file

	return nil
}

func (sn *ShortNameFile) Add(id string, value string) error {

	_, err := os.Stat(sn.file.Name())
	if err != nil {
		panic(err)
	}

	sn.Lock()

	_, exists := sn.URLMap[id]
	if exists {
		sn.Unlock()
		return NewUniqueError(id, nil)
	}

	for _, mapValue := range sn.URLMap {
		if mapValue == value {
			sn.Unlock()
			return NewUniqueError(id, nil)
		}
	}

	data, err := json.Marshal(&shortURL{ID: id, URL: value})
	if err != nil {
		sn.Unlock()
		return err
	}

	sn.URLMap[id] = value
	data = append(data, '\n')

	_, err = sn.file.Write(data)

	sn.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (sn *ShortNameFile) Get(id string) (value string, error error) {

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.URLMap[id]
	if !exists {
		return ``, errors.New(ErrIDNotFound)
	}

	return url, nil
}

func (sn *ShortNameFile) Remove(id string) {
	panic(errors.New(`method "Remove" from file repository is restricted`))
}

func (sn *ShortNameFile) StorageIsReady() bool {

	return sn.file != nil
}

func (sn *ShortNameFile) AddBatch(in []BatchEl) error {
	_, err := os.Stat(sn.file.Name())
	if err != nil {
		panic(err)
	}

	sn.Lock()

	var res []byte
	for _, v := range in {

		data, err := json.Marshal(&shortURL{ID: v.ShortURL, URL: v.OriginalURL})
		if err != nil {
			sn.Unlock()
			return err
		}

		_, exists := sn.URLMap[v.ShortURL]
		if exists {
			sn.Unlock()
			return errors.New(ErrValueAlreadyExist)
		}

		for _, mapValue := range sn.URLMap {
			if mapValue == v.OriginalURL {
				sn.Unlock()
				return errors.New(ErrValueAlreadyExist)
			}
		}

		res = append(res, append(data, '\n')...)
	}

	_, err = sn.file.Write(res)

	sn.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func load(sn *ShortNameFile, path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	sn.Lock()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		el := shortURL{}
		err = json.Unmarshal(data, &el)
		if err != nil {
			continue
		}

		sn.URLMap[el.ID] = el.URL
	}

	sn.Unlock()

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

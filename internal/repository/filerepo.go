// Package repository
// file implementation
package repository

import (
	"bufio"
	"context"
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

func (sn *ShortNameFile) Init(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sn.URLMap = make(map[string]string)

	err := load(ctx, sn, config.Options.FileStoragePath)
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

func (sn *ShortNameFile) Add(ctx context.Context, id string, value string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := os.Stat(sn.file.Name())
	if err != nil {
		panic(err)
	}

	sn.Lock()
	defer sn.Unlock()

	_, exists := sn.URLMap[id]
	if exists {
		return NewUniqueError(id, nil)
	}

	for _, mapValue := range sn.URLMap {
		if mapValue == value {
			return NewUniqueError(id, nil)
		}
	}

	data, err := json.Marshal(&shortURL{ID: id, URL: value})
	if err != nil {
		return err
	}

	sn.URLMap[id] = value
	data = append(data, '\n')

	_, err = sn.file.Write(data)

	if err != nil {
		return err
	}

	return nil
}

func (sn *ShortNameFile) Get(ctx context.Context, id string) (value string, error error) {

	select {
	case <-ctx.Done():
		return ``, errors.New(ErrIDNotFound)
	default:
	}

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.URLMap[id]
	if !exists {
		return ``, errors.New(ErrIDNotFound)
	}

	return url, nil
}

func (sn *ShortNameFile) Remove(ctx context.Context, id string) {
	panic(errors.New(`method "Remove" from file repository is restricted`))
}

func (sn *ShortNameFile) IsReady(ctx context.Context) bool {

	select {
	case <-ctx.Done():
		return false
	default:
	}

	return sn.file != nil
}

func (sn *ShortNameFile) AddBatch(ctx context.Context, in []BatchEl) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := os.Stat(sn.file.Name())
	if err != nil {
		panic(err)
	}

	sn.Lock()
	defer sn.Unlock()

	var res []byte
	for _, v := range in {

		data, err := json.Marshal(&shortURL{ID: v.ShortURL, URL: v.OriginalURL})
		if err != nil {
			return err
		}

		_, exists := sn.URLMap[v.ShortURL]
		if exists {
			return errors.New(ErrValueAlreadyExist)
		}

		for _, mapValue := range sn.URLMap {
			if mapValue == v.OriginalURL {
				return errors.New(ErrValueAlreadyExist)
			}
		}

		res = append(res, append(data, '\n')...)
	}

	_, err = sn.file.Write(res)

	if err != nil {
		return err
	}

	return nil
}

func load(ctx context.Context, sn *ShortNameFile, path string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	sn.Lock()
	defer sn.Unlock()

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

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

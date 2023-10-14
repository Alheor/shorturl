// Package repository
// Short url repository
package repository

import (
	"errors"
	"sync"
)

// ShortName struct
type ShortName struct {
	urlMap map[string]string
	sync.RWMutex
}

// InitMap repository constructor
func InitMap() *ShortName {
	instance := new(ShortName)
	instance.urlMap = make(map[string]string)

	return instance
}

func (sn *ShortName) Add(id string, value string) error {

	_, exists := sn.urlMap[id]
	if exists {
		return errors.New(ErrorValueAlreadyExist)
	}

	for _, mapValue := range sn.urlMap {
		if mapValue == value {
			return errors.New(ErrorValueAlreadyExist)
		}
	}

	sn.Lock()
	defer sn.Unlock()

	sn.urlMap[id] = value

	return nil
}

func (sn *ShortName) Get(id string) (value string, error error) {

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.urlMap[id]
	if !exists {
		return ``, errors.New(ErrorIdNotFound)
	}

	return url, nil
}

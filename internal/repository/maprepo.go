// Package repository
// map implementation
package repository

import (
	"errors"
	"sync"
)

// ShortNameMap struct
type ShortNameMap struct {
	URLMap map[string]string
	sync.RWMutex
}

func (sn *ShortNameMap) Init() error {
	sn.URLMap = make(map[string]string)
	return nil
}

func (sn *ShortNameMap) Add(id string, value string) error {

	sn.Lock()

	_, exists := sn.URLMap[id]
	if exists {
		sn.Unlock()
		return errors.New(ErrValueAlreadyExist)
	}

	for _, mapValue := range sn.URLMap {
		if mapValue == value {
			sn.Unlock()
			return errors.New(ErrValueAlreadyExist)
		}
	}

	sn.URLMap[id] = value
	sn.Unlock()

	return nil
}

func (sn *ShortNameMap) Get(id string) (value string, error error) {

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.URLMap[id]
	if !exists {
		return ``, errors.New(ErrIDNotFound)
	}

	return url, nil
}

func (sn *ShortNameMap) Remove(id string) {

	sn.Lock()
	defer sn.Unlock()

	delete(sn.URLMap, id)
}

func (sn *ShortNameMap) StorageIsReady() bool {

	return sn.URLMap != nil
}

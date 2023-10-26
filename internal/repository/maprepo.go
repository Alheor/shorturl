// Package repository
// map implementation
package repository

import (
	"errors"
	"sync"
)

// ShortNameMap struct
type ShortNameMap struct {
	urlMap map[string]string
	sync.RWMutex
}

func (sn *ShortNameMap) Add(id string, value string) error {

	sn.RLock()

	_, exists := sn.urlMap[id]
	if exists {
		sn.RUnlock()
		return errors.New(ErrorValueAlreadyExist)
	}

	for _, mapValue := range sn.urlMap {
		if mapValue == value {
			sn.RUnlock()
			return errors.New(ErrorValueAlreadyExist)
		}
	}

	sn.RUnlock()

	sn.Lock()
	defer sn.Unlock()

	sn.urlMap[id] = value

	return nil
}

func (sn *ShortNameMap) Get(id string) (value string, error error) {

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.urlMap[id]
	if !exists {
		return ``, errors.New(ErrorIDNotFound)
	}

	return url, nil
}

func (sn *ShortNameMap) Remove(id string) {

	sn.RLock()
	defer sn.RUnlock()

	delete(sn.urlMap, id)
}

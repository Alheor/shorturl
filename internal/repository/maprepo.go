// Package repository
// map implementation
package repository

import (
	"errors"
	"sync"
)

// ShortNameMap struct
type ShortNameMap struct {
	UrlMap map[string]string
	sync.RWMutex
}

func (sn *ShortNameMap) Init() error {
	sn.UrlMap = make(map[string]string)
	return nil
}

func (sn *ShortNameMap) Add(id string, value string) error {

	sn.RLock()

	_, exists := sn.UrlMap[id]
	if exists {
		sn.RUnlock()
		return errors.New(ErrorValueAlreadyExist)
	}

	for _, mapValue := range sn.UrlMap {
		if mapValue == value {
			sn.RUnlock()
			return errors.New(ErrorValueAlreadyExist)
		}
	}

	sn.RUnlock()

	sn.Lock()
	defer sn.Unlock()

	sn.UrlMap[id] = value

	return nil
}

func (sn *ShortNameMap) Get(id string) (value string, error error) {

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.UrlMap[id]
	if !exists {
		return ``, errors.New(ErrorIDNotFound)
	}

	return url, nil
}

func (sn *ShortNameMap) Remove(id string) {

	sn.Lock()
	defer sn.Unlock()

	delete(sn.UrlMap, id)
}

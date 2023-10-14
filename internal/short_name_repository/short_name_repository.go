// Package short_name_repository
// Short url repository
package short_name_repository

import (
	"errors"
	"sync"
)

const (
	// ErrorURLNotFound error message
	ErrorURLNotFound = `url not found`

	// ErrorURLAlreadyExist error message
	ErrorURLAlreadyExist = `url already exist`
)

// ShortName struct
type ShortName struct {
	urlMap map[string]string
	sync.RWMutex
}

// Init repository constructor
func Init() *ShortName {
	instance := new(ShortName)
	instance.urlMap = make(map[string]string)

	return instance
}

func (sn *ShortName) Add(id string, value string) error {

	_, exists := sn.urlMap[id]
	if exists {
		return errors.New(ErrorURLAlreadyExist)
	}

	for _, value := range sn.urlMap {
		if value == value {
			return errors.New(ErrorURLAlreadyExist)
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
		return ``, errors.New(ErrorURLNotFound)
	}

	return url, nil
}

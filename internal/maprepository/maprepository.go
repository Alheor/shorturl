// Package maprepository
// Short url repository
package maprepository

import (
	"errors"
	"github.com/Alheor/shorturl/internal/repository"
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

	sn.RLock()
	//если сразу вызывать defer sn.RUnlock(), возникает deadlock

	_, exists := sn.urlMap[id]
	if exists {
		sn.RUnlock()
		return errors.New(repository.ErrorValueAlreadyExist)
	}

	for _, mapValue := range sn.urlMap {
		if mapValue == value {
			sn.RUnlock()
			return errors.New(repository.ErrorValueAlreadyExist)
		}
	}

	sn.RUnlock()

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
		return ``, errors.New(repository.ErrorIDNotFound)
	}

	return url, nil
}

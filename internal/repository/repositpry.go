// Package repository
// Short url repository
package repository

import (
	"errors"
)

const (
	// ErrorURLNotFound error message
	ErrorURLNotFound = `url not found`

	// ErrorURLAlreadyExist error message
	ErrorURLAlreadyExist = `url already exist`
)

// Repository interface
type Repository interface {
	AddURL(shortName string, url string) error
	GetURL(shortName string) (url string, error error)
}

// ShortName struct
type ShortName struct {
	urlMap map[string]string
}

// Init repository constructor
func (sn ShortName) Init() *ShortName {
	instance := new(ShortName)
	instance.urlMap = make(map[string]string)

	return instance
}

func (sn ShortName) AddURL(shortName string, url string) error {

	_, exists := sn.urlMap[shortName]
	if exists {
		return errors.New(ErrorURLAlreadyExist)
	}

	for _, value := range sn.urlMap {
		if value == url {
			return errors.New(ErrorURLAlreadyExist)
		}
	}

	sn.urlMap[shortName] = url

	return nil
}

func (sn ShortName) GetURL(shortName string) (url string, error error) {

	url, exists := sn.urlMap[shortName]
	if !exists {
		return ``, errors.New(ErrorURLNotFound)
	}

	return url, nil
}

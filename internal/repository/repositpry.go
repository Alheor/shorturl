// Package repository
// Short url repository
package repository

import (
	"errors"
)

const (
	// ErrorUrlNotFound error message
	ErrorUrlNotFound = `url not found`

	// ErrorUrlAlreadyExist error message
	ErrorUrlAlreadyExist = `url already exist`
)

// Repository interface
type Repository interface {
	AddUrl(shortName string, url string) error
	GetUrl(shortName string) (url string, error error)
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

func (sn ShortName) AddUrl(shortName string, url string) error {

	_, exists := sn.urlMap[shortName]
	if exists {
		return errors.New(ErrorUrlAlreadyExist)
	}

	for _, value := range sn.urlMap {
		if value == url {
			return errors.New(ErrorUrlAlreadyExist)
		}
	}

	sn.urlMap[shortName] = url

	return nil
}

func (sn ShortName) GetUrl(shortName string) (url string, error error) {

	url, exists := sn.urlMap[shortName]
	if !exists {
		return ``, errors.New(ErrorUrlNotFound)
	}

	return url, nil
}

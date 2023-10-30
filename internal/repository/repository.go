// Package repository
// Short url repository
package repository

import "github.com/Alheor/shorturl/internal/config"

const (
	// ErrorIDNotFound error message
	ErrorIDNotFound = `id not found`

	// ErrorValueAlreadyExist error message
	ErrorValueAlreadyExist = `value already exist`
)

// Repository interface
type Repository interface {
	Add(id string, value string) error
	Get(id string) (value string, error error)
	Remove(id string)
	Init() error
}

// Init repository constructor
func Init() Repository {

	var instance Repository

	if config.Options.FileStoragePath == `` {
		instance = new(ShortNameMap)
	} else {
		instance = new(ShortNameFile)
	}

	err := instance.Init()
	if err != nil {
		panic(err)
	}

	return instance
}

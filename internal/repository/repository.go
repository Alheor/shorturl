// Package repository
// Short url repository
package repository

import (
	"github.com/Alheor/shorturl/internal/config"
)

const (
	// ErrIDNotFound error message
	ErrIDNotFound = `id not found`

	// ErrValueAlreadyExist error message
	ErrValueAlreadyExist = `value already exist`
)

type BatchEl struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"-"`
	ShortURL      string `json:"short_url"`
}

type UniqueError struct {
	ShortKey string
	Err      error
}

func (e *UniqueError) Error() string {
	return e.Err.Error()
}

func NewUniqueError(shortKey string, err error) error {
	return &UniqueError{
		ShortKey: shortKey,
		Err:      err,
	}
}

// Repository interface
type Repository interface {
	Add(id string, value string) error
	AddBatch(in []BatchEl) error
	Get(id string) (value string, error error)
	Remove(id string)
	Init() error
	StorageIsReady() bool
}

// Init repository constructor
func Init() Repository {

	var instance Repository

	if config.Options.DatabaseDsn != `` {
		instance = new(Postgres)

	} else {
		if config.Options.FileStoragePath == `` {
			instance = new(ShortNameMap)

		} else {
			instance = new(ShortNameFile)
		}
	}

	err := instance.Init()
	if err != nil {
		panic(err)
	}

	return instance
}

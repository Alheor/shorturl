// Package repository
// Short url repository
package repository

import (
	"context"
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
	Add(ctx context.Context, id string, value string) error
	AddBatch(ctx context.Context, in []BatchEl) error
	Get(ctx context.Context, id string) (value string, error error)
	Remove(ctx context.Context, id string)
	Init(ctx context.Context) error
	IsReady(ctx context.Context) bool
}

// Init repository constructor
func Init(ctx context.Context) Repository {

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

	err := instance.Init(ctx)
	if err != nil {
		panic(err)
	}

	return instance
}

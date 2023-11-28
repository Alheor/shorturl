// Package repository
// Short url repository
package repository

import (
	"context"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/userauth"
)

const (
	// ErrIDNotFound error message
	ErrIDNotFound = `id not found`

	// ErrNotFound error message
	ErrNotFound = `not found`

	// ErrValueAlreadyExist error message
	ErrValueAlreadyExist = `value already exist`
)

type BatchEl struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"-"`
	ShortURL      string `json:"short_url"`
}

type HistoryEl struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type UniqueErr struct {
	ShortKey string
	Err      error
}

func (e *UniqueErr) Error() string {
	return e.Err.Error()
}

func NewUniqueError(shortKey string, err error) error {
	return &UniqueErr{
		ShortKey: shortKey,
		Err:      err,
	}
}

// Repository interface
type Repository interface {
	Add(ctx context.Context, ser *userauth.User, id string, value string) error
	AddBatch(ctx context.Context, user *userauth.User, in []BatchEl) error
	Get(ctx context.Context, user *userauth.User, id string) (value string, isDeleted bool, error error)
	GetAll(ctx context.Context, user *userauth.User) (list []HistoryEl, error error)
	Remove(ctx context.Context, user *userauth.User, id string)
	RemoveBatch(ctx context.Context, user *userauth.User, ids []string) error
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

package repository

import (
	"context"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"time"
)

var repo Repository

type Repository interface {
	Add(ctx context.Context, user *models.User, name string) (string, error)
	AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error
	GetByShortName(ctx context.Context, user *models.User, name string) (string, error)
	IsReady(ctx context.Context) bool
	RemoveByOriginalURL(ctx context.Context, user *models.User, url string) error
}

func Init(ctx context.Context, config *config.Options, repository Repository) error {

	if repository != nil {
		repo = repository
		return nil
	}

	if config.DatabaseDsn != `` {
		logger.Info(`Repository starting in database mode`)

		var db *pgxpool.Pool
		var err error

		if db, err = pgxpool.New(ctx, config.DatabaseDsn); err != nil {
			return err
		}

		logger.Info(`Running migrations ...`)

		schemaCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = createDBSchema(schemaCtx, db)
		if err != nil {
			return err
		}

		_ = goose.Up(stdlib.OpenDBFromPool(db), "./internal/migrations")

		repo = &PostgresRepo{Conn: db}

	} else if config.FileStoragePath != `` {
		logger.Info(`Repository starting in file mode`)

		fRepo := &FileRepo{list: make(map[string]map[string]string)}

		err := fRepo.Load(ctx, config.FileStoragePath)
		if err != nil {
			return err
		}

		repo = fRepo

	} else {
		logger.Info(`Repository starting in memory mode`)

		repo = &MemoryRepo{list: make(map[string]map[string]string)}
	}

	return nil
}

func GetRepository() Repository {
	return repo
}

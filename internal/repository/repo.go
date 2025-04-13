package repository

import (
	"context"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var repo Repository

type Repository interface {
	Add(ctx context.Context, name string) (string, error)
	AddBatch(ctx context.Context, list *[]models.BatchEl) error
	GetByShortName(ctx context.Context, name string) (string, error)
	IsReady(ctx context.Context) bool
	RemoveByOriginalURL(ctx context.Context, url string) error
}

func Init(ctx context.Context, config *config.Options, repository Repository) error {

	if repository != nil {
		repo = repository
		return nil
	}

	if config.DatabaseDsn != `` {
		logger.Info(`Repository starting in database mode ` + config.DatabaseDsn)

		db, err := pgxpool.New(ctx, config.DatabaseDsn)

		if err != nil {
			return err
		}

		repo = &PostgresRepo{Conn: db}

		schemaCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
		defer cancel()

		err = createDBSchema(schemaCtx, db)
		if err != nil {
			return err
		}

	} else if config.FileStoragePath != `` {
		logger.Info(`Repository starting in file mode`)

		fRepo := &FileRepo{list: make(map[string]string)}

		err := fRepo.Load(ctx, config.FileStoragePath)
		if err != nil {
			return err
		}

		repo = fRepo

	} else {
		logger.Info(`Repository starting in memory mode`)

		repo = &MemoryRepo{list: make(map[string]string)}
	}

	return nil
}

func GetRepository() Repository {
	return repo
}

package repository

import (
	"context"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var repo Repository

type Repository interface {
	Add(ctx context.Context, name string) (string, error)
	GetByShortName(ctx context.Context, name string) (string, error)
	IsReady(ctx context.Context) bool
}

func Init(ctx context.Context, repository Repository) error {

	if repository != nil {
		repo = repository
		return nil
	}

	if config.GetOptions().DatabaseDsn != `` {
		logger.Info(`Repository starting in database mode:` + config.GetOptions().DatabaseDsn)

		db, err := pgxpool.New(ctx, config.GetOptions().DatabaseDsn)

		if err != nil {
			return err
		}

		repo = &PostgresRepo{Conn: db}

		schemaCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
		defer cancel()

		createDBSchema(schemaCtx, db)

	} else if config.GetOptions().FileStoragePath != `` {
		logger.Info(`Repository starting in file mode`)

		fRepo := &FileRepo{list: make(map[string]string)}

		err := fRepo.Load(ctx, config.GetOptions().FileStoragePath)
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

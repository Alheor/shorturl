package repository

import (
	"context"
	"os"
	"time"

	"github.com/Alheor/shorturl/internal/config"

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

		db, err := pgxpool.New(ctx, config.GetOptions().DatabaseDsn)
		if err != nil {
			return err
		}

		pgRepo = &Postgres{Conn: db}

		schemaCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
		defer cancel()

		createDBSchema(schemaCtx, db)

		repo = pgRepo

	} else {
		fileRepo = &FileRepo{list: make(map[string]string)}

		err := load(ctx, config.GetOptions().FileStoragePath)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(config.GetOptions().FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		fileRepo.file = file
		repo = fileRepo
	}

	return nil
}

func GetRepository() Repository {
	return repo
}

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
var Connection *pgxpool.Pool

type Repository interface {
	Add(ctx context.Context, user *models.User, name string) (string, error)
	AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error
	GetByShortName(ctx context.Context, user *models.User, name string) (string, bool, error)
	IsReady(ctx context.Context) bool
	RemoveByOriginalURL(ctx context.Context, user *models.User, url string) error
	GetAll(ctx context.Context, user *models.User) (*map[string]string, error)
	RemoveBatch(ctx context.Context, user *models.User, list []string) error
}

func Init(ctx context.Context, config *config.Options, repository Repository) error {

	if repository != nil {
		repo = repository
		return nil
	}

	if config.DatabaseDsn != `` {
		logger.Info(`Repository starting in database mode`)

		var err error

		if Connection, err = pgxpool.New(ctx, config.DatabaseDsn); err != nil {
			return err
		}

		logger.Info(`Apply DB schema ...`)

		repo = &PostgresRepo{Conn: Connection}

		schemaCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
		defer cancel()

		err = createDBSchema(schemaCtx, Connection)
		if err != nil {
			return err
		}

		logger.Info(`done`)

	} else if config.FileStoragePath != `` {
		logger.Info(`Repository starting in file mode`)

		fRepo := &FileRepo{list: make(map[string]map[string]string)}

		err := fRepo.Load(ctx, config.FileStoragePath)
		if err != nil {
			return err
		}

		repo = fRepo

		logger.Info(`done`)

	} else {
		logger.Info(`Repository starting in memory mode`)

		repo = &MemoryRepo{list: make(map[string]map[string]string)}

		logger.Info(`done`)
	}

	return nil
}

func GetRepository() Repository {
	return repo
}

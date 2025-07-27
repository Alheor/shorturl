package repository

import (
	"context"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Экземпляр репозитория.
var repo IRepository

// Connection - активное подключение к БД(если используется БД).
var Connection *pgxpool.Pool

// IRepository - интерфейс репозитория.
type IRepository interface {
	// Add - добавить URL.
	Add(ctx context.Context, user *models.User, name string) (string, error)

	// AddBatch - добавить несколько URL.
	AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error

	// GetByShortName - получить 1 URL.
	GetByShortName(ctx context.Context, user *models.User, name string) (string, bool, error)

	// IsReady - проверка работоспособности репозитория.
	IsReady(ctx context.Context) bool

	// RemoveByOriginalURL - удалить URL.
	RemoveByOriginalURL(ctx context.Context, user *models.User, url string) error

	// GetAll - получить все URL.
	GetAll(ctx context.Context, user *models.User) (<-chan models.HistoryEl, <-chan error)

	// RemoveBatch - удалить несколько URL.
	RemoveBatch(ctx context.Context, user *models.User, list []string) error

	// GetStats Статистика по пользователям и сокращенным URL
	GetStats(ctx context.Context) (*models.APIStatsResponse, error)

	Close()
}

// URL - структура URL элемента.
type URL struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
	URL    string `json:"url"`
}

// Init - инициализация репозитория, определение типа.
func Init(ctx context.Context, config *config.Options, repository IRepository) error {

	if repository != nil {
		repo = repository
		return nil
	}

	if config.DatabaseDsn != `` {
		logger.Info(`IRepository starting in database mode`)

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

	} else if config.FileStoragePath != `` {
		logger.Info(`IRepository starting in file mode`)

		fRepo := &FileRepo{list: make(map[string]map[string]string)}

		err := fRepo.load(ctx, config.FileStoragePath)
		if err != nil {
			return err
		}

		repo = fRepo

	} else {
		logger.Info(`IRepository starting in memory mode`)

		repo = &MemoryRepo{list: make(map[string]map[string]string)}
	}

	logger.Info(`done`)

	return nil
}

// GetRepository - метод получения текущего экземпляра репозитория.
func GetRepository() IRepository {
	return repo
}

// Package service - основные функции сервиса сокращения URL адресов.
//
// # Описание
//
// Представляет собой набор функций - моделей бизнес процесса.
//
// • Добавление 1 URL и получение его сокращенной версии в ответ.
//
// • Получение 1 URL по сокращенной версии.
//
// • Массовое добавление URL и получение их сокращенной версии в ответ.
//
// • Получение всех сокращенных URL.
//
// • Массовое удаление URL.
//
// • Проверка работоспособности репозитория.
package service

import (
	"context"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

var baseHost string

// Init Подготовка сервиса к работе
func Init(config *config.Options) {
	baseHost = config.BaseHost
}

// Add Добавление 1 URL и получение его сокращенной версии в ответ.
func Add(ctx context.Context, user *models.User, URL string) (string, error) {

	var err error
	var shortURL string
	if shortURL, err = repository.GetRepository().Add(ctx, user, URL); err != nil {
		logger.Error(`add url error: `, err)
		return ``, err
	}

	return shortURL, nil
}

// Get Получение 1 URL по сокращенной версии.
func Get(ctx context.Context, user *models.User, shortName string) (URL string, isRemoved bool) {
	str, isRemoved, err := repository.GetRepository().GetByShortName(ctx, user, shortName)
	if err != nil {
		logger.Error(`get url error: `, err)
		return ``, false
	}

	return str, isRemoved
}

// AddBatch Массовое добавление URL и получение их сокращенной версии в ответ.
func AddBatch(ctx context.Context, user *models.User, batch []models.APIBatchRequestEl) ([]models.APIBatchResponseEl, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	list := make([]models.BatchEl, 0, len(batch))

	for _, v := range batch {
		list = append(list, models.BatchEl{
			CorrelationID: v.CorrelationID,
			OriginalURL:   v.OriginalURL,
			ShortURL:      urlhasher.GetHash(v.OriginalURL),
		})
	}

	err := repository.GetRepository().AddBatch(ctx, user, &list)
	if err != nil {
		logger.Error(`add batch url error: `, err)
		return nil, err
	}

	resList := make([]models.APIBatchResponseEl, 0, len(batch))
	for _, v := range list {
		resList = append(resList, models.APIBatchResponseEl{
			CorrelationID: v.CorrelationID,
			ShortURL:      baseHost + `/` + v.ShortURL,
		})
	}

	return resList, nil
}

// GetAll Получение всех сокращенных URL.
func GetAll(ctx context.Context, user *models.User) (<-chan models.HistoryEl, <-chan error) {
	return repository.GetRepository().GetAll(ctx, user)
}

// RemoveBatch Массовое удаление URL.
func RemoveBatch(ctx context.Context, user *models.User, list []string) error {
	err := repository.GetRepository().RemoveBatch(ctx, user, list)
	if err != nil {
		logger.Error(`remove batch url error: `, err)
		return err
	}

	return nil
}

// IsDBReady Проверка работоспособности репозитория.
func IsDBReady(ctx context.Context) bool {
	return repository.GetRepository().IsReady(ctx)
}

// GetStats Статистика по пользователям и сокращенным URL
func GetStats(ctx context.Context) (*models.APIStatsResponse, error) {
	return repository.GetRepository().GetStats(ctx)
}

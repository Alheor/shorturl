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

func Init(config *config.Options) {
	baseHost = config.BaseHost
}

func Add(ctx context.Context, user *models.User, URL string) (string, error) {

	var err error
	var shortURL string
	if shortURL, err = repository.GetRepository().Add(ctx, user, URL); err != nil {
		return ``, err
	}

	return shortURL, nil
}

func Get(ctx context.Context, user *models.User, shortName string) string {
	str, err := repository.GetRepository().GetByShortName(ctx, user, shortName)
	if err != nil {
		logger.Error(`get url error: `, err)
		return ``
	}

	return str
}

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

func IsDBReady(ctx context.Context) bool {
	return repository.GetRepository().IsReady(ctx)
}

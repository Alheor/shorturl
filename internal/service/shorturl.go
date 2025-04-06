package service

import (
	"context"
	"time"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

func Add(ctx context.Context, URL string) (string, error) {

	var err error
	var shortURL string
	if shortURL, err = repository.GetRepository().Add(ctx, URL); err != nil {
		return ``, err
	}

	return shortURL, nil
}

func Get(ctx context.Context, shortName string) string {
	str, err := repository.GetRepository().GetByShortName(ctx, shortName)
	if err != nil {
		logger.Error(`get url error: `, err)
		return ``
	}

	return str
}

func AddBatch(ctx context.Context, batch []models.APIBatchRequestEl) ([]models.APIBatchResponseEl, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	list := make([]models.BatchEl, 0, len(batch))
	snGen := urlhasher.GetShortNameGenerator()

	for _, v := range batch {
		list = append(list, models.BatchEl{
			CorrelationID: v.CorrelationID,
			OriginalURL:   v.OriginalURL,
			ShortURL:      snGen.Generate(),
		})
	}

	err := repository.GetRepository().AddBatch(ctx, &list)
	if err != nil {
		return nil, err
	}

	resList := make([]models.APIBatchResponseEl, 0, len(batch))
	for _, v := range list {
		resList = append(resList, models.APIBatchResponseEl{
			CorrelationID: v.CorrelationID,
			ShortURL:      v.ShortURL,
		})
	}

	return resList, nil
}

func IsDBReady(ctx context.Context) bool {
	return repository.GetRepository().IsReady(ctx)
}

package service

import (
	"context"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/repository"
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

func IsDBReady(ctx context.Context) bool {
	return repository.GetRepository().IsReady(ctx)
}

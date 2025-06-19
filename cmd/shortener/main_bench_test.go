package main

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

var triesN = 100

func BenchmarkApiAddUrlWithFile(b *testing.B) {
	cfg := config.Load()

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	b.ResetTimer()
	b.Run("AddUrl file storage", func(b *testing.B) {
		b.StopTimer()

		_ = os.Remove(cfg.FileStoragePath)
		repository.Init(ctx, &cfg, nil)

		b.StartTimer()

		for i := 0; i < triesN; i++ {
			service.Add(ctx, user, targetURL+`/test`+strconv.Itoa(i))
		}
	})
}

func BenchmarkApiAddUrlWithMap(b *testing.B) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	b.ResetTimer()
	b.Run("AddUrl map storage", func(b *testing.B) {
		b.StopTimer()

		_ = os.Remove(cfg.FileStoragePath)
		repository.Init(ctx, &cfg, nil)

		b.StartTimer()

		for i := 0; i < triesN; i++ {
			service.Add(ctx, user, targetURL+`/test`+strconv.Itoa(i))
		}
	})
}

func BenchmarkApiAddUrlWithDB(b *testing.B) {

	b.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	repository.Connection.Exec(ctx, `TRUNCATE short_url`)

	b.ResetTimer()
	b.Run("AddUrl DB storage", func(b *testing.B) {
		b.StopTimer()

		repository.Connection.Exec(ctx, `TRUNCATE short_url`)

		b.StartTimer()

		for i := 0; i < triesN; i++ {
			service.Add(ctx, user, targetURL+`/test`+strconv.Itoa(i))
		}
	})
}

func BenchmarkApiAddBatchUrlsWithFile(b *testing.B) {
	cfg := config.Load()

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	b.ResetTimer()
	b.Run("AddBatchUrls file storage", func(b *testing.B) {
		b.StopTimer()

		_ = os.Remove(cfg.FileStoragePath)
		repository.Init(ctx, &cfg, nil)

		b.StartTimer()

		service.AddBatch(ctx, user, request)
	})
}

func BenchmarkApiAddBatchUrlsWithMap(b *testing.B) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	b.ResetTimer()
	b.Run("AddBatchUrls map storage", func(b *testing.B) {
		b.StopTimer()

		_ = os.Remove(cfg.FileStoragePath)
		repository.Init(ctx, &cfg, nil)

		b.StartTimer()

		service.AddBatch(ctx, user, request)
	})
}

func BenchmarkApiAddBatchUrlsWithDB(b *testing.B) {

	b.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	b.ResetTimer()
	b.Run("AddBatchUrls DB storage", func(b *testing.B) {
		b.StopTimer()

		repository.Connection.Exec(ctx, `TRUNCATE short_url`)

		b.StartTimer()

		service.AddBatch(ctx, user, request)
	})
}

func BenchmarkApiGetAllUrlsWithFile(b *testing.B) {
	cfg := config.Load()

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Get all urls file storage", func(b *testing.B) {
		for i := 0; i < triesN; i++ {
			chList, chErr := service.GetAll(ctx, user)
			for range chList {
			}
			for range chErr {
			}
		}
	})
}

func BenchmarkApiGetAllUrlsWithMap(b *testing.B) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Get all urls map storage", func(b *testing.B) {
		for i := 0; i < triesN; i++ {
			chList, chErr := service.GetAll(ctx, user)
			for range chList {
			}
			for range chErr {
			}
		}
	})
}

func BenchmarkApiGetAllUrlsWithDB(b *testing.B) {

	b.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	repository.Connection.Exec(ctx, `TRUNCATE short_url`)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Get all urls DB storage", func(b *testing.B) {
		for i := 0; i < triesN; i++ {
			chList, chErr := service.GetAll(ctx, user)
			for range chList {
			}
			for range chErr {
			}
		}
	})
}

func BenchmarkGetUrlWithFile(b *testing.B) {
	cfg := config.Load()

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Get url file storage", func(b *testing.B) {
		for i := 0; i < triesN; i++ {
			hash := urlhasher.GetHash(targetURL + `/test` + strconv.Itoa(i))
			service.Get(ctx, user, hash)
		}
	})
}

func BenchmarkGetUrlWithMap(b *testing.B) {
	cfg := config.Load()
	cfg.FileStoragePath = ``

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Get url map storage", func(b *testing.B) {
		for i := 0; i < triesN; i++ {
			hash := urlhasher.GetHash(targetURL + `/test` + strconv.Itoa(i))
			service.Get(ctx, user, hash)
		}
	})
}

func BenchmarkGetUrlWithDB(b *testing.B) {

	b.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	repository.Connection.Exec(ctx, `TRUNCATE short_url`)

	var request []models.APIBatchRequestEl

	for i := 0; i < triesN; i++ {
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   targetURL + `/test` + strconv.Itoa(i),
			CorrelationID: strconv.Itoa(i),
		})
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Get url map storage", func(b *testing.B) {
		for i := 0; i < triesN; i++ {
			hash := urlhasher.GetHash(targetURL + `/test` + strconv.Itoa(i))
			service.Get(ctx, user, hash)
		}
	})
}

func BenchmarkDeleteUrlWithDB(b *testing.B) {
	cfg := config.Load()
	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	logger.Init(nil)
	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx := context.Background()

	repository.Init(ctx, &cfg, nil)

	repository.Connection.Exec(ctx, `TRUNCATE short_url`)

	var request []models.APIBatchRequestEl
	var list []string
	for i := 0; i < triesN; i++ {
		url := targetURL + `/test` + strconv.Itoa(i)
		request = append(request, models.APIBatchRequestEl{
			OriginalURL:   url,
			CorrelationID: strconv.Itoa(i),
		})

		list = append(list, urlhasher.GetHash(url))
	}

	service.AddBatch(ctx, user, request)

	b.ResetTimer()
	b.Run("Remove batch DB storage", func(b *testing.B) {
		service.RemoveBatch(ctx, user, list)
	})
}

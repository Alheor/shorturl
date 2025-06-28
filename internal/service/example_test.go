package service

import (
	"context"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
)

func ExampleAdd() {
	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}
	ctx := context.Background()

	shortURL, err := Add(ctx, user, `https://example.com/?var1=value1&var2=value2`)
	if err != nil {
		logger.Error(`add url error`, err)
		return
	}

	println(shortURL)
}

func ExampleGet() {
	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}
	ctx := context.Background()

	originalURL, isRemoved := Get(ctx, user, `short_name`)
	if originalURL == `` {
		println(`URL not found`)
		return
	}

	println(originalURL)
	println(isRemoved)
}

func ExampleAddBatch() {
	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}
	ctx := context.Background()

	var list []models.APIBatchRequestEl
	var res []models.APIBatchResponseEl

	list = append(list, models.APIBatchRequestEl{CorrelationID: `1`, OriginalURL: `https://example.com/?var1=value1&var2=value2`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `2`, OriginalURL: `https://example.com/?var2=value2&var3=value3`})
	list = append(list, models.APIBatchRequestEl{CorrelationID: `3`, OriginalURL: `https://example.com/?var3=value3&var4=value4`})

	res, err := AddBatch(ctx, user, list)
	if err != nil {
		logger.Error(`add batch url error`, err)
		return
	}

	for _, v := range res {
		println(v.CorrelationID + `:` + v.ShortURL)
	}
}

func ExampleGetAll() {
	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}
	ctx := context.Background()

	chList, chErr := GetAll(ctx, user)

	for el := range chList {
		println(el.ShortURL + `:` + el.OriginalURL)
	}

	for err := range chErr {
		logger.Error(`get all urls error`, err)
	}
}

func ExampleRemoveBatch() {
	user := &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}
	ctx := context.Background()

	var list []string
	list = append(list, `short_name1`)
	list = append(list, `short_name2`)
	list = append(list, `short_name3`)

	err := RemoveBatch(ctx, user, list)
	if err != nil {
		logger.Error(`remove batch urls error`, err)
		return
	}

	println(`done`)
}

func ExampleIsDBReady() {
	isReady := IsDBReady(context.Background())

	println(isReady)
}

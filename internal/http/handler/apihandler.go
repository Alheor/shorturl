package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/userauth"
)

// AddShorten API обработчик запроса на добавление URL пользователя.
func AddShorten(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "AddShorten" handler`)

	var body []byte
	var err error
	var request models.APIRequest
	var response models.APIResponse

	defer req.Body.Close()
	if body, err = io.ReadAll(req.Body); err != nil || len(body) == 0 {
		response = models.APIResponse{Error: `url required`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		response = models.APIResponse{Error: `url required`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	if request.URL == `` {
		response = models.APIResponse{Error: `url required`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	if _, err = url.ParseRequestURI(request.URL); err != nil {
		response = models.APIResponse{Error: `Url invalid`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		response = models.APIResponse{Error: `Unauthorized`, StatusCode: http.StatusUnauthorized}
		sendAPIResponse(resp, &response)
		return
	}

	var shortURL string
	if shortURL, err = service.Add(ctx, user, request.URL); err != nil {

		var uniqErr *models.UniqueErr
		if errors.As(err, &uniqErr) {
			response = models.APIResponse{
				Result:     baseHost + `/` + uniqErr.ShortKey,
				StatusCode: http.StatusConflict,
			}

			sendAPIResponse(resp, &response)
			return
		}

		response = models.APIResponse{Error: `Internal error`, StatusCode: http.StatusInternalServerError}
		sendAPIResponse(resp, &response)
		return
	}

	response = models.APIResponse{
		Result:     baseHost + `/` + shortURL,
		StatusCode: http.StatusCreated,
	}

	sendAPIResponse(resp, &response)
}

// AddShortenBatch API обработчик запроса на массовое добавление URL пользователя.
func AddShortenBatch(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "AddShortenBatch" handler`)

	var body []byte
	var err error
	var request []models.APIBatchRequestEl
	var response []models.APIBatchResponseEl

	defer req.Body.Close()
	if body, err = io.ReadAll(req.Body); err != nil || len(body) == 0 {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusBadRequest})
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusBadRequest})
		return
	}

	if len(request) == 0 {
		sendAPIResponse(resp, &models.APIResponse{Error: `empty url list`, StatusCode: http.StatusBadRequest})
		return
	}

	for _, v := range request {
		if _, err = url.ParseRequestURI(v.OriginalURL); err != nil {
			sendAPIResponse(resp, &models.APIResponse{Error: `Url '` + v.OriginalURL + `' invalid`, StatusCode: http.StatusBadRequest})
			return
		}

		if v.CorrelationID == `` {
			sendAPIResponse(resp, &models.APIResponse{Error: `empty correlation_id`, StatusCode: http.StatusBadRequest})
			return
		}
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err = service.AddBatch(ctx, user, request)
	if err != nil {
		logger.Error(`Add batch error`, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawByte, err := json.Marshal(response)
	if err != nil {
		logger.Error(`response marshal error`, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add(HeaderContentType, HeaderContentTypeJSON)
	resp.WriteHeader(http.StatusCreated)

	_, err = resp.Write(rawByte)
	if err != nil {
		logger.Error(`write response error`, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetAllShorten API обработчик запроса на получение всех URL пользователя.
func GetAllShorten(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "GetAllShorten" handler`)

	var response models.APIResponse

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		response = models.APIResponse{Error: `Unauthorized`, StatusCode: http.StatusUnauthorized}
		sendAPIResponse(resp, &response)
		return
	}

	resp.Header().Add(HeaderContentType, HeaderContentTypeJSON)

	chList, chErr := service.GetAll(ctx, user)
	first := true
	hasEls := false

	for el := range chList {
		hasEls = true

		if first {
			_, err := resp.Write([]byte("["))
			if err != nil {
				logger.Error(`write response error`, err)
				resp.WriteHeader(http.StatusInternalServerError)

				return
			}

		} else {
			_, err := resp.Write([]byte(","))
			if err != nil {
				logger.Error(`write response error`, err)
				resp.WriteHeader(http.StatusInternalServerError)

				return
			}
		}
		first = false

		short := strings.TrimRight(baseHost, `/`) + `/` + el.ShortURL
		h := models.HistoryEl{OriginalURL: el.OriginalURL, ShortURL: short}
		rawByte, err := json.Marshal(h)
		if err != nil {
			logger.Error(`response marshal error`, err)
			resp.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = resp.Write(rawByte)
		if err != nil {
			logger.Error(`write response error`, err)
			resp.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	if hasEls {
		_, err := resp.Write([]byte("]"))
		if err != nil {
			logger.Error(`write response error`, err)
			resp.WriteHeader(http.StatusInternalServerError)

			return
		}

	}

	for err := range chErr {
		logger.Error(`Get all urls error`, err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !hasEls {
		resp.WriteHeader(http.StatusNoContent)
	}
}

// DeleteShorten API обработчик запроса на удаление URL пользователя.
func DeleteShorten(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "DeleteShorten" handler`)

	var request []string
	var response models.APIResponse
	var body []byte
	var err error

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		response = models.APIResponse{Error: `Unauthorized`, StatusCode: http.StatusUnauthorized}
		sendAPIResponse(resp, &response)
		return
	}

	defer req.Body.Close()
	if body, err = io.ReadAll(req.Body); err != nil || len(body) == 0 {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusBadRequest})
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusBadRequest})
		return
	}

	err = service.RemoveBatch(ctx, user, request)
	if err != nil {
		logger.Error(`Remove batch urls error`, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusAccepted)
}

// Stats Статистика по пользователям и сокращенным URL
func Stats(resp http.ResponseWriter, req *http.Request) {
	logger.Info(`Used "Stats" handler`)

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	stats, err := service.GetStats(ctx)
	if err != nil {
		logger.Error(`Get stats error`, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add(HeaderContentType, HeaderContentTypeJSON)

	rawByte, err := json.Marshal(stats)
	if err != nil {
		logger.Error(`response marshal error`, err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = resp.Write(rawByte)
	if err != nil {
		logger.Error(`write response error`, err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
}

// Подготовка ответа.
func sendAPIResponse(respWr http.ResponseWriter, resp *models.APIResponse) {
	rawByte, err := json.Marshal(resp)
	if err != nil {
		logger.Error(`response marshal error`, err)
		respWr.WriteHeader(http.StatusInternalServerError)
		return
	}

	respWr.Header().Add(HeaderContentType, HeaderContentTypeJSON)
	respWr.WriteHeader(resp.StatusCode)

	_, err = respWr.Write(rawByte)
	if err != nil {
		logger.Error(`write response error`, err)
		respWr.WriteHeader(http.StatusInternalServerError)
		return
	}
}

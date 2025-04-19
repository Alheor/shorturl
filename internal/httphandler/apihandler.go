package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Alheor/shorturl/internal/auth"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/service"
)

func AddShorten(resp http.ResponseWriter, req *http.Request) {

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

	user := auth.GetUser(ctx)
	if user == nil {
		response = models.APIResponse{Error: `Unauthorized`, StatusCode: http.StatusUnauthorized}
		sendAPIResponse(resp, &response)
		return
	}

	var shortURL string
	if shortURL, err = service.Add(ctx, request.URL); err != nil {

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

func AddShortenBatch(resp http.ResponseWriter, req *http.Request) {
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
			sendAPIResponse(resp, &models.APIResponse{Error: `Url ` + v.OriginalURL + ` invalid`, StatusCode: http.StatusBadRequest})
			return
		}

		if v.CorrelationID == `` {
			sendAPIResponse(resp, &models.APIResponse{Error: `empty correlation_id`, StatusCode: http.StatusBadRequest})
			return
		}
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user := auth.GetUser(ctx)
	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{Error: `Unauthorized`, StatusCode: http.StatusUnauthorized})
		return
	}

	response, err = service.AddBatch(req.Context(), request)
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

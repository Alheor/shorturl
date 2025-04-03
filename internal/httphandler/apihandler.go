package httphandler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
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

	var shortURL *string
	if shortURL, err = repository.GetRepository().Add(request.URL); err != nil {
		response = models.APIResponse{Error: `Internal error`, StatusCode: http.StatusInternalServerError}
		sendAPIResponse(resp, &response)
		return
	}

	response = models.APIResponse{
		Result:     config.GetOptions().BaseHost + `/` + *shortURL,
		StatusCode: http.StatusCreated,
	}

	sendAPIResponse(resp, &response)
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
